#!/bin/bash
# serve.sh - Usage: ./serve.sh testResponseBody.txt [html port] [css_file] [css_port]
#
HTML_FILE="$1"
HTML_PORT="${2:-1234}"
CSS_FILE="$3"
CSS_PORT="${4:-1235}"

if [ ! -f "$HTML_FILE" ]; then
    echo "Usage: $0 <body-file> [html_port] [css_file] [css_port]"
    echo "Example: $0 index.html 1234 style.css 1235"
    exit 1
fi

serve_html() {
    while true; do
        if [ -f "$HTML_FILE" ]; then
            BODY=$(cat "$HTML_FILE")
            LENGTH=${#BODY}
            (
                echo -e "HTTP/1.1 200 OK\r"
                echo -e "Content-Type: text/html\r"
                echo -e "Content-Length: $LENGTH\r"
                echo -e "Connection: close\r"
                echo -e "\r"
                echo -n "$BODY"
            ) | nc -l -p "$HTML_PORT" -q 0
        else
            echo "Error: HTML file '$HTML_FILE' not found"
            sleep 1
        fi
    done
}

serve_css() {
    if [ -z "$CSS_FILE" ] || [ ! -f "$CSS_FILE" ]; then
        return
    fi

    while true; do
        if [ -f "$CSS_FILE" ]; then
            CSS_BODY=$(cat "$CSS_FILE")
            CSS_LENGTH=${#CSS_BODY}
            (
                echo -e "HTTP/1.1 200 OK\r"
                echo -e "Content-Type: text/css\r"
                echo -e "Content-Length: $CSS_LENGTH\r"
                echo -e "Connection: close\r"
                echo -e "\r"
                echo -n "$CSS_BODY"
            ) | nc -l -p "$CSS_PORT" -q 0
        else
            echo "Error: CSS file '$CSS_FILE' not found"
            sleep 1
        fi
    done
}

# Start HTML server
serve_html &
HTML_PID=$!
echo "Serving HTML on port $HTML_PORT (PID: $HTML_PID)"

# Start CSS server if CSS file is provided
if [ -n "$CSS_FILE" ] && [ -f "$CSS_FILE" ]; then
    serve_css &
    CSS_PID=$!
    echo "Serving CSS on port $CSS_PORT (PID: $CSS_PID)"
else
    if [ -n "$CSS_FILE" ]; then
        echo "Warning: CSS file '$CSS_FILE' not found, skipping CSS server"
    fi
    CSS_PID=""
fi

echo "NOTE: make client fetch 2 times to reload the change"

# Function to cleanup
cleanup() {
    echo ""
    echo "Shutting down..."
    kill $HTML_PID 2>/dev/null
    [ -n "$CSS_PID" ] && kill $CSS_PID 2>/dev/null
    exit 0
}

trap cleanup SIGINT SIGTERM

wait
