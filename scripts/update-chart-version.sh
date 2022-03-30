#!/usr/bin/env bash

CHART_PATH="./deploy/chart/kubeserial"
PLACEHOLDER_VALUE="APP_VERSION"
VERSION="works!"

sed="sed"

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     sed=sed;;
    Darwin*)    sed=gsed;;
    *)          sed="UNKNOWN:${unameOut}"
esac

find ${CHART_PATH} \( -type d -name .git -prune \) -o -type f | \
    xargs ${sed} -i 's/'"${PLACEHOLDER_VALUE}"'/'"${VERSION}"'/g'
