DOCKERHUB=janekbaraniewski/kubeserial
TARGET_PLATFORMS=$(shell cat TARGET_PLATFORMS)
VERSION ?= $(shell git rev-parse --short HEAD)
DOCKERBUILD_EXTRA_OPTS ?=
DOCKERBUILD_PLATFORM_OPT=--platform
GO_BUILD_OUTPUT_PATH ?= build/_output/bin/kubeserial
RELEASE_NAME ?= kubeserial

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
test: manifests generate fmt vet envtest ## Run tests.
	KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path)" go test ./... -coverprofile cover.out

##@ Build

.PHONY: all
all: generate kubeserial ## Run codegen and build all components.

PHONY: .kubeserial
kubeserial: ## Build manager binary.
	go build -o ${GO_BUILD_OUTPUT_PATH} cmd/manager/main.go

.PHONY: run
run: generate fmt vet ## Run codegen and start controller from your host.
	go run ./cmd/manager/main.go

##@ Docker

PHONY: .kubeserial-docker-local
kubeserial-docker-local: PLATFORMS=
kubeserial-docker-local: DOCKERBUILD_PLATFORM_OPT=
kubeserial-docker-local: DOCKERBUILD_ACTION=--load
kubeserial-docker-local: VERSION=local
kubeserial-docker-local: kubeserial-docker ## Build image for local development, tag local, supports only builder platform

PHONY: .kubeserial-docker-all
kubeserial-docker-all: PLATFORMS=${TARGET_PLATFORMS}
kubeserial-docker-all: DOCKERBUILD_ACTION=--push
kubeserial-docker-all: kubeserial-docker ## Build and push image for all target platforms

PHONY: .kubeserial-docker
kubeserial-docker: DOCKERFILE=Dockerfile
kubeserial-docker: docker-build

PHONY: .docker-build
docker-build:
	docker buildx build . -f ${DOCKERFILE} ${DOCKERBUILD_EXTRA_OPTS} ${DOCKERBUILD_PLATFORM_OPT} ${PLATFORMS} -t $(DOCKERHUB):$(VERSION) ${DOCKERBUILD_ACTION}

##@ Helm

PHONY: .update-chart-version
update-chart-version: ## Update version used in chart. Requires VERSION var to be set
	@./hack/update-chart-version.sh

PHONY: .helm-lint
helm-lint: ## Run chart-testing to lint kubeserial chart.
	@ct lint --charts deploy/chart/kubeserial

PHONY: .update-crds-labels
update-crds-labels:
	@python3 ./hack/update-crd-metadata.py deploy/chart/kubeserial-crds/templates/app.kubeserial.com_kubeserials.yaml hack/crd_metadata_template.yaml

##@ Deployment

.PHONY: uninstall
uninstall: ## Uninstall release.
	helm uninstall ${RELEASE_NAME}

.PHONY: deploy
deploy: manifests update-chart-version ## Install release in current context/namespace.
	helm upgrade --install ${RELEASE_NAME} ${CHART_PATH}
