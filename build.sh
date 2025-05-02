#!/bin/sh
set -e

# Build the Go binary
echo "Building Go binary..."
go build -o sqlserve_bin main.go

# Create initial database if it doesn't exist
if [ ! -f sqlserve.db ]; then
    echo "Creating initial database..."
    sqlite3 sqlserve.db < schema.sql
fi

# Create the combined file
echo "Creating combined server..."
# First, copy the shell script
cp sqlserve.sh sqlserve
# Then append the database
cat sqlserve.db >> sqlserve

# Make it executable
chmod +x sqlserve

echo "Done! You can now run ./sqlserve" 