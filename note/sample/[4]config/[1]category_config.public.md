# Category Config

```toml
# Category Config

# display_name of the category
# if empty, defaults to name
display_name = "display_name"

# name of the category
# if empty, defaults to the directory name in file system
name = "category_name"

# index file of the category
# the content of index file fills {{ .Content }} of the index template or category template
index = "index.md"

# overrides note_file_pattern in site_config.toml
note_file_pattern = "^(?:\\[.*?\\])*(.*)\\.public\\.(?:txt|html|md)$"
```
