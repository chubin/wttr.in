#!/bin/bash
# scripts/update_golden.sh
# Update golden files by fetching current output from local wttr.in server

set -euo pipefail

SERVER="http://127.0.0.1:9005"
TESTCASES_FILE="internal/renderer/v1/testdata/testcases.json"
TESTDATA_DIR="internal/renderer/v1/testdata"

echo "=== Updating golden files from local server ==="
echo "Server: $SERVER"
echo "Test cases: $TESTCASES_FILE"
echo "Output directory: $TESTDATA_DIR"
echo

# Check if server is running
if ! curl -s --max-time 2 "$SERVER/heideck" > /dev/null; then
    echo "❌ Error: Server is not responding at $SERVER"
    echo "   Please start the server with: go run ."
    exit 1
fi

# Read test cases using jq (more robust than pure bash)
if ! command -v jq >/dev/null 2>&1; then
    echo "❌ Error: jq is required but not installed."
    echo "   Install jq with: sudo apt install jq   (or brew install jq on macOS)"
    exit 1
fi

# Create testdata directory if it doesn't exist
mkdir -p "$TESTDATA_DIR"

count=0
updated=0

while read -r testcase; do
    name=$(echo "$testcase" | jq -r '.name')
    golden=$(echo "$testcase" | jq -r '.golden')
    options=$(echo "$testcase" | jq -r '.query.Options.options // ""')

    golden_path="$TESTDATA_DIR/${golden}.txt"
    url="$SERVER/heideck"

    if [[ -n "$options" ]]; then
        url="${url}?${options}"
    fi

    echo -n "Fetching: $name → ${golden}.txt ... "

    # Fetch the output and trim trailing whitespace/newlines for consistency with tests
    if output=$(curl -s --max-time 10 "$url"); then
        # Trim trailing newlines and whitespace (matches test behavior)
        echo "$output" | sed 's/[[:space:]]*$//' > "$golden_path"
        
        echo "✓ updated"
        ((updated++)) || true
    else
        echo "✗ failed"
        echo "   URL: $url"
    fi

    ((count++)) || true
done <<< "$(jq -c '.testcases[]' "$TESTCASES_FILE")"

echo
echo "=== Summary ==="
echo "Processed: $count test cases"
echo "Updated:   $updated golden files"
echo "Done! You can now run the tests with:"
echo "   go test ./internal/renderer/v1 -run TestV1Renderer_Render"
