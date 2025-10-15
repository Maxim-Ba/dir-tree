# dir-tree

A Go utility and library for generating directory trees in various formats (JSON, YAML, XML, TXT).

## Features

- Generate directory trees with configurable depth
- Support for multiple output formats (JSON, YAML, XML, TXT)
- Flexible filtering options (exclude paths, file types, node fields)
- Symbolic link handling with follow option
- Both CLI and library APIs available

## Installation

```bash
go get github.com/Maxim-Ba/dir-tree
```

## Usage

### As a CLI Tool

```bash
# Basic usage with default settings
dir-tree

# With custom path and depth
dir-tree -p /path/to/dir -d 3

# Generate YAML output
dir-tree -f yaml -o output

# Exclude specific paths and file types
dir-tree -ep ".git,node_modules" -et ".log,.tmp"
```

### As a Library

```go
package main

import (
    "log"
    
    "github.com/Maxim-Ba/dir-tree/dirtree"
    "github.com/Maxim-Ba/dir-tree/configs"
)

func main() {
    // Quick JSON generation
    data, err := dirtree.GenerateJSON(".", 2)
    if err != nil {
        log.Fatal(err)
    }
    
    // Advanced configuration
    cfg := configs.New().
        WithPath(".").
        WithMaxDepth(2).
        WithFormat(configs.JSON).
        Build()
    
    data, err := dirtree.Generate(cfg)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(data))
}
```

## CLI Flags
- p - Target directory path (default: ".")
- d - Maximum tree depth (default: 1)
- f - Output format: json, yaml, xml, txt (default: json)
- o - Output file path (without extension)
- if - Include files in output (default: true)
- fl - Follow symbolic links (default: false)
- ep - Exclude paths (regex patterns, comma separated)
- et - Exclude file types (extensions, comma separated)
- enf - Exclude node fields (comma separated)
- c - Path to config file

## Config File

```yaml
path: "."
max_depth: 2
include_files: true
follow_links: false
exclude_paths:
  - ".git"
  - "node_modules"
exclude_types:
  - ".tmp"
  - ".log"
format:
  type: "json"
  output_path: "output"
  indent: 2
  exclude_node_fields:
    - "size"
    - "is_hidden"
```
## Output Formats
- JSON: Structured JSON output
- YAML: YAML format for human-readable output
- XML: XML structured output
- TXT: Simple text tree with emoji indicators

## Building from Source

```bash 
git clone https://github.com/Maxim-Ba/dir-tree
cd dir-tree
go build -o dir-tree main.go
```
