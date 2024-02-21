#!/bin/bash

# Validate usage (optional, enhances robustness)
if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <chunk_size_in_lines> <input_file>"
    exit 1
fi

chunk_size="$1"
input_file="$2"
output_dir="data"

# Create output directory if it doesn't exist (optional, ensures correct behavior)
if [[ ! -d "$output_dir" ]]; then
    mkdir -p "$output_dir"
fi

# Use efficient `split` command directly
split -l "$chunk_size" "$input_file" "$output_dir/chunk_"

# Alternative splitting logic (consider for more custom requirements)
# (Optional. Enable if needed)
# cat "$input_file" | split -l "$chunk_size" -a 5 --numeric-suffixes - "" - | while read -r chunk_lines; do
#     filename="$output_dir/chunk_$((++count)).txt"
#     echo "$headers" > "$filename"
#     echo "$chunk_lines" >> "$filename"
# done

echo "Split into chunks (using split command)"
