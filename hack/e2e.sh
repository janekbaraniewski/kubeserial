#!/usr/bin/env bash
#
# e2e.sh - orchestrate the kubeserial end-to-end test suite against a kind
# cluster. Drives: cluster create -> build+load images -> cert-manager ->
# helm install -> `go test -tags e2e ./test/e2e` -> (always) export logs +
# teardown.
#
# Used by `make test-e2e` and .github/workflows/e2e.yml.
#
# Environment knobs (all optional):
#   E2E_KIND_CLUSTER         kind cluster name        (default: kubeserial-e2e)
#   E2E_NAMESPACE            install namespace        (default: kubeserial)
#   E2E_IMAGE_TAG            local image tag          (default: local)
#   E2E_NODE_IMAGE           kindest/node image       (default: pinned v1.27.1)
#   E2E_CERT_MANAGER_VERSION cert-manager version     (default: v1.14.5)
#   E2E_SKIP_CLUSTER_SETUP   reuse existing cluster   (default: false)
#   E2E_SKIP_TEARDOWN        keep cluster after run   (default: false)
#   E2E_SKIP_DEVICE_SIM      skip E4/E5 device specs  (default: false)
#   E2E_KUBECONTEXT          kube context to target   (default: kind-<cluster>)
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${REPO_ROOT}"

KIND_CLUSTER="${E2E_KIND_CLUSTER:-kubeserial-e2e}"
NAMESPACE="${E2E_NAMESPACE:-kubeserial}"
IMAGE_TAG="${E2E_IMAGE_TAG:-local}"
SKIP_CLUSTER_SETUP="${E2E_SKIP_CLUSTER_SETUP:-false}"
SKIP_TEARDOWN="${E2E_SKIP_TEARDOWN:-false}"
SKIP_DEVICE_SIM="${E2E_SKIP_DEVICE_SIM:-false}"
KUBECONTEXT="${E2E_KUBECONTEXT:-kind-${KIND_CLUSTER}}"
KIND_CONFIG="${REPO_ROOT}/test/e2e/kind-config.yaml"
CERT_MANAGER_VERSION="${E2E_CERT_MANAGER_VERSION:-v1.14.5}"
SIMULATOR_IMG="alpine/socat:latest"

# Pinned node image compatible with kind v0.18 (newest it ships). Set
# E2E_NODE_IMAGE to a different tag to override, or to the empty string to let
# kind pick its own default (use this with newer kind in CI). The variable is
# only defaulted when entirely unset, so an explicit empty value is honored.
NODE_IMAGE="${E2E_NODE_IMAGE-kindest/node:v1.27.1@sha256:9915f5629ef4d29f35b478e819249e89cfaffcbfeebda4324e5c01d53d937b09}"

CTRL_IMG="ghcr.io/janekbaraniewski/kubeserial"
MONITOR_IMG="ghcr.io/janekbaraniewski/kubeserial-device-monitor"
WEBHOOK_IMG="ghcr.io/janekbaraniewski/kubeserial-injector-webhook"

log() { echo -e "\n\033[1;34m==> $*\033[0m"; }

require() {
  command -v "$1" >/dev/null 2>&1 || { echo "ERROR: '$1' not found on PATH" >&2; exit 1; }
}

teardown() {
  local rc=$?
  if [[ "${SKIP_TEARDOWN}" != "true" && "${SKIP_CLUSTER_SETUP}" != "true" ]]; then
    log "Exporting cluster logs (best-effort) before teardown"
    kind export logs "/tmp/kubeserial-e2e-logs" --name "${KIND_CLUSTER}" || true
    log "Deleting kind cluster ${KIND_CLUSTER}"
    kind delete cluster --name "${KIND_CLUSTER}" || true
  fi
  exit "${rc}"
}
trap teardown EXIT

require go
require kubectl
require helm

if [[ "${SKIP_CLUSTER_SETUP}" != "true" ]]; then
  require kind
  require docker

  log "Creating kind cluster ${KIND_CLUSTER} (node ${NODE_IMAGE:-kind default})"
  if kind get clusters | grep -qx "${KIND_CLUSTER}"; then
    echo "cluster already exists, reusing"
  else
    mkdir -p /tmp/kubeserial-e2e-dev
    image_args=()
    [[ -n "${NODE_IMAGE}" ]] && image_args=(--image "${NODE_IMAGE}")
    kind create cluster --name "${KIND_CLUSTER}" "${image_args[@]}" \
      --config "${KIND_CONFIG}" --wait 180s
  fi

  # Build images directly with buildx. NOTE: the Makefile's monitor/webhook
  # docker targets hardcode a registry build-cache that needs auth and fails
  # offline, so we invoke buildx here instead of `make docker-local`.
  log "Building images (tag ${IMAGE_TAG})"
  docker buildx build . -f Dockerfile         --platform linux/$(go env GOARCH) -t "${CTRL_IMG}:${IMAGE_TAG}"    --load
  docker buildx build . -f Dockerfile.monitor --platform linux/$(go env GOARCH) -t "${MONITOR_IMG}:${IMAGE_TAG}" --load
  docker buildx build . -f Dockerfile.webhook --platform linux/$(go env GOARCH) -t "${WEBHOOK_IMG}:${IMAGE_TAG}" --load

  log "Loading images into kind"
  kind load docker-image "${CTRL_IMG}:${IMAGE_TAG}"    --name "${KIND_CLUSTER}"
  kind load docker-image "${MONITOR_IMG}:${IMAGE_TAG}" --name "${KIND_CLUSTER}"
  kind load docker-image "${WEBHOOK_IMG}:${IMAGE_TAG}" --name "${KIND_CLUSTER}"

  # Pre-load the device-simulator image so the E4/E5 simulator pod starts
  # without needing registry access from inside the cluster.
  if [[ "${SKIP_DEVICE_SIM}" != "true" ]]; then
    log "Pre-loading device-simulator image ${SIMULATOR_IMG}"
    docker pull "${SIMULATOR_IMG}"
    kind load docker-image "${SIMULATOR_IMG}" --name "${KIND_CLUSTER}"
  fi

  log "Installing cert-manager ${CERT_MANAGER_VERSION} (webhook certs depend on it)"
  kubectl --context "${KUBECONTEXT}" apply -f \
    "https://github.com/cert-manager/cert-manager/releases/download/${CERT_MANAGER_VERSION}/cert-manager.yaml"
  kubectl --context "${KUBECONTEXT}" -n cert-manager wait --for=condition=Available deploy --all --timeout=180s

  log "Creating namespace and self-signed Issuer"
  kubectl --context "${KUBECONTEXT}" create namespace "${NAMESPACE}" --dry-run=client -o yaml \
    | kubectl --context "${KUBECONTEXT}" apply -f -
  cat <<EOF | kubectl --context "${KUBECONTEXT}" apply -f -
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: kubeserial-selfsigned
  namespace: ${NAMESPACE}
spec:
  selfSigned: {}
EOF

  log "Installing CRDs chart"
  helm --kube-context "${KUBECONTEXT}" upgrade --install kubeserial-crds charts/kubeserial-crds \
    --namespace "${NAMESPACE}" --wait

  log "Installing kubeserial chart (images tag=${IMAGE_TAG}, pullPolicy=Never)"
  helm --kube-context "${KUBECONTEXT}" upgrade --install kubeserial charts/kubeserial \
    --namespace "${NAMESPACE}" \
    --set image.tag="${IMAGE_TAG}"          --set image.pullPolicy=Never \
    --set monitor.image.tag="${IMAGE_TAG}"  --set monitor.image.pullPolicy=Never \
    --set webhook.image.tag="${IMAGE_TAG}"  --set webhook.image.pullPolicy=Never \
    --set certManagerIssuer.name=kubeserial-selfsigned --set certManagerIssuer.kind=Issuer \
    --wait --timeout 5m
fi

log "Running e2e specs"
E2E_SKIP_CLUSTER_SETUP=true \
E2E_KUBECONTEXT="${KUBECONTEXT}" \
E2E_NAMESPACE="${NAMESPACE}" \
E2E_IMAGE_TAG="${IMAGE_TAG}" \
E2E_KIND_CLUSTER="${KIND_CLUSTER}" \
E2E_SKIP_DEVICE_SIM="${SKIP_DEVICE_SIM}" \
  go test -tags e2e ./test/e2e/... -v -timeout 20m

log "e2e suite passed"
