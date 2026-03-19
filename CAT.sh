#!/usr/bin/env bash

ROOT="${1:-.}"
OUTPUT="${2:-project.txt}"

# Clear output file first
> "$OUTPUT"

find "$ROOT" \
  -type d \( -path "$ROOT/.git" -o -path "$ROOT/tes" \) -prune \
  -o -type f -print | while IFS= read -r file; do

    echo "===== FILE: $file =====" >> "$OUTPUT"
    cat "$file" >> "$OUTPUT"
    echo -e "\n" >> "$OUTPUT"

done

echo "Done. Output saved to $OUTPUT"

