#!/bin/bash
check_env_variable() {
    if [ -z "$1" ]; then
        echo "error, environment variable \"$2\" is empty."
        exit 1
    fi
}
