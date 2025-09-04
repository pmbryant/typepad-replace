#!/bin/bash

set -e

DEST_DIR="_site"

# Remove all contents of destination directory if it exists
if [ -d "$DEST_DIR" ]; then
    rm -rf "$DEST_DIR"/*
fi

# Create the destination directory if it doesn't exist
mkdir -p "$DEST_DIR"

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

# Copy contents of docs into destination directory
if [ -d "docs" ]; then
    cp -a docs/. "$DEST_DIR/"
fi

echo "Merge complete."
