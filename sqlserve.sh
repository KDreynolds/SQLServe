#!/bin/sh
# This is a self-extracting database server
# The actual database starts after this script

# Get the directory where this script is located
SCRIPT_DIR="$(dirname "$0")"
SCRIPT_PATH="$(realpath "$0")"

# Create a temporary directory for our files
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

# Get the size of this script in bytes
SCRIPT_SIZE=$(wc -c < "$0")
# Skip the script portion to get to the database
dd if="$0" of="$TMP_DIR/db.sqlite" bs=1 skip="$SCRIPT_SIZE" 2>/dev/null

# Extract the binary if it doesn't exist
BIN_PATH="$SCRIPT_DIR/sqlserve_bin"
if [ ! -f "$BIN_PATH" ]; then
    # Find the binary in the same directory as the script
    if [ -f "$SCRIPT_DIR/sqlserve_bin" ]; then
        BIN_PATH="$SCRIPT_DIR/sqlserve_bin"
    else
        echo "Error: sqlserve_bin not found in $SCRIPT_DIR"
        exit 1
    fi
fi

# Execute the binary with the extracted database
exec "$BIN_PATH" --db "$TMP_DIR/db.sqlite" "$@"

# End of script 