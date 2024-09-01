# markdown-journal

A markdown journaling helper

# Quickstart

```bash
# Install
go install github.com/taylorskalyo/markdown-journal@latest

# Create a new journal entry
echo "# Hello World" > "$(date +%Y-%m-%d).md"

# Display journal entries in a timeline view
markdown-journal timeline

# Learn more
markdown-journal --help
```

Works with vim too!

```vim
" Install
Plug 'taylorskalyo/markdown-journal'

" Create a new journal entry
:JournalToday

" Display journal entries in a timeline view
:JournalTimeline

" Learn more
:help journal
```

# Features

`markdown-journal` makes viewing and exploring markdown journals easier.

## Timeline View

The timeline view is simply a markdown formatted index of entries, listed in reverse chronological order. It's an easy way to see your most recent entries or go back in time to revisit old entries.

## Labels View

markdown-journal provides a way to label markdown files. Labels can also be thought of as keywords or categories.
- They can appear anywhere in a markdown file (except for code blocks).
- Labels look like this: `:label:`.
- Any combination of letters, digits, underscores (`_`), and dashes (`-`) between two colons (`:`) creates a label.

The `labels` command generates a markdown formatted list of entries, grouped by label.

## Ctags Integration

By default, the metadata used to generate the timeline and labels views are generated on the fly. However, they can also be cached in a ctags tags file. Other programs can use the tags file to provide additional functionality (e.g. `tags` command or [tagbar](https://github.com/majutsushi/tagbar) plugin in vim)

## Vim Integration

This repo includes a plugin for integrating markdown-journal with vim. See [doc/journal.txt](../blob/master/doc/journal.txt) for a description of the plugin and the commands that it provides.

# Anti-features

markdown-journal...

- is not a markdown editor.
- is not a markdown renderer (see [pandoc](https://pandoc.org/), [hugo](https://gohugo.io/), [jekyll](https://jekyllrb.com/), etc).
- does not impose a particular directory structure. Journal entries can be located anywhere. However, entry file names must begin with a date and end with the markdown extension (i.e. `YYYY-MM-DD*.md`).
- does not force you to use a particular markdown flavor.
