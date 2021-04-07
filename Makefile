# include the bingo binary variables. This enables the bingo versions to be 
# referenced here as make variables. For example: $(GOLANGCI_LINT)
include .bingo/Variables.mk

# CLUSTER_LOGGING_VERSION
# defines the version of the OpenShift Cluster Logging product.
# Updates this value when a new version of the product should include this operator and its bundle.
CLUSTER_LOGGING_VERSION ?= 5.1.preview.1

# CLUSTER_LOGGING_NS
# defines the default namespace of the OpenShift Cluster Logging product.
CLUSTER_LOGGING_NS ?= openshift-logging

# OLM_VERSION
# defines the version of the operator-lifecycle-manager to deploy
# on local kind clusters. OLM is used to install the operator
# via the operator bundle image.
OLM_VERSION ?= v0.17.0

# VERSION
# defines the project version for the bundle. 
# Update this value when you upgrade the version of your project.
# To re-generate a bundle for another specific version without changing the standard setup, you can:
# - use the VERSION as arg of the bundle target (e.g make bundle VERSION=0.0.2)
# - use environment variables to overwrite this value (e.g export VERSION=0.0.2)
VERSION ?= 0.0.1
CHANNELS ?= "tech-preview"
DEFAULT_CHANNELS ?= "tech-preview"

# CHANNELS define the bundle channels used in the bundle. 
# Add a new line here if you would like to change its default config. (E.g CHANNELS = "preview,fast,stable")
# To re-generate a bundle for other specific channels without changing the standard setup, you can:
# - use the CHANNELS as arg of the bundle target (e.g make bundle CHANNELS=preview,fast,stable)
# - use environment variables to overwrite this value (e.g export CHANNELS="preview,fast,stable")
ifneq ($(origin CHANNELS), undefined)
BUNDLE_CHANNELS := --channels=$(CHANNELS)
endif

# DEFAULT_CHANNEL defines the default channel used in the bundle. 
# Add a new line here if you would like to change its default config. (E.g DEFAULT_CHANNEL = "stable")
# To re-generate a bundle for any other default channel without changing the default setup, you can:
# - use the DEFAULT_CHANNEL as arg of the bundle target (e.g make bundle DEFAULT_CHANNEL=stable)
# - use environment variables to overwrite this value (e.g export DEFAULT_CHANNEL="stable")
ifneq ($(origin DEFAULT_CHANNEL), undefined)
BUNDLE_DEFAULT_CHANNEL := --default-channel=$(DEFAULT_CHANNEL)
endif
BUNDLE_METADATA_OPTS ?= $(BUNDLE_CHANNELS) $(BUNDLE_DEFAULT_CHANNEL)

REGISTRY_ORG ?= openshift-logging

# BUNDLE_IMG defines the image:tag used for the bundle. 
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)
BUNDLE_IMG ?= quay.io/$(REGISTRY_ORG)/loki-operator-bundle:v$(VERSION)

GO_FILES := $(shell find . -type f -name '*.go')

# Image URL to use all building/pushing image targets
IMG ?= quay.io/$(REGISTRY_ORG)/loki-operator:v$(VERSION)

# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true,preserveUnknownFields=false"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: lint generate manager bin/loki-broker

OCI_RUNTIME ?= $(shell which podman || which docker)

# Run tests
ENVTEST_ASSETS_DIR=$(CURDIR)/testbin
test: generate lint manifests
	mkdir -p ${ENVTEST_ASSETS_DIR}
	test -f ${ENVTEST_ASSETS_DIR}/setup-envtest.sh || curl -sSLo ${ENVTEST_ASSETS_DIR}/setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.7.0/hack/setup-envtest.sh
	source ${ENVTEST_ASSETS_DIR}/setup-envtest.sh; fetch_envtest_tools $(ENVTEST_ASSETS_DIR); setup_envtest_env $(ENVTEST_ASSETS_DIR); go test ./... -coverprofile cover.out

# Build manager binary
manager: generate
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate manifests
	go run ./main.go

# Install CRDs into a cluster
install: manifests $(KUSTOMIZE)
	$(KUSTOMIZE) build config/crd | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests $(KUSTOMIZE)
	$(KUSTOMIZE) build config/crd | kubectl delete -f -

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests $(KUSTOMIZE)
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/default | kubectl apply -f -

# UnDeploy controller from the configured Kubernetes cluster in ~/.kube/config
undeploy:
	$(KUSTOMIZE) build config/default | kubectl delete -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests: $(CONTROLLER_GEN)
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases


# Download golangci-lint locally if necessary
GOLANGCI_LINT = $(CURDIR)/bin/golangci-lint
golangci-lint: $(CURDIR)/bin/golangci-lint
$(CURDIR)/bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.38.0

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

fmt: $(GOFUMPT)
	find . -type f -name '*.go' -not -path './vendor/*' -exec $(GOFUMPT) -w {} \;

# Generate code
generate: $(CONTROLLER_GEN)
	go generate ./...
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build the image
oci-build:
	$(OCI_RUNTIME) build -t ${IMG} .

# Push the image
oci-push:
	$(OCI_RUNTIME) push ${IMG}

legacy-manifest-dir: manifests/$(CLUSTER_LOGGING_VERSION)/.keep
manifests/$(CLUSTER_LOGGING_VERSION)/.keep:
	@mkdir -p $(shell dirname $@) && touch $@

# legacy-manifest-files acts as an intermediary target by depending on all
# bundle/manifest files being copied to the manifests directory. The manifest
# yaml files are handled by the target immediately following this, where they
# simply cp source to target using a wildcard. Put simply, legacy-manifest-files
# maps to every yaml under bundle/manifest remapped to it's legacy path e.g.
#   'bundle/manifests/loki.openshift.io_lokistacks.yaml' -> 'manifests/5.1/loki.openshift.io_lokistacks.yaml'
# where manifests/5.1/loki.openshift.io_lokistacks.yaml is the target.
legacy-manifest-files: legacy-manifest-dir
legacy-manifest-files: legacy-csv
legacy-manifest-files: $(filter-out %.clusterserviceversion.yaml, $(patsubst bundle/manifests/%.yaml, manifests/$(CLUSTER_LOGGING_VERSION)/%.yaml, $(wildcard bundle/manifests/*.yaml)))
manifests/$(CLUSTER_LOGGING_VERSION)/%.yaml: bundle/manifests/%.yaml
	@cp -v $^ $@


# The CSV does not work under the previous rule, because the
# CLUSTER_LOGGING_VERSION is part of the target filename so it requires it's
# own target
legacy-csv: manifests/$(CLUSTER_LOGGING_VERSION)/loki-operator.v$(CLUSTER_LOGGING_VERSION).clusterserviceversion.yaml
manifests/$(CLUSTER_LOGGING_VERSION)/loki-operator.v$(CLUSTER_LOGGING_VERSION).clusterserviceversion.yaml: bundle/manifests/loki-operator.clusterserviceversion.yaml
	@cp -v $^ $@

# Generate bundle manifests and metadata, then validate generated files.
.PHONY: bundle
bundle: manifests $(KUSTOMIZE) $(OPERATOR_SDK)
	$(OPERATOR_SDK) generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	$(OPERATOR_SDK) bundle validate ./bundle
	$(MAKE) legacy-manifest-files

# Build the bundle image.
.PHONY: bundle-build
bundle-build:
	$(OCI_RUNTIME) build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

# Create a new local kind cluster.
KIND_CLUSTER_NAME = loki-cluster
.PHONY: kind-create-cluster
kind-create-cluster: $(KIND)
	$(KIND) create cluster --name $(KIND_CLUSTER_NAME)

# Destroy the local kind cluster.
.PHONY: kind-delete-cluster
kind-delete-cluster: $(KIND)
	$(KIND) delete cluster --name $(KIND_CLUSTER_NAME)

# Deploy the operator-lifecycle-manager
# on local kind cluster.
.PHONY: kind-install-olm
kind-install-olm: kind-create-cluster
	./hack/olm-install.sh $(OLM_VERSION)

# Deploy the operator bundle and the operator via OLM
# on local kind cluster.
LOKI_OPERATOR_SECRET_NAME = loki-operator-metrics
.PHONY: kind-olm-deploy
ifneq ($(findstring $(KIND_CLUSTER_NAME),$(shell $(KIND) get clusters)),$(KIND_CLUSTER_NAME))
kind-olm-deploy: kind-install-olm olm-prepare-deploy $(CFSSL) $(CFSSLJSON)
else
kind-olm-deploy: olm-prepare-deploy $(CFSSL) $(CFSSLJSON)
endif
	kubectl apply -f hack/kind/prometheus-operator-0servicemonitorCustomResourceDefinition.yaml
	kubectl apply -f hack/kind/prometheus-operator-0prometheusruleCustomResourceDefinition.yaml
	$(CFSSL) gencert -initca hack/kind/ca.json | $(CFSSLJSON) -bare hack/kind/ca
	$(CFSSL) gencert -ca hack/kind/ca.pem -ca-key hack/kind/ca-key.pem -config hack/kind/cfssl.json -profile=peer hack/kind/loki-operator-monitor-service.json | $(CFSSLJSON) -bare hack/kind/loki-operator-monitor-service
	kubectl get secret -n $(CLUSTER_LOGGING_NS) $(LOKI_OPERATOR_SECRET_NAME) || kubectl create secret generic $(LOKI_OPERATOR_SECRET_NAME) -n $(CLUSTER_LOGGING_NS) --from-file=tls.crt=hack/kind/loki-operator-monitor-service.pem --from-file=tls.key=hack/kind/loki-operator-monitor-service-key.pem
	$(MAKE) olm-deploy REGISTRY_ORG=loki-operator-dev VERSION=$(VERSION)-$(shell whoami)

# Cleanup deployments of the operator bundle and the operator
# via OLM on local kind cluster.
.PHONY: kind-olm-cleanup
kind-olm-cleanup:
	$(MAKE) olm-cleanup

# Build and push the bundle image to a container registry.
.PHONY: olm-deploy-bundle
olm-deploy-bundle: bundle bundle-build
	$(MAKE) oci-push IMG=$(BUNDLE_IMG)

# Build and push the operator image to a container registry.
.PHONY: olm-deploy-operator
olm-deploy-operator: oci-build oci-push

.PHONY: olm-prepare-deploy
olm-prepare-deploy:
	kubectl apply -f hack/kind/namespace.yaml

# Deploy the operator bundle and the operator via OLM into
# an Kubernetes cluster selected via KUBECONFIG.
.PHONY: olm-deploy
ifeq ($(or $(findstring openshift-logging,$(IMG)),$(findstring openshift-logging,$(BUNDLE_IMG))),openshift-logging)
olm-deploy:
	$(error Set variable REGISTRY_ORG to use a custom container registry org account for local development)
else
olm-deploy: olm-deploy-bundle olm-deploy-operator olm-prepare-deploy $(OPERATOR_SDK)
	$(OPERATOR_SDK) run bundle -n $(CLUSTER_LOGGING_NS) $(BUNDLE_IMG)
endif

# Cleanup deployments of the operator bundle and the operator via OLM
# on an OpenShift cluster selected via KUBECONFIG.
.PHONY: olm-cleanup
olm-cleanup: $(OPERATOR_SDK)
	$(OPERATOR_SDK) cleanup -n $(CLUSTER_LOGGING_NS) loki-operator
	kubectl delete -f hack/kind/namespace.yaml

cli: bin/loki-broker
bin/loki-broker: $(GO_FILES) | generate
	go build -o $@ ./cmd/loki-broker/
