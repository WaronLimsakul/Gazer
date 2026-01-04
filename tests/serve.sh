#!/bin/bash
# serve.sh - Usage: ./serve.sh testResponseBody.txt [port]
BODY_FILE="$1"
PORT="${2:-1234}"

if [ ! -f "$BODY_FILE" ]; then
    echo "Usage: $0 <body-file> [port]"
    exit 1
fi

BODY=$(cat "$BODY_FILE")
LENGTH=${#BODY}

while true; do
    (
        echo -e "HTTP/1.1 200 OK\r"
        echo -e "Content-Type: text/html\r"
        echo -e "Content-Length: $LENGTH\r"
        echo -e "Connection: close\r"
        echo -e "\r"
        echo -n "$BODY"
    ) | nc -l -p "$PORT"
done
