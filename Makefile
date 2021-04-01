# include the bingo binary variables. This enables the bingo versions to be 
# referenced here as make variables. For example: $(GOLANGCI_LINT)
include .bingo/Variables.mk

# CLUSTER_LOGGING_VERSION
# defines the version of the OpenShift Cluster Logging product.
# Updates this value when a new version of the product should include this operator and its bundle.
CLUSTER_LOGGING_VERSION ?= 5.1.preview.1

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

# BUNDLE_IMG defines the image:tag used for the bundle. 
# You can use it as an arg. (E.g make bundle-build BUNDLE_IMG=<some-registry>/<project-name-bundle>:<tag>)
BUNDLE_IMG ?= quay.io/openshift-logging/loki-operator-bundle:$(VERSION)

GO_FILES := $(shell find . -type f -name '*.go')

# Image URL to use all building/pushing image targets
IMG ?= quay.io/openshift-logging/loki-operator:latest
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
bundle: manifests legacy-manifest-dir $(KUSTOMIZE) $(OPERATOR_SDK)
	$(OPERATOR_SDK) generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(IMG)
	$(KUSTOMIZE) build config/manifests | $(OPERATOR_SDK) generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	$(OPERATOR_SDK) bundle validate ./bundle
	$(MAKE) legacy-manifest-files

# Build the bundle image.
.PHONY: bundle-build
bundle-build:
	$(OCI_RUNTIME) build -f bundle.Dockerfile -t $(BUNDLE_IMG) .

cli: bin/loki-broker
bin/loki-broker: $(GO_FILES) | generate
	go build -o $@ ./cmd/loki-broker/
