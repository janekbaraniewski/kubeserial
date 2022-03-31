#!/usr/bin/env bash

exit_on_error () {
    printf "\n${1}\n\n"
    exit 1
}

replace_or_compare () {
    if [[ "${COPY_OR_DIFF}" == "copy" ]]; then
        cp -r $1 $2
    elif [[ "${COPY_OR_DIFF}" == "diff" ]]; then
        diff -qr $1 $2 || exit_on_error "To fix run:\n    make codegen"
    fi
}
