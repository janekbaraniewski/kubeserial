#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

: ${GO=go}

export GOPATH=$("$GO" env GOPATH | awk -F ':' '{print $1}')
export PATH=$PATH:$GOPATH/bin

source hack/utils.sh

printf "Start code-gen script \n"

# Install the required binaries with modules enabled
go install k8s.io/code-generator/cmd/client-gen@v0.35.1
go install k8s.io/code-generator/cmd/lister-gen@v0.35.1
go install k8s.io/code-generator/cmd/informer-gen@v0.35.1

# Source kube_codegen.sh for helper functions
CODEGEN_PKG_PATH=$(go env GOMODCACHE)/k8s.io/code-generator@v0.35.1
chmod +x "${CODEGEN_PKG_PATH}/kube_codegen.sh"
source "${CODEGEN_PKG_PATH}/kube_codegen.sh"

printf "Got all dependencies \n"

printf "Running code generators...\n"

TEMP_OUTPUT_DIR=$(mktemp -d)

# Generate deepcopy helpers
kube::codegen::gen_helpers ./pkg/apis

# Generate client, listers, and informers
kube::codegen::gen_client \
    --output-dir "${TEMP_OUTPUT_DIR}/pkg/generated" \
    --output-pkg github.com/janekbaraniewski/kubeserial/pkg/generated \
    --with-watch \
    ./pkg/apis

printf "Finished kube_codegen.sh, updating source files...\n"

replace_or_compare "${TEMP_OUTPUT_DIR}/pkg/generated/" ./pkg/generated/

rm -rf "${TEMP_OUTPUT_DIR}"

printf "All generators have completed.\n"
