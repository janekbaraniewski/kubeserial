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
go install k8s.io/code-generator

printf "Got all dependencies \n"

# Ensure proper package tags are in place (See: https://pkg.go.dev/k8s.io/code-generator/cmd/deepcopy-gen)
printf "Running code generators...\n"

chmod +x "$GOPATH"/pkg/mod/k8s.io/code-generator@v0.30.0/kube_codegen.sh

# Use generate-groups.sh helper script to run all code generators
"$GOPATH"/pkg/mod/k8s.io/code-generator@v0.30.0/kube_codegen.sh all \
    github.com/janekbaraniewski/kubeserial/pkg/generated \
    github.com/janekbaraniewski/kubeserial/pkg/apis \
    v1alpha1

# Manual copy might be required if output paths are incorrectly set by the tools, adjust paths as needed
if [[ "${COPY_OR_DIFF}" == "copy" ]]; then
    rm -rf ./pkg/generated
    mkdir -p ./pkg/generated
    cp -r "$GOPATH/src/github.com/janekbaraniewski/kubeserial/pkg/generated/"* ./pkg/generated/
fi

replace_or_compare "$GOPATH/src/github.com/janekbaraniewski/kubeserial/pkg/generated/" ./pkg/generated/

printf "All generators have completed.\n"
