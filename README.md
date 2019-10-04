# markdown-journal

A markdown journaling helper

# Quickstart

```bash
# Install
go get github.com/taylorskalyo/markdown-journal

# Create a new journal entry
echo "# Hello World" > "$(date +%Y-%m-%d).md"

# Display journal entries in a timeline view
markdown-journal timeline

# Learn more
markdown-journal --help
```

# What markdown-journal is

`markdown-journal` makes viewing and exploring markdown journals easier.

- It provides commands for creating indexes of dated markdown files.
  - The `timeline` command generates a markdown formated index of entries in chronological order.
  - The `taglist` command (not yet implemented) generates a markdown formated list of entries, grouped by tag.
- It provides a way to tag markdown files. Tags can also be thought of a keywords or categories.
  - Tags can appear anywhere in a markdown file (except for code blocks and code fences).
  - Tags look like this: `:tag:`.
  - Any combination of letters, digits, underscores (`_`), and dashes (`-`) between two colons (`:`) creates a tag.
- It generates ctags compatible tag files.
  - A ctags file can be used to quickly locate journal tags and headers.
  - Other programs can use these ctags files to provide additional functionality (e.g. `tags` command or [tagbar](https://github.com/majutsushi/tagbar) plugin in vim)

# What markdown-journal is not

- It is not a markdown editor.
- It is not a markdown renderer (see [pandoc](https://pandoc.org/), [hugo](https://gohugo.io/), [jekyll](https://jekyllrb.com/), etc).
- It does not impose a particular directory structure.
  - Journal entries can be located anywhere.
  - However, entry file names must begin with a date and end with the markdown extension (i.e. `YYYY-MM-DD*.md`).
- It does not force you to use a particular markdown flavor.

# Planned features

- [ ] Add a `taglist` command to generate a markdown formated list of entries, grouped by tag.
- [ ] Optionally pull tags from hugo/jekyll front matter.
- [ ] Better feature parity with ctags/universl-ctags/exuberant-ctags
