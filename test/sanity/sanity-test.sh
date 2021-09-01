#!/bin/bash

set -eou pipefail

trap 'kill $(jobs -p); ps ax | grep exe/main | grep -v grep | awk '\''{print $1}'\'' | xargs kill; exit 0' EXIT

NAMESPACE="${1:-openshift-logging}"
KUBECTL="${KUBECTL:-kubectl}"
DEPLOY_WITH_GATEWAY="${DEPLOY_WITH_GATEWAY:-false}"

deploy_loki_operator() {
    echo "--> deploying loki operator"
    if $DEPLOY_WITH_GATEWAY; then
      go run ./main.go --with-lokistack-gateway &
    else
      go run ./main.go &
    fi
  	sleep 10
}

wait_for_loki_stack() {
    echo "--> waiting for loki stack to be ready"
    loop_count=0
    while [[ $(kubectl -n "$NAMESPACE" get "deploy/loki-query-frontend-lokistack-dev" -o 'jsonpath={..status.conditions[?(@.type=="Available")].status}') != "True" ]]; do
      echo "waiting for loki deployment ..." && sleep 10; ((loop_count=loop_count+1))
      if [[ "$loop_count" -gt 30 ]]; then
           echo "Exit, fail waiting for loki deployment !"
           exit 1
      fi
    done

    $KUBECTL -n "$NAMESPACE" rollout status "deploy/loki-query-frontend-lokistack-dev" --timeout=300s
    $KUBECTL -n "$NAMESPACE" rollout status "deploy/loki-distributor-lokistack-dev" --timeout=300s
    $KUBECTL -n "$NAMESPACE" rollout status "statefulsets/loki-ingester-lokistack-dev" --timeout=300s
    $KUBECTL -n "$NAMESPACE" rollout status "statefulsets/loki-querier-lokistack-dev" --timeout=300s
    $KUBECTL -n "$NAMESPACE" rollout status "statefulsets/loki-compactor-lokistack-dev" --timeout=300s
}

expose_services() {
    echo "--> exposing services"
    if $DEPLOY_WITH_GATEWAY; then
      $KUBECTL -n "$NAMESPACE" port-forward svc/lokistack-gateway-http-lokistack-dev 13100:3100 &
      $KUBECTL -n "$NAMESPACE" port-forward svc/lokistack-gateway-http-lokistack-dev 23100:3100 &
    else
      $KUBECTL -n "$NAMESPACE" port-forward svc/loki-distributor-http-lokistack-dev 13100:3100 &
      $KUBECTL -n "$NAMESPACE" port-forward svc/loki-query-frontend-http-lokistack-dev 23100:3100 &
    fi
    $KUBECTL -n "$NAMESPACE" port-forward svc/loki-ingester-http-lokistack-dev 33100:3100 &
    sleep 10
}

deploy_minio() {
    echo "--> deploying minio"
    $KUBECTL -n "$NAMESPACE" delete pod minio ||:
    $KUBECTL -n "$NAMESPACE" run minio --env MINIO_ACCESS_KEY="user" --env MINIO_SECRET_KEY="password" --image=bitnami/minio:latest
    $KUBECTL -n "$NAMESPACE" wait --for=condition=ready pod/minio
    sleep 20
    $KUBECTL -n "$NAMESPACE" exec -i minio mc mb local/bucket
    $KUBECTL -n "$NAMESPACE" expose pod/minio --port 9000
    $KUBECTL -n "$NAMESPACE" port-forward  pod/minio 19000:9000 &
}

deploy_storage_secret() {
    echo "--> deploying minio storage secret"
    $KUBECTL -n "$NAMESPACE" delete secret test ||:
    $KUBECTL -n "$NAMESPACE" create secret generic test \
      --from-literal=endpoint="s3://user:password@minio.openshift-logging.svc.cluster.local:9000/bucket"\
      --from-literal=region="local" \
      --from-literal=bucketnames="bucket" \
      --from-literal=access_key_id="user" \
      --from-literal=access_key_secret="password"
}

deploy_lokistack-crd() {
    echo "--> deploying lokistack crd"
    make olm-deploy-example-lokistack-crd
}

push_log_line() {
    echo "-->  sending a log line to Loki"
    timestamp=$(date +"%s000000000")
    curl -v -H "Content-Type: application/json" -XPOST -s "http://localhost:13100/loki/api/v1/push" --data-raw '{"streams": [{ "stream": { "foo": "bar2" }, "values": [ [ "'$timestamp'", "fizzbuzz" ] ] }]}' ||:
    sleep 1
}

flush_to_storage () {
    echo "--> flush log lines to storage (minio)"
    curl -v -H "Content-Type: application/json" -XPOST -s "http://localhost:33100/flush" ||:
    sleep 1
    ls_result=$($KUBECTL -n "$NAMESPACE" exec -i minio mc ls local/bucket/fake)
    # storage must include 'fake' folder with one file sized 255 bytes
    if [[ $ls_result != *"255B"* ]]; then
      echo "***************************"
      echo "*** Sanity test failed! ***"
      echo "***************************"
      exit 1
    fi
}

query_logs() {
    echo "--> reading log lines from Loki"
    query_result=$(curl -G -s  "http://localhost:23100/loki/api/v1/query_range" --data-urlencode 'query={foo="bar2"}')
    # query must include log line with text 'fizzbuzz'
    if [[ $query_result != *"fizzbuzz"* ]]; then
      echo "***************************"
      echo "*** Sanity test failed! ***"
      echo "***************************"
      exit 2
    fi
}

clean_namespace (){
    echo "--> cleaning the namespace"
    $KUBECTL delete namespace "$NAMESPACE" ||:
    $KUBECTL create namespace "$NAMESPACE"
}

main() {
    clean_namespace
    deploy_storage_secret
    deploy_lokistack-crd
    deploy_minio
    deploy_loki_operator
    wait_for_loki_stack
    expose_services
    push_log_line
    flush_to_storage
    query_logs
    echo "***************************"
    echo "*** Sanity test passed! ***"
    echo "***************************"
    exit 0
}

main
