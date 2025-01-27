#!/bin/bash

# Directory where the script is located (staging/src/k8s.io/code-generator/validation-gen/)
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Run the command and capture the JSON output
json_output=$(go run "${script_dir}" --docs)

# Create a Markdown file in the same directory
markdown_file="${script_dir}/validation_tags.md"

# Write the header to the Markdown file
{
    echo "# Kubernetes Validation Tags Documentation"
    echo ""
    echo "This document lists the supported validation tags and their related information."
    echo ""
    echo "## Tags Overview"
    echo ""
    echo "| Tag | Description | Contexts |"
    echo "|-----|-------------|----------|"
} > "$markdown_file"

# Process the JSON output and populate the main table
{
    echo "$json_output" | jq -c '.[]' | while read -r line; do
        tag=$(echo "$line" | jq -r '.Tag')
        description=$(echo "$line" | jq -r '.Description')
        contexts=$(echo "$line" | jq -r '.Contexts | join(", ")')
        payloads=$(echo "$line" | jq -r '.Payloads')

        # Add row to main table with link to payloads if they exist
        if [[ "$payloads" != "null" ]]; then
            tag_link=$(echo "$tag" | sed 's/k8s:/k8s/' | tr '[:upper:]' '[:lower:]')
            echo "| [\`$tag\`](#$tag_link) | $description | $contexts |"
        else
            echo "| \`$tag\` | $description | $contexts |"
        fi
    done

    echo ""
    echo "## Tag Details"
    echo ""
} >> "$markdown_file"

# Create separate sections for each tag's payloads
echo "$json_output" | jq -c '.[]' | while read -r line; do
    tag=$(echo "$line" | jq -r '.Tag')
    payloads=$(echo "$line" | jq -r '.Payloads')

    if [[ "$payloads" != "null" ]]; then
        # Create section for this tag
        {
            echo "### $tag"
            echo ""
            echo "| Description | Docs | Schema |"
            echo "|-------------|------|---------|"
            echo "$payloads" | jq -r '.[] | 
                "| **" + 
                (.Description | gsub("<"; "\\<") | gsub(">"; "\\>")) + 
                "** | " +
                (if .Docs then .Docs else "" end) +
                " | " +
                (if .Schema then
                    (.Schema | map(
                        "- `" + .Key + "`: `" + .Value + "`" +
                        (if .Docs then " (" + .Docs + ")" else "" end)
                    ) | join("<br>"))
                else "None" end) +
                " |"'
            echo ""
        } >> "$markdown_file"
    fi
done

echo "Markdown documentation generated at $markdown_file"