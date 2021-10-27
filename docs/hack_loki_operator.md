# Loki Operator Hack

Loki Operator is the Kubernetes Operator for [Loki](https://grafana.com/docs/loki/latest/) provided by the Red Hat OpenShift engineering team.

## Hacking on Loki Operator using kind

[kind](https://kind.sigs.k8s.io/docs/user/quick-start/) is a tool for running local Kubernetes clusters using Docker container "nodes". kind was primarily designed for testing Kubernetes itself, but may be used for local development or CI.

### Requirements

* Install [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) or [Openshift CLI](https://docs.openshift.com/container-platform/latest/cli_reference/openshift_cli/getting-started-cli.html) for communicating with the cluster. The guide below will be using `kubectl` for the same.
* Create a running Kubernetes cluster using kind.
* A container registry that you and your Kubernetes cluster can reach. We recommend  [quay.io](https://quay.io/signin/).

### Installation of Loki Operator

* Build and push the container image and then deploy the operator with:

  ```console
  make oci-build oci-push deploy REGISTRY_ORG=$YOUR_QUAY_ORG VERSION=latest
  ```

  where `$YOUR_QUAY_ORG` is your personal [quay.io](http://quay.io/) account where you can push container images.

  The above command will deploy the operator to your active Kubernetes cluster defined by your local [kubeconfig](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/). The operator will be running in the `default` namespace.

* You can confirm that the operator is up and running using:

  ```console
  kubectl get pods
  ```

* Now create a LokiStack instance to get the various components of Loki up and running:

  ```console
  kubectl apply -f hack/lokistack_gateway_dev.yaml
  ```

  This will create `distributor`, `compactor`, `ingester`, `querier`, `query-frontend` and `lokistack-gateway` components.

  Confirm that all are up and running using:

  ```console
  kubectl get pods
  ```

  _Note:_  `lokistack-gateway` is an optional component deployed as part of Loki Operator. It provides secure access to Loki's distributor (i.e. for pushing logs) and query-frontend (i.e. for querying logs) via consulting an OAuth/OIDC endpoint for the request subject.

  If you don't want `lokistack-gateway` component then you can skip it by removing the `--with-lokistack-gateway` args from the `controller-manager` deployment:

  ```console
  kubectl edit deployment/controller-manager
  ```

  Delete the `args` part in it and save the file. This will update the deployment and now you can create LokiStack instance using:

  ```console
  kubectl apply -f hack/lokistack_dev.yaml
  ```

  This will create `distributor`, `compactor`, `ingester`, `querier` and `query-frontend` components only.

## Hacking on Loki Operator on OpenShift

### Requirements

* Install [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) or [Openshift CLI](https://docs.openshift.com/container-platform/latest/cli_reference/openshift_cli/getting-started-cli.html) for communicating with the cluster. The guide below will be using `kubectl` for the same.
* Create a running OpenShift cluster.
* A container registry that you and your OpenShift cluster can reach. We recommend  [quay.io](https://quay.io/signin/).

### Installation of Loki Operator

* Build and push the container image and then deploy the operator with:

  ```console
  make olm-deploy REGISTRY_ORG=$YOUR_QUAY_ORG VERSION=$VERSION
  ```

  where `$YOUR_QUAY_ORG` is your personal [quay.io](http://quay.io/) account where you can push container images and `$VERSION` can be any random version number such as `v0.0.1`.

  The above command will deploy the operator to your active Openshift cluster defined by your local [kubeconfig](https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/). The operator will be running in the `openshift-logging` namespace.

* You can confirm that the operator is up and running using:

  ```console
  kubectl -n openshift-logging get pods
  ```

* Now you need to create a storage secret for the operator. This can be done using:

  ```
  make olm-deploy-example-storage-secret
  ```

  OR

  ```console
  ./hack/deploy-example-secret.sh openshift-logging
  ```

  This secret will be available in openshift-logging namespace. You can check the `hack/deploy-example-secret.sh` file to check the content of the secret.

* Once the object storage secret is created, you can now create a LokiStack instance to get the various components of Loki up and running:

  ```console
  kubectl apply -f hack/lokistack_gateway_dev.yaml
  ```

  This will create `distributor`, `compactor`, `ingester`, `querier`, `query-frontend` and `lokistack-gateway` components.

  Confirm that all are up and running using:

  ```console
  kubectl get pods
  ```

  _Note:_  `lokistack-gateway` is an optional component deployed as part of Loki Operator. It provides secure access to Loki's distributor (i.e. for pushing logs) and query-frontend (i.e. for querying logs) via consulting an OAuth/OIDC endpoint for the request subject.

  If you don't want `lokistack-gateway` component then you can skip it by removing the `--with-lokistack-gateway` args from the `loki-operator-controller-manager` deployment:

  ```console
  kubectl edit deployment/loki-operator-controller-manager
  ```

  Delete the `args` part in it and save the file. This will update the deployment and now you can create LokiStack instance using:

  ```console
  kubectl apply -f hack/lokistack_dev.yaml
  ```

  This will create `distributor`, `compactor`, `ingester`, `querier` and `query-frontend` components only.

## Basic Troubleshooting on Hacking on Loki Operator

### New changes are not detected by Loki Operator

Suppose you made some changes to the Loki Operator's code and deployed it but the changes are not visible when it runs. This happens when the deployment pulls the old image of the operator because of the `imagePullPolicy` being set to `IfNotPresent`. There, you need to make some changes to make your deployment pull new image always:

* Go to `config/manager/manager.yaml` file.
* Set the `imagePullPolicy` to `Always` i.e.,

  ```yaml
  imagePullPolicy: Always
  ```

* Deploy the operator again.

### kubectl using old context

It is possible that when you use two different clusters - one is kind cluster and the other is OpenShift cluster, you might need to switch between clusters to test your changes. There is a possibility that once you switch between clusters, the kubectl might not switch the context automatically and hence you might need to do this manually to correctly communicate with your cluster.

* List all the available context:

  ```console
  kubectl config get-contexts
  ```

  The `*` mark against the context shows the one in use currently.
* Set the context name you want to use now:

  ```console
  kubectl config use-context $CONTEXTNAME
  ```

  where `$CONTEXTNAME` is the context name you want to use now from the previous step.
