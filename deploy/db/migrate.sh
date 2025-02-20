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
VER="$2"
re='^[0-9]+$'
if ! [[ $VER =~ $re ]] ; then
   echo "Error: version not a number" >&2; exit 1
fi

BACKUP_FILE=$(/bin/bash ./test.sh "$ENVR" "$VER" | grep -oP '(?<=Backup created: ).*')
if [[ "$BACKUP_FILE" == "" ]]; then
    echo "Error: backup failed"
    exit 1
fi

# Get current datetime in YYYY-MM-DD-HHMM format
TIMESTAMP=$(date +"%Y-%m-%d-%H%M")

# REAL VALUES
# ACTIVE_DIR="/home/deploy/$ENVR"
# DATA_DIR="/home/deploy/data/$ENVR"
# BACKUP_DIR="/home/deploy/data/backups/$ENVR"
# UPDATED_BACKUP="$BACKUP_DIR/${VER}-${TIMESTAMP}.db"
# UPDATED_COPY="$DATA_DIR/${VER}.db"
# UPDATED_LINK="$ACTIVE_DIR/${VER}.db"
# #####################################################################
# TEST VALUES
ACTIVE_DIR="./deploy/$ENVR"
DATA_DIR="./deploy/data/$ENVR"
BACKUP_DIR="./deploy/data/backups/$ENVR"
UPDATED_BACKUP="$BACKUP_DIR/${VER}-${TIMESTAMP}.db"
UPDATED_COPY="$DATA_DIR/${VER}.db"
UPDATED_LINK="$ACTIVE_DIR/${VER}.db"

# back to real code
# ####################################################################
cp $BACKUP_FILE $UPDATED_BACKUP
# ####################################################################
# do migration
echo "Migration in progress"





# TODO: if failed delete updated_backup
echo "Migration completed"
# end migration
# ####################################################################

cp $UPDATED_BACKUP $UPDATED_COPY
ln -s $UPDATED_COPY $UPDATED_LINK
echo "Upgraded database linked and ready for deploy"
exit 0
