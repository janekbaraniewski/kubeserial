#! /bin/bash

mkdir -p /tmp/kubeserial/tests-assets
helm template --name-template kubeserial --show-only templates/config_map.yaml charts/kubeserial > /tmp/kubeserial/tests-assets/config-map.yaml
ls charts/kubeserial/specs | xargs -I{} bash -c "yq eval '.data[\"{}\"]' /tmp/kubeserial/tests-assets/config-map.yaml > pkg/assets/{}"
