#!/bin/bash
OUTPUT="Markdowns/AI-Agent-Full-Prompt.md"

echo "The following is the Software Design Document and the complete codebase to reproduce the project." > "$OUTPUT"
echo "" >> "$OUTPUT"

cat Markdowns/AI-agent-SDD.md >> "$OUTPUT"
echo "" >> "$OUTPUT"
echo "## 16. Codebase Files" >> "$OUTPUT"
echo "" >> "$OUTPUT"

# List all files and append them to the markdown
find . -type f \
    -not -path '*/\.git/*' \
    -not -path '*/\.idea/*' \
    -not -path '*/Markdowns/*' \
    -not -name 'generate_prompt.sh' \
    -not -name 'go.sum' \
    | sort | while read -r file; do
    
    echo "### File: ${file#./}" >> "$OUTPUT"
    
    ext="${file##*.}"
    lang=""
    case "$ext" in
        go) lang="go" ;;
        sql) lang="sql" ;;
        yml|yaml) lang="yaml" ;;
        mod) lang="go" ;;
        md) lang="markdown" ;;
        *) lang="text" ;;
    esac
    
    if [[ $(basename "$file") == "Makefile" || $(basename "$file") == "Dockerfile" ]]; then
        lang="makefile"
        [[ $(basename "$file") == "Dockerfile" ]] && lang="dockerfile"
    fi

    echo '```'$lang >> "$OUTPUT"
    cat "$file" >> "$OUTPUT"
    echo '```' >> "$OUTPUT"
    echo "" >> "$OUTPUT"
done

echo "Done generating $OUTPUT"
