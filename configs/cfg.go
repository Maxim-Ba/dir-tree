package configs

import (
	"fmt"
	"strings"
)

// OutputFormat represents supported output formats
type OutputFormat string

const (
	JSON OutputFormat = "json" // JSON format
	YAML OutputFormat = "yaml" // YAML format
	XML  OutputFormat = "xml"  // XML format
	TXT  OutputFormat = "txt"  // Plain text format
)

// FormatCfg contains formatting configuration options
type FormatCfg struct {
	Type             OutputFormat `json:"type" yaml:"type"`                             // Output format type
	OutputPath       string       `json:"output_path" yaml:"output_path"`               // Output file path (without extension)
	Indent           int          `json:"indent" yaml:"indent"`                         // Indentation for pretty formatting
	ExcludeNodeFields []string    `json:"exclude_node_fields" yaml:"exclude_node_fields"` // Node fields to exclude from output
}

// GetOutputPath returns the output path with appropriate file extension
func (f *FormatCfg) GetOutputPath() string {
    if f.OutputPath == "" {
        return "" // indicates stdout output
    }
    
    // Add extension if missing
    ext := fmt.Sprintf(".%s", string(f.Type))
    if !hasExtension(f.OutputPath, ext) {
        return f.OutputPath + ext
    }
    return f.OutputPath
}

// Config contains all configuration options for directory tree generation
type Config struct {
	Path         string    `json:"path" yaml:"path"`                   // Root directory path
	ExcludeTypes []string  `json:"exclude_types" yaml:"exclude_types"` // File extensions to exclude (e.g., [".tmp", ".log"])
	ExcludePaths []string  `json:"exclude_paths" yaml:"exclude_paths"` // Path patterns to exclude (regex)
	IncludeFiles bool      `json:"include_files" yaml:"include_files"` // Whether to include files or only directories
	MaxDepth     int       `json:"max_depth" yaml:"max_depth"`         // Maximum traversal depth (-1 for unlimited)
	FollowLinks  bool      `json:"follow_links" yaml:"follow_links"`   // Whether to follow symbolic links
	Format       FormatCfg `json:"format" yaml:"format"`               // Formatting configuration
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	if c.MaxDepth < -1 {
		return fmt.Errorf("max depth cannot be less than -1")
	}

	switch c.Format.Type {
	case JSON, YAML, XML, TXT:
		// valid formats
	default:
		return fmt.Errorf("unsupported output format: %s", c.Format.Type)
	}

	return nil
}

// New creates a new ConfigBuilder with default values
func New() *ConfigBuilder {
    return &ConfigBuilder{
        config: &Config{
            Path:         ".",
            MaxDepth:     1,
            IncludeFiles: true,
            FollowLinks:  false,
            Format: FormatCfg{
                Type:       JSON,
                OutputPath: "output-dir",
                Indent:     2,
                ExcludeNodeFields: []string{"size", "is_hidden", "type", "path"},
            },
        },
    }
}

// ConfigBuilder provides a fluent interface for building Config
type ConfigBuilder struct {
    config *Config
}

// WithPath sets the target directory path
func (b *ConfigBuilder) WithPath(path string) *ConfigBuilder {
    b.config.Path = path
    return b
}

// WithMaxDepth sets the maximum traversal depth
func (b *ConfigBuilder) WithMaxDepth(maxDepth int) *ConfigBuilder {
    b.config.MaxDepth = maxDepth
    return b
}

// WithIncludeFiles sets whether to include files in output
func (b *ConfigBuilder) WithIncludeFiles(includeFiles bool) *ConfigBuilder {
    b.config.IncludeFiles = includeFiles
    return b
}

// WithFollowLinks sets whether to follow symbolic links
func (b *ConfigBuilder) WithFollowLinks(followLinks bool) *ConfigBuilder {
    b.config.FollowLinks = followLinks
    return b
}

// WithExcludePaths sets the path exclusion patterns
func (b *ConfigBuilder) WithExcludePaths(excludePaths []string) *ConfigBuilder {
    b.config.ExcludePaths = excludePaths
    return b
}

// WithExcludeTypes sets the file type exclusions
func (b *ConfigBuilder) WithExcludeTypes(excludeTypes []string) *ConfigBuilder {
    b.config.ExcludeTypes = excludeTypes
    return b
}

// WithFormat sets the output format
func (b *ConfigBuilder) WithFormat(format OutputFormat) *ConfigBuilder {
    b.config.Format.Type = format
    return b
}

// WithOutputPath sets the output file path
func (b *ConfigBuilder) WithOutputPath(outputPath string) *ConfigBuilder {
    b.config.Format.OutputPath = outputPath
    return b
}

// WithIndent sets the indentation for formatted output
func (b *ConfigBuilder) WithIndent(indent int) *ConfigBuilder {
    b.config.Format.Indent = indent
    return b
}

// WithExcludeNodeFields sets the node fields to exclude from output
func (b *ConfigBuilder) WithExcludeNodeFields(fields []string) *ConfigBuilder {
    b.config.Format.ExcludeNodeFields = fields
    return b
}

// AddExcludePath adds a path to the exclusion list
func (b *ConfigBuilder) AddExcludePath(path string) *ConfigBuilder {
    b.config.ExcludePaths = append(b.config.ExcludePaths, path)
    return b
}

// AddExcludeType adds a file type to the exclusion list
func (b *ConfigBuilder) AddExcludeType(fileType string) *ConfigBuilder {
    b.config.ExcludeTypes = append(b.config.ExcludeTypes, fileType)
    return b
}

// AddExcludeNodeField adds a node field to the exclusion list
func (b *ConfigBuilder) AddExcludeNodeField(field string) *ConfigBuilder {
    b.config.Format.ExcludeNodeFields = append(b.config.Format.ExcludeNodeFields, field)
    return b
}

// Build returns the final configuration
func (b *ConfigBuilder) Build() *Config {
    // Return a copy to avoid modifications after Build
    return &Config{
        Path:         b.config.Path,
        ExcludeTypes: append([]string{}, b.config.ExcludeTypes...),
        ExcludePaths: append([]string{}, b.config.ExcludePaths...),
        IncludeFiles: b.config.IncludeFiles,
        MaxDepth:     b.config.MaxDepth,
        FollowLinks:  b.config.FollowLinks,
        Format: FormatCfg{
            Type:             b.config.Format.Type,
            OutputPath:       b.config.Format.OutputPath,
            Indent:           b.config.Format.Indent,
            ExcludeNodeFields: append([]string{}, b.config.Format.ExcludeNodeFields...),
        },
    }
}

// hasExtension checks if a path has the specified file extension
func hasExtension(path, ext string) bool {
    if len(path) < len(ext) {
        return false
    }
    // Compare case-insensitively
    return strings.EqualFold(path[len(path)-len(ext):], ext)
}
