#!/bin/bash

# Exit on error
set -e

if [[ -z "$1" ]]; then
    echo "Usage: $0 <environment>"
    exit 1
fi

ENVR="$1"
if [[ "$ENVR" != "production" && "$ENVR" != "staging" ]]; then
    echo "Error: environment must be 'production' or 'staging'."
    exit 1
fi
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "random output"
echo "Error"
FILE="00001-2025-02-19-0927.db"
ACTIVE_DIR="./deploy/$ENVR"
DATA_DIR="./deploy/data/$ENVR"
BACKUP_DIR="./deploy/data/backups/$ENVR"
mkdir -p $ACTIVE_DIR
mkdir -p $DATA_DIR
mkdir -p $BACKUP_DIR
touch "$BACKUP_DIR/$FILE"
echo "Backup created: $BACKUP_DIR/$FILE"
