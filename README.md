# markdown-journal

A markdown journaling system

# Quickstart

```bash
# Install
go get github.com/taylorskalyo/markdown-journal

# Create a new journal entry wherever you want
echo "# Hello World" > "$(date +%Y-%m-%d).md"

# Display journal entries a timeline view
markdown-journal timeline *.md

# Learn more
markdown-journal --help
```
