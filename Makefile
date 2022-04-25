KUBESERIAL_REGISTRY=janekbaraniewski/kubeserial
DEVICE_MONITOR_REGISTRY=janekbaraniewski/kubeserial-device-monitor
TARGET_PLATFORMS=$(shell cat TARGET_PLATFORMS)
VERSION ?= 0.0.1-$(shell git rev-parse --short HEAD)
DOCKERBUILD_EXTRA_OPTS ?=
DOCKERBUILD_PLATFORM_OPT=--platform
KUBESERIAL_BUILD_OUTPUT_PATH ?= build/_output/bin/kubeserial
DEVICE_MONITOR_BUILD_OUTPUT_PATH ?= build/_output/bin/device-monitor
RELEASE_NAME ?= kubeserial
ENVTEST_K8S_VERSION = 1.23

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: manifests-gen
manifests-gen: COPY_OR_DIFF=copy
manifests-gen: manifests-gen-script

.PHONY: check-manifests-gen
check-manifests-gen: COPY_OR_DIFF=diff
check-manifests-gen: manifests-gen-script

manifests-gen-script:
	@COPY_OR_DIFF=${COPY_OR_DIFF} ./hack/manifests-gen.sh

.PHONY: code-gen
code-gen: COPY_OR_DIFF=copy
code-gen: code-gen-script

.PHONY: check-code-gen
check-code-gen: COPY_OR_DIFF=diff
check-code-gen: code-gen-script

code-gen-script:
	@COPY_OR_DIFF=${COPY_OR_DIFF} ./hack/code-gen.sh

.PHONY: generate
generate: manifests-gen code-gen

.PHONY: check-generated
check-generated: check-code-gen check-manifests-gen

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./... -coverprofile coverage.txt -covermode atomic

# ENVTEST = $(shell pwd)/bin/setup-envtest
# .PHONY: envtest
# envtest: ## Download envtest-setup locally if necessary.
# 	$(call go-get-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

##@ Build

.PHONY: all
all: generate kubeserial ## Run codegen and build all components.

PHONY: .kubeserial
kubeserial: ## Build manager binary.
	go build -o ${KUBESERIAL_BUILD_OUTPUT_PATH} cmd/manager/main.go

PHONY: .device-monitor
device-monitor: ## Build device monitor binary
	go build -o ${DEVICE_MONITOR_BUILD_OUTPUT_PATH} cmd/device-monitor/main.go

.PHONY: run
run: generate fmt vet ## Run codegen and start controller from your host.
	go run ./cmd/manager/main.go

##@ Docker

PHONY: .kubeserial-docker-local
kubeserial-docker-local: PLATFORMS=
kubeserial-docker-local: DOCKERBUILD_PLATFORM_OPT=
kubeserial-docker-local: DOCKERBUILD_ACTION=--load
kubeserial-docker-local: VERSION ?= local
kubeserial-docker-local: kubeserial-docker ## Build image for local development, tag local, supports only builder platform

PHONY: .kubeserial-docker-all
kubeserial-docker-all: PLATFORMS=${TARGET_PLATFORMS}
kubeserial-docker-all: DOCKERBUILD_ACTION=--push
kubeserial-docker-all: kubeserial-docker ## Build and push image for all target platforms

PHONY: .kubeserial-docker
kubeserial-docker: DOCKERFILE=Dockerfile
kubeserial-docker: REGISTRY=${KUBESERIAL_REGISTRY}
kubeserial-docker: docker-build

PHONY: .device-monitor-docker-local
device-monitor-docker-local: PLATFORMS=
device-monitor-docker-local: DOCKERBUILD_PLATFORM_OPT=
device-monitor-docker-local: DOCKERBUILD_ACTION=--load
device-monitor-docker-local: VERSION ?= local
device-monitor-docker-local: device-monitor-docker ## Build image for local development, tag local, supports only builder platform

PHONY: .device-monitor-docker-all
device-monitor-docker-all: PLATFORMS=$(TARGET_PLATFORMS)
device-monitor-docker-all: DOCKERBUILD_ACTION=--push
device-monitor-docker-all: device-monitor-docker ## Build and push image for all target platforms

PHONY: .device-monitor-docker
device-monitor-docker: DOCKERFILE=Dockerfile.monitor
device-monitor-docker: REGISTRY=${DEVICE_MONITOR_REGISTRY}
device-monitor-docker: docker-build

PHONY: .docker-build
docker-build:
	docker buildx build . -f ${DOCKERFILE} ${DOCKERBUILD_EXTRA_OPTS} ${DOCKERBUILD_PLATFORM_OPT} ${PLATFORMS} -t $(REGISTRY):$(VERSION) ${DOCKERBUILD_ACTION}

##@ Helm

PHONY: .update-kubeserial-chart-version
update-kubeserial-chart-version: CHART_PATH=./deploy/chart/kubeserial
update-kubeserial-chart-version: ## Update version used in chart. Requires VERSION var to be set
	@CHART_PATH=${CHART_PATH} VERSION=${VERSION} ./hack/update-chart-version.sh

PHONY: .update-kubeserial-crds-chart-version
update-kubeserial-crds-chart-version: CHART_PATH=./deploy/chart/kubeserial-crds
update-kubeserial-crds-chart-version: ## Update version used in chart. Requires VERSION var to be set
	@CHART_PATH=${CHART_PATH} VERSION=${VERSION} ./hack/update-chart-version.sh

PHONY: .helm-lint
helm-lint: ## Run chart-testing to lint kubeserial chart.
	@ct lint --chart-dirs deploy/chart/ --check-version-increment=false

PHONY: .update-crds-labels
update-crds-labels:
	@python3 ./hack/update-crd-metadata.py deploy/chart/kubeserial-crds/templates/app.kubeserial.com_kubeserials.yaml hack/crd_metadata_template.yaml

##@ Deployment

.PHONY: uninstall
uninstall: ## Uninstall release.
	helm uninstall ${RELEASE_NAME}

.PHONY: .deploy-dev
deploy-dev: manifests-gen update-kubeserial-chart-version update-kubeserial-crds-chart-version ## Install dev release in current context/namespace.
	helm upgrade --install ${RELEASE_NAME}-crds ./deploy/chart/kubeserial-crds
	helm upgrade --install ${RELEASE_NAME} ./deploy/chart/kubeserial


# # go-get-tool will 'go get' any package $2 and install it to $1.
# PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
# define go-get-tool
# @[ -f $(1) ] || { \
# set -e ;\
# TMP_DIR=$$(mktemp -d) ;\
# cd $$TMP_DIR ;\
# go mod init tmp ;\
# echo "Downloading $(2)" ;\
# GOBIN=$(PROJECT_DIR)/bin go get $(2) ;\
# rm -rf $$TMP_DIR ;\
# }
# endef
