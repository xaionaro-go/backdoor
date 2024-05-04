#!/bin/bash

errorExit() {
    EXIT_CODE="$1"; shift
    ERROR_MSG="$1"; shift

    echo "$ERROR_MSG" >&2
    exit "$EXIT_CODE"
}

FILE_PATH="$1"; shift
if [ "$FILE_PATH" = '' ]; then
    errorExit 1 'No file provided. Syntax: upload_file.sh <file> <backdoor host> <backdoor port> <destination path>'
fi

BACKDOOR_HOST="$1"; shift
if [ "$BACKDOOR_HOST" = '' ]; then
    errorExit 1 'No host provided. Syntax: upload_file.sh <file> <backdoor host> <backdoor port> <destination path>'
fi

BACKDOOR_PORT="$1"; shift
if [ "$BACKDOOR_PORT" = '' ]; then
    errorExit 1 'No port provided. Syntax: upload_file.sh <file> <backdoor host> <backdoor port> <destination path>'
fi

DEST_PATH="$1"; shift
if [ "$DEST_PATH" = '' ]; then
    errorExit 1 'No destination dir provided. Syntax: upload_file.sh <file> <backdoor host> <backdoor port> <destination path>'
fi

(
    sleep 1
    echo "set -x"
    echo "cat <<EOF | base64 -d > '$DEST_PATH'"
    base64 "$FILE_PATH"
    echo "EOF"
    echo "exit"
) | ssh -ttp "$BACKDOOR_PORT" "$BACKDOOR_HOST" > /dev/null
