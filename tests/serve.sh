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
        echo "HTTP/1.1 200 OK"
        echo "Content-Type: text/html"
        echo "Content-Length: $LENGTH"
        echo "Connection: close"
        echo ""
        echo -n "$BODY"
    ) | nc -l -p "$PORT"
done
