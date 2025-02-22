#!/bin/bash

if [[ -z "$1" ]]; then
    echo "Usage: $0 <environment> <version> <commit-hash>"
    exit 1
fi
ENVR="$1"
if [[ "$ENVR" != "production" && "$ENVR" != "staging" ]]; then
    echo "Error: environment must be 'production' or 'staging'."
    exit 1
fi
if [[ -z "$2" ]]; then
    echo "Usage: $0 <environment> <version> <commit-hash>"
    exit 1
fi
TGT_VER="$2"
re='^[0-9]+$'
if ! [[ $TGT_VER =~ $re ]] ; then
   echo "Error: version not a number" >&2
   exit 1
fi
if [ -z "$3" ]; then
    echo "Usage: $0 <environment> <version> <commit-hash>"
  exit 1
fi
COMMIT_HASH=$3
MIGRATION_BIN="/home/deploy/migration-bin"
BACKUP_OUTPUT=$(/bin/bash ${MIGRATION_BIN}/backup.sh "$ENVR" 2>&1)
echo "$BACKUP_OUTPUT"
if [[ $? -ne 0 ]]; then
    exit 1
fi
BACKUP_FILE=$(echo "$BACKUP_OUTPUT" | grep -oP '(?<=Backup created: ).*')
if [[ -z "$BACKUP_FILE" ]]; then
    echo "Error: backup failed"
    exit 1
fi

FILE_NAME=${BACKUP_FILE##*/}
CUR_VER=${FILE_NAME%%-*}
if [[ $((+$TGT_VER)) == $((+$CUR_VER)) ]]; then
    echo "Version same, skipping migration"
    exit 0
fi
if [[ $((+$TGT_VER)) > $((+$CUR_VER)) ]]; then
    CMD="up-to"
fi
if [[ $((+$TGT_VER)) < $((+$CUR_VER)) ]]; then
    CMD="down-to"
fi
TIMESTAMP=$(date +"%Y-%m-%d-%H%M")

ACTIVE_DIR="/home/deploy/$ENVR"
DATA_DIR="/home/deploy/data/$ENVR"
BACKUP_DIR="/home/deploy/data/backups/$ENVR"
UPDATED_BACKUP="$BACKUP_DIR/${TGT_VER}-${TIMESTAMP}.db"
UPDATED_COPY="$DATA_DIR/${TGT_VER}.db"
UPDATED_LINK="$ACTIVE_DIR/${TGT_VER}.db"

cp $BACKUP_FILE $UPDATED_BACKUP
failed_cleanup() {
    rm $UPDATED_BACKUP
}
trap 'if [ $? -ne 0 ]; then failed_cleanup; fi' EXIT

echo "Migration in progress from $CUR_VER to $TGT_VER"
${MIGRATION_BIN}/prmigrate-$COMMIT_HASH $UPDATED_BACKUP $CMD $TGT_VER
if [ $? -ne 0 ]; then
    echo "Migration failed"
    exit 1
fi
echo "Migration completed"

cp $UPDATED_BACKUP $UPDATED_COPY
ln -s $UPDATED_COPY $UPDATED_LINK
echo "Upgraded database linked and ready for deploy"
exit 0
