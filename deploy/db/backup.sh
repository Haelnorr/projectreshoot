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
ACTIVE_DIR="/home/deploy/$ENVR"
DATA_DIR="/home/deploy/data/$ENVR"
BACKUP_DIR="/home/deploy/data/backups/$ENVR"
if [[ "$ENVR" == "production" ]]; then
    SERVICE_NAME="projectreshoot"
    declare -a PORTS=("3000" "3001" "3002")
else
    SERVICE_NAME="$ENVR.projectreshoot"
    declare -a PORTS=("3005" "3006" "3007")
fi

# Send SIGUSR2 to release maintenance mode
release_maintenance() {
    echo "Releasing maintenance mode..."
    for PORT in "${PORTS[@]}"; do
        sudo systemctl kill -s SIGUSR2 "$SERVICE_NAME@$PORT.service"
    done
}

shopt -s nullglob
DB_FILES=("$ACTIVE_DIR"/*.db)
DB_COUNT=${#DB_FILES[@]}

if [[ $DB_COUNT -gt 1 ]]; then
    echo "Error: More than one .db file found in $ACTIVE_DIR. Manual intervention required."
    exit 1
elif [[ $DB_COUNT -eq 0 ]]; then
    echo "Error: No .db file found in $ACTIVE_DIR."
    exit 1
fi

# Extract the filename without extension
DB_FILE="${DB_FILES[0]}"
DB_VER=$(basename "$DB_FILE" .db)

# Send SIGUSR1 to trigger maintenance mode only for active services
declare -a ACTIVE_PORTS=()
for PORT in "${PORTS[@]}"; do
    if systemctl is-active --quiet "$SERVICE_NAME@$PORT.service"; then
        sudo systemctl kill -s SIGUSR1 "$SERVICE_NAME@$PORT.service"
        ACTIVE_PORTS+=("$PORT")
    fi
done
trap release_maintenance EXIT

# Function to check logs for success or failure
check_logs() {
    local port="$1"
    local service="$SERVICE_NAME@$port.service"

    echo "Waiting for $service to enter maintenance mode..."

    # Check the last few lines first in case the message already appeared
    if sudo journalctl -u "$service" -n 20 --no-pager | grep -q "Global database lock acquired"; then
        echo "$service successfully entered maintenance mode."
        return 0
    elif sudo journalctl -u "$service" -n 20 --no-pager | grep -q "Timeout: Global database lock abandoned"; then
        echo "Error: $service failed to enter maintenance mode."
        return 1
    fi

    # If not found, continuously watch logs until we get a success or failure message
    sudo journalctl -u "$service" -f --no-pager | while read -r line; do
        if echo "$line" | grep -q "Global database lock acquired"; then
            echo "$service successfully entered maintenance mode."
            pkill -P $$ journalctl  # Kill journalctl process once we have success
            return 0
        elif echo "$line" | grep -q "Timeout: Global database lock abandoned"; then
            echo "Error: $service failed to enter maintenance mode."
            pkill -P $$ journalctl  # Kill journalctl process on failure
            return 1
        fi
    done
}

# Check logs for each service
for PORT in "${ACTIVE_PORTS[@]}"; do
    check_logs "$PORT"
done

# Get current datetime in YYYY-MM-DD-HHMM format
TIMESTAMP=$(date +"%Y-%m-%d-%H%M")

# Define source and destination paths
SOURCE_DB="$DATA_DIR/$DB_VER.db"
BACKUP_DB="$BACKUP_DIR/${DB_VER}-${TIMESTAMP}.db"

# Copy the database file
if [[ -f "$SOURCE_DB" ]]; then
    cp "$SOURCE_DB" "$BACKUP_DB"
    echo "Backup created: $BACKUP_DB"
else
    echo "Error: Source database file $SOURCE_DB not found."
    exit 1
fi


