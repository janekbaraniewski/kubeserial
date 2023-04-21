#!/usr/bin/env bash

exit_on_error () {
    printf "\n${1}\n\n"
    exit 1
}

replace_or_compare () {
    if [[ "${COPY_OR_DIFF}" == "copy" ]]; then
        if [[ -d $1 ]]; then
            cp -r $1/* $2
        elif [[ -f $1 ]]; then
            cp -r $1 $2
        else
            echo "$1 is not valid"
            exit 1
        fi
    elif [[ "${COPY_OR_DIFF}" == "diff" ]]; then
        diff -r --exclude="_helpers.tpl" --exclude="**/fake/*" $1 $2 || exit_on_error "To fix run:\n    make generate"
    fi
}
