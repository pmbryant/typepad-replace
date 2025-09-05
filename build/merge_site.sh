#!/bin/bash

set -e

DEST_DIR="_site"

# Remove all contents of destination directory if it exists
if [ -d "$DEST_DIR" ]; then
    rm -rf "$DEST_DIR"/*
fi

# Create the destination directory if it doesn't exist
mkdir -p "$DEST_DIR"

# Always copy contents of docs into destination directory
if [ -d "docs" ]; then
    cp -a docs/. "$DEST_DIR/"
fi

# If parameter is 'disabled', exit after copying docs
if [ "$1" = "disabled" ]; then
    echo "Script disabled: copied docs and exited."
    exit 0
fi

# Function to merge subdirectories from a source into destination
merge_subdirs() {
    src_dir="$1"
    for sub in "$src_dir"/*; do
        if [ -d "$sub" ]; then
            subname=$(basename "$sub")
            # Rename pmbryant.typepad.com to x
            if [ "$subname" = "pmbryant.typepad.com" ]; then
                subname="x"
            fi
            mkdir -p "$DEST_DIR/$subname"
            cp -a "$sub"/. "$DEST_DIR/$subname/"
        fi
    done
}

# Merge subdirs under bb-blog
merge_subdirs "bb-blog"

# Merge subdirs under lyg-blog
merge_subdirs "lyg-blog"


echo "Merge complete."
