#!/bin/bash

# Exit on error
set -e

if [[ -z "$1" ]]; then
    echo "Usage: $0 <environment> <version>"
    exit 1
fi

ENVR="$1"
if [[ "$ENVR" != "production" && "$ENVR" != "staging" ]]; then
    echo "Error: environment must be 'production' or 'staging'."
    exit 1
fi
if [[ -z "$2" ]]; then
    echo "Usage: $0 <environment> <version>"
    exit 1
fi
TGT_VER="$2"
re='^[0-9]+$'
if ! [[ $TGT_VER =~ $re ]] ; then
   echo "Error: version not a number" >&2
   exit 1
fi
ACTIVE_DIR="/home/deploy/$ENVR"
find "$ACTIVE_DIR" -type l -name "*.db" ! -name "${TGT_VER}.db" -exec rm -v {} +
