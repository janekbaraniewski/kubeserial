#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

: ${GO=go}

export GOPATH=$("$GO" env GOPATH | awk -F ':' '{print $1}')
export PATH=$PATH:$GOPATH/bin

source hack/utils.sh

GO111MODULE=on "$GO" install \
    sigs.k8s.io/controller-tools/cmd/controller-gen

printf "controller-gen... "

controller-gen rbac:roleName=manager-role crd paths=./pkg/apis/kubeserial/... output:crd:dir=/tmp/deploy/crds
find /tmp/deploy/crds -name "*.yaml" | xargs -I % python3 ./hack/update-crd-metadata.py % ./hack/crd_metadata_template.yaml
replace_or_compare /tmp/deploy/crds/* deploy/chart/kubeserial-crds/templates/
rm -r /tmp/deploy/crds

printf "Done!\n"
