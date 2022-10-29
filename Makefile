KUBESERIAL_REGISTRY=ghcr.io/janekbaraniewski/kubeserial
DEVICE_MONITOR_REGISTRY=ghcr.io/janekbaraniewski/kubeserial-device-monitor
INJECTOR_WEBHOOK_REGISTRY=ghcr.io/janekbaraniewski/kubeserial-injector-webhook
TARGET_PLATFORMS=$(shell cat TARGET_PLATFORMS)
APP_VERSION?=$(shell git describe --dirty --tags --match "[0-9]*" )
CRDS_VERSION?=$(shell git describe --dirty --tags --match "crds*" | sed 's/crds-//')
DOCKERBUILD_PLATFORM_OPT=--platform
RELEASE_NAME ?= kubeserial
ENVTEST_K8S_VERSION = 1.23.3
MINIKUBE_PROFILE=kubeserial
ARM_PLATFORM=linux/arm64

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

include Makefile.build

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

.PHONY: manifests-gen-script
manifests-gen-script:
	@COPY_OR_DIFF=${COPY_OR_DIFF} ./hack/manifests-gen.sh

.PHONY: code-gen
code-gen: COPY_OR_DIFF=copy
code-gen: code-gen-script

.PHONY: check-code-gen
check-code-gen: COPY_OR_DIFF=diff
check-code-gen: code-gen-script

.PHONY: code-gen-script
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

.PHONY: get-test-assets
get-test-assets:
	@echo "Genereting test assets"
	@./hack/get-test-assets.sh

.PHONY: test
test: fmt vet envtest-render-crds get-test-assets ## Run tests.
	go test ./... -coverprofile coverage.txt.tmp -covermode atomic
	@cat coverage.txt.tmp | grep -v "fake_api.go" > coverage.txt
	@rm coverage.txt.tmp
	@rm -r test-assets

.PHONY: test-fswatch
test-fswatch: ## Use fswatch to watch source files and run tests on chamnge
	fswatch -or pkg Makefile Makefile.build cmd go.mod go.sum | xargs -n1 -I{} make test

.PHONY: envtest-render-crds
envtest-render-crds:
	@rm -rf build/_output/kubeserial-crds || echo ""
	@helm template charts/kubeserial-crds --name-template kubeserial --output-dir build/_output

.PHONY: check
check: check-generated
	golangci-lint run

# ENVTEST = $(shell pwd)/bin/setup-envtest
# .PHONY: envtest
# envtest: ## Download envtest-setup locally if necessary.
# 	$(call go-get-tool,$(ENVTEST),sigs.k8s.io/controller-runtime/tools/setup-envtest@latest)

##@ Run

.PHONY: run
run: generate fmt vet ## Run codegen and start controller from your host.
	go run ./cmd/manager/main.go

##@ Docker

.PHONY: docker-local
docker-local: kubeserial-docker-local device-monitor-docker-local injector-webhook-docker-local

.PHONY: docker-all
docker-all: kubeserial-docker-all device-monitor-docker-all injector-webhook-docker-all

.PHONY: kubeserial-docker-local
kubeserial-docker-local: PLATFORMS?=$(ARM_PLATFORM)
kubeserial-docker-local: DOCKERBUILD_PLATFORM_OPT?=
kubeserial-docker-local: DOCKERBUILD_ACTION?=--load
kubeserial-docker-local: APP_VERSION ?= local
kubeserial-docker-local: kubeserial-docker ## Build image for local development, tag local, supports only builder platform

.PHONY: kubeserial-docker-all
kubeserial-docker-all: PLATFORMS=${TARGET_PLATFORMS}
kubeserial-docker-all: DOCKERBUILD_ACTION?=--push
kubeserial-docker-all: kubeserial-docker ## Build and push image for all target platforms

.PHONY: kubeserial-docker
kubeserial-docker: DOCKERFILE=Dockerfile
kubeserial-docker: REGISTRY=${KUBESERIAL_REGISTRY}
kubeserial-docker: DOCKERBUILD_EXTRA_OPTS?=--cache-to=type=registry,mode=max,ref=${KUBESERIAL_REGISTRY}:cache --cache-from ${KUBESERIAL_REGISTRY}:cache
kubeserial-docker:
	docker buildx build . -f ${DOCKERFILE} ${DOCKERBUILD_EXTRA_OPTS} ${DOCKERBUILD_PLATFORM_OPT} ${PLATFORMS} -t $(REGISTRY):$(APP_VERSION) ${DOCKERBUILD_ACTION}

.PHONY: device-monitor-docker-local
device-monitor-docker-local: PLATFORMS?=$(ARM_PLATFORM)
device-monitor-docker-local: DOCKERBUILD_PLATFORM_OPT?=
device-monitor-docker-local: DOCKERBUILD_ACTION?=--load
device-monitor-docker-local: APP_VERSION ?= local
device-monitor-docker-local: device-monitor-docker ## Build image for local development, tag local, supports only builder platform

.PHONY: device-monitor-docker-all
device-monitor-docker-all: PLATFORMS=$(TARGET_PLATFORMS)
device-monitor-docker-all: DOCKERBUILD_ACTION?=--push
device-monitor-docker-all: device-monitor-docker ## Build and push image for all target platforms

.PHONY: device-monitor-docker
device-monitor-docker: DOCKERFILE=Dockerfile.monitor
device-monitor-docker: REGISTRY=${DEVICE_MONITOR_REGISTRY}
device-monitor-docker: DOCKERBUILD_EXTRA_OPTS=--cache-to=type=registry,mode=max,ref=${DEVICE_MONITOR_REGISTRY}:cache --cache-from ${DEVICE_MONITOR_REGISTRY}:cache
device-monitor-docker:
	docker buildx build . -f ${DOCKERFILE} ${DOCKERBUILD_EXTRA_OPTS} ${DOCKERBUILD_PLATFORM_OPT} ${PLATFORMS} -t $(REGISTRY):$(APP_VERSION) ${DOCKERBUILD_ACTION}

.PHONY: injector-webhook-docker-local
injector-webhook-docker-local: PLATFORMS?=$(ARM_PLATFORM)
injector-webhook-docker-local: DOCKERBUILD_PLATFORM_OPT?=
injector-webhook-docker-local: DOCKERBUILD_ACTION?=--load
injector-webhook-docker-local: APP_VERSION ?= local
injector-webhook-docker-local: injector-webhook-docker ## Build image for local development, tag local, supports only builder platform

.PHONY: injector-webhook-docker-all
injector-webhook-docker-all: PLATFORMS=$(TARGET_PLATFORMS)
injector-webhook-docker-all: DOCKERBUILD_ACTION?=--push
injector-webhook-docker-all: injector-webhook-docker ## Build and push image for all target platforms

.PHONY: injector-webhook-docker
injector-webhook-docker: DOCKERFILE=Dockerfile.webhook
injector-webhook-docker: REGISTRY=${INJECTOR_WEBHOOK_REGISTRY}
injector-webhook-docker: DOCKERBUILD_EXTRA_OPTS=--cache-to=type=registry,mode=max,ref=${INJECTOR_WEBHOOK_REGISTRY}:cache --cache-from ${INJECTOR_WEBHOOK_REGISTRY}:cache
injector-webhook-docker:
	docker buildx build . -f ${DOCKERFILE} ${DOCKERBUILD_EXTRA_OPTS} ${DOCKERBUILD_PLATFORM_OPT} ${PLATFORMS} -t $(REGISTRY):$(APP_VERSION) ${DOCKERBUILD_ACTION}

##@ Helm

.PHONY: update-kubeserial-chart-version
update-kubeserial-chart-version: CHART_PATH=./charts/kubeserial
update-kubeserial-chart-version:
	@CHART_PATH=${CHART_PATH} VERSION=${APP_VERSION} ./hack/update-chart-version.sh

.PHONY: update-kubeserial-crds-chart-version
update-kubeserial-crds-chart-version: CHART_PATH=./charts/kubeserial-crds
update-kubeserial-crds-chart-version:
	@CHART_PATH=${CHART_PATH} VERSION=${CRDS_VERSION} ./hack/update-chart-version.sh

.PHONY: helm-lint
helm-lint: ## Run chart-testing to lint kubeserial chart.
	@ct lint --chart-dirs charts/ --check-version-increment=false

.PHONY: update-crds-labels
update-crds-labels:
	@python3 ./hack/update-crd-metadata.py charts/kubeserial-crds/templates/app.kubeserial.com_kubeserials.yaml hack/crd_metadata_template.yaml

.PHONY: update-webhook-template
update-webhook-template:
	@python3 ./hack/update-webhook-template.py charts/kubeserial/templates/webhooks.yaml hack/webhook_template.yaml

.PHONY: update-version
update-version: ## Update charts version.
update-version: update-kubeserial-crds-chart-version update-kubeserial-chart-version

##@ Kind

.PHONY: kind
kind: kind-create install-certmanager docker-local kind-load-images install-dev ## Create kind cluster, install certmanager, build and load images, install dev release.

kind-create:
	kind create cluster --name kubeserial

.PHONY: install-certmanager
install-certmanager:
	helm upgrade --install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.10.0 --set installCRDs=true

kind-load-images:
	kind load docker-image --name kubeserial ghcr.io/janekbaraniewski/kubeserial:${APP_VERSION}
	kind load docker-image --name kubeserial ghcr.io/janekbaraniewski/kubeserial-device-monitor:${APP_VERSION}
	kind load docker-image --name kubeserial ghcr.io/janekbaraniewski/kubeserial-injector-webhook:${APP_VERSION}

##@ Minikube

.PHONY: minikube
minikube: minikube-start minikube-build-controller-image minikube-build-monitor-image install-certmanager install-dev ## Start local cluster, build image and deploy

.PHONY: minikube-start
minikube-start: ## Start minikube cluster
	@minikube -p ${MINIKUBE_PROFILE}  start

.PHONY: minikube-set-context
minikube-set-context: ## Set context to use minikube cluster
	@minikube -p ${MINIKUBE_PROFILE} update-context

.PHONY: minikube-build-controller-image
minikube-build-controller-image: DOCKERFILE=Dockerfile
minikube-build-controller-image: REGISTRY=${KUBESERIAL_REGISTRY}
minikube-build-controller-image: minikube-build-image

.PHONY: minikube-build-monitor-image
minikube-build-monitor-image: DOCKERFILE=Dockerfile.monitor
minikube-build-monitor-image: REGISTRY=${DEVICE_MONITOR_REGISTRY}
minikube-build-monitor-image: minikube-build-image

.PHONY: minikube-build-image
minikube-build-image: DOCKERBUILD_EXTRA_OPTS=--load
minikube-build-image:
	@eval $$(minikube -p ${MINIKUBE_PROFILE} docker-env) ;\
	docker buildx build . -f ${DOCKERFILE} ${DOCKERBUILD_EXTRA_OPTS} -t $(REGISTRY):$(APP_VERSION)
	@echo "Finished building image ${REGISTRY}:${APP_VERSION}"
	@echo "Available images:"
	@eval $$(minikube -p ${MINIKUBE_PROFILE} docker-env) ;\
	docker images

##@ Deployment

.PHONY: uninstall
uninstall: ## Uninstall release.
	helm uninstall ${RELEASE_NAME}

.PHONY: deploy-dev
deploy-dev: manifests-gen update-version
	helm upgrade --create-namespace --namespace kubeserial --install ${RELEASE_NAME}-crds ./charts/kubeserial-crds
	helm upgrade --create-namespace --namespace kubeserial --install ${RELEASE_NAME} ./charts/kubeserial -f ./charts/kubeserial/values-local.yaml

.PHONY: install-dev
install-dev: update-version deploy-dev ## Install dev release in current context/namespace.

##@ Docs

.PHONY: docs-install-deps
docs-deps: ## Install mdbook (requires rust and cargo) + plugins
	cargo install mdbook mdbook-mermaid mdbook-open-on-gh mdbook-toc

.PHONY: serve-docs
docs-serve: ## Build docs, start server and open in browser
	cd docs && mdbook serve --open
