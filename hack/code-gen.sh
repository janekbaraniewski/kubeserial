#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

: "${GO:=go}"
: "${COPY_OR_DIFF:=copy}"

SCRIPT_ROOT="$(pwd)"
MODULE="github.com/janekbaraniewski/kubeserial"
BOILERPLATE="${SCRIPT_ROOT}/hack/boilerplate.go.txt"

# Resolve the code-generator package from the module graph so the generator
# version always tracks go.mod instead of being hard-coded.
CODEGEN_PKG="$("${GO}" list -m -f '{{.Dir}}' k8s.io/code-generator)"

# shellcheck source=/dev/null
source "${CODEGEN_PKG}/kube_codegen.sh"

printf "Generating deepcopy helpers...\n"
kube::codegen::gen_helpers \
    --boilerplate "${BOILERPLATE}" \
    "${SCRIPT_ROOT}/pkg/apis"

printf "Generating register...\n"
kube::codegen::gen_register \
    --boilerplate "${BOILERPLATE}" \
    "${SCRIPT_ROOT}/pkg/apis"

printf "Generating clientset, listers and informers...\n"
# The input dir is pkg (not pkg/apis): kube_codegen derives the client group
# name from the parent directory of the version package, so pkg/apis/v1alpha1
# yields the "apis" group that existing consumers import
# (pkg/generated/clientset/versioned/typed/apis/v1alpha1).
kube::codegen::gen_client \
    --with-watch \
    --output-dir "${SCRIPT_ROOT}/pkg/generated" \
    --output-pkg "${MODULE}/pkg/generated" \
    --boilerplate "${BOILERPLATE}" \
    "${SCRIPT_ROOT}/pkg"

printf "Code generation complete.\n"

# In CI (diff mode) fail if regeneration produced uncommitted changes.
if [[ "${COPY_OR_DIFF}" == "diff" ]]; then
    printf "Verifying generated code is up to date...\n"
    git diff --exit-code -- pkg/apis pkg/generated
fi
