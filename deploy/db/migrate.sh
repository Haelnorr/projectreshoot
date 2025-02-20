#!/bin/bash

# Exit on error
set -e

if [[ -z "$1" ]]; then
    echo "Usage: $0 <environment> up-to|down-to <version>"
    exit 1
fi
ENVR="$1"
if [[ "$ENVR" != "production" && "$ENVR" != "staging" ]]; then
    echo "Error: environment must be 'production' or 'staging'."
    exit 1
fi
if [[ -z "$2" ]]; then
    echo "Usage: $0 <environment> up-to|down-to <version>"
    exit 1
fi
CMD="$2"
if [[ "$CMD" != "up-to" && "$CMD" != "down-to" ]]; then
    echo "Error: Command must be 'up-to' or 'down-to'."
    exit 1
fi
if [[ -z "$3" ]]; then
    echo "Usage: $0 <environment> up-to|down-to <version>"
    exit 1
fi
VER="$3"
re='^[0-9]+$'
if ! [[ $VER =~ $re ]] ; then
   echo "Error: version not a number" >&2; exit 1
fi

BACKUP_FILE=$(/bin/bash ./backup.sh "$ENVR" "$VER" | grep -oP '(?<=Backup created: ).*')
if [[ "$BACKUP_FILE" == "" ]]; then
    echo "Error: backup failed"
    exit 1
fi
TIMESTAMP=$(date +"%Y-%m-%d-%H%M")

ACTIVE_DIR="/home/deploy/$ENVR"
DATA_DIR="/home/deploy/data/$ENVR"
BACKUP_DIR="/home/deploy/data/backups/$ENVR"
UPDATED_BACKUP="$BACKUP_DIR/${VER}-${TIMESTAMP}.db"
UPDATED_COPY="$DATA_DIR/${VER}.db"
UPDATED_LINK="$ACTIVE_DIR/${VER}.db"

cp $BACKUP_FILE $UPDATED_BACKUP
failed_cleanup() {
    rm $UPDATED_BACKUP
}
trap 'if [ $? -ne 0 ]; then failed_cleanup; fi' EXIT

echo "Migration in progress"
echo $UPDATED_BACKUP $CMD $VER
./psmigrate $UPDATED_BACKUP $CMD $VER
if [ $? -ne 0 ]; then
    echo "Migration failed"
    exit 1
fi
echo "Migration completed"

cp $UPDATED_BACKUP $UPDATED_COPY
ln -s $UPDATED_COPY $UPDATED_LINK
echo "Upgraded database linked and ready for deploy"
exit 0
