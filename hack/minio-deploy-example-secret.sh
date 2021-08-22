#!/bin/bash

set -eou pipefail

NAMESPACE=$1

REGION=""
ENDPOINT=""
ACCESS_KEY_ID=""
SECRET_ACCESS_KEY=""

set_credentials_from_aws() {
  REGION="eu1"
  ENDPOINT="https://minio.loki.svc.cluster.local:9000"
  ACCESS_KEY_ID=$(kubectl get secret minio-tenant-1-creds-secret -n loki -o jsonpath={.data.accesskey} |base64 -d)
  SECRET_ACCESS_KEY=$(kubectl get secret minio-tenant-1-creds-secret -n loki -o jsonpath={.data.secretkey} |base64 -d)
}

create_secret() {
  kubectl -n $NAMESPACE delete secret test ||:
  kubectl -n $NAMESPACE create secret generic test \
    --from-literal=endpoint=$(echo -n "$ENDPOINT") \
    --from-literal=region=$(echo -n "$REGION") \
    --from-literal=bucketnames=$(echo -n "loki") \
    --from-literal=access_key_id=$(echo -n "$ACCESS_KEY_ID") \
    --from-literal=access_key_secret=$(echo -n "$SECRET_ACCESS_KEY")
}

main() {
  set_credentials_from_aws
  create_secret
}

main
