#!/usr/bin/env bash

CHART_PATH="./deploy/chart/kubeserial"
PLACEHOLDER_VALUE="APP_VERSION"
VERSION="works!"

find ./deploy/chart/kubeserial \( -type d -name .git -prune \) -o -type f | \
    xargs -0 sed -i 's/${PLACEHOLDER_VALUE}/${VERSION}/g'
