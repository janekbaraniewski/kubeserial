#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

: ${GO=go}

export GOPATH=$("$GO" env GOPATH | awk -F ':' '{print $1}')
export PATH=$PATH:$GOPATH/bin

source hack/utils.sh

# Install the required binaries with modules enabled
GO111MODULE=on "$GO" install \
    k8s.io/code-generator/cmd/deepcopy-gen \
    k8s.io/code-generator/cmd/register-gen \
    k8s.io/code-generator/cmd/client-gen \
    k8s.io/code-generator/cmd/lister-gen \
    k8s.io/code-generator/cmd/informer-gen

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
