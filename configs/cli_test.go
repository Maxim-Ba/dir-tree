package configs

import (
	"testing"
)

// TestGetOutputPath tests the GetOutputPath method
func TestGetOutputPath(t *testing.T) {
	tests := []struct {
		name     string
		format   FormatCfg
		expected string
	}{
		{
			name: "Empty output path",
			format: FormatCfg{
				Type:       JSON,
				OutputPath: "",
			},
			expected: "",
		},
		{
			name: "JSON with existing extension",
			format: FormatCfg{
				Type:       JSON,
				OutputPath: "output.json",
			},
			expected: "output.json",
		},
		{
			name: "JSON without extension",
			format: FormatCfg{
				Type:       JSON,
				OutputPath: "output",
			},
			expected: "output.json",
		},
		{
			name: "YAML with existing extension",
			format: FormatCfg{
				Type:       YAML,
				OutputPath: "output.yaml",
			},
			expected: "output.yaml",
		},
		{
			name: "YAML without extension",
			format: FormatCfg{
				Type:       YAML,
				OutputPath: "output",
			},
			expected: "output.yaml",
		},
		{
			name: "XML with existing extension",
			format: FormatCfg{
				Type:       XML,
				OutputPath: "output.xml",
			},
			expected: "output.xml",
		},
		{
			name: "XML without extension",
			format: FormatCfg{
				Type:       XML,
				OutputPath: "output",
			},
			expected: "output.xml",
		},
		{
			name: "TXT with existing extension",
			format: FormatCfg{
				Type:       TXT,
				OutputPath: "output.txt",
			},
			expected: "output.txt",
		},
		{
			name: "TXT without extension",
			format: FormatCfg{
				Type:       TXT,
				OutputPath: "output",
			},
			expected: "output.txt",
		},
		{
			name: "Different extension case - uppercase",
			format: FormatCfg{
				Type:       JSON,
				OutputPath: "output.JSON",
			},
			expected: "output.JSON",
		},
		{
			name: "Different extension case - mixed case",
			format: FormatCfg{
				Type:       JSON,
				OutputPath: "output.JsOn",
			},
			expected: "output.JsOn",
		},
		{
			name: "Path with directory",
			format: FormatCfg{
				Type:       JSON,
				OutputPath: "dir/output",
			},
			expected: "dir/output.json",
		},
		{
			name: "Path with directory and existing extension",
			format: FormatCfg{
				Type:       JSON,
				OutputPath: "dir/output.json",
			},
			expected: "dir/output.json",
		},
		{
			name: "Wrong extension for format",
			format: FormatCfg{
				Type:       JSON,
				OutputPath: "output.txt",
			},
			expected: "output.txt.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format.GetOutputPath()
			if result != tt.expected {
				t.Errorf("GetOutputPath() = %s, want %s", result, tt.expected)
			}
		})
	}
}
// TestConfigValidate tests the Validate method
func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		shouldError bool
	}{
		{
			name: "Valid config",
			config: &Config{
				Path:     "/valid/path",
				MaxDepth: 5,
				Format: FormatCfg{
					Type: JSON,
				},
			},
			shouldError: false,
		},
		{
			name: "Empty path",
			config: &Config{
				Path:     "",
				MaxDepth: 1,
				Format: FormatCfg{
					Type: JSON,
				},
			},
			shouldError: true,
		},
		{
			name: "Max depth less than -1",
			config: &Config{
				Path:     "/valid/path",
				MaxDepth: -2,
				Format: FormatCfg{
					Type: JSON,
				},
			},
			shouldError: true,
		},
		{
			name: "Valid max depth -1",
			config: &Config{
				Path:     "/valid/path",
				MaxDepth: -1,
				Format: FormatCfg{
					Type: JSON,
				},
			},
			shouldError: false,
		},
		{
			name: "Unsupported output format",
			config: &Config{
				Path:     "/valid/path",
				MaxDepth: 1,
				Format: FormatCfg{
					Type: "invalid",
				},
			},
			shouldError: true,
		},
		{
			name: "Valid YAML format",
			config: &Config{
				Path:     "/valid/path",
				MaxDepth: 1,
				Format: FormatCfg{
					Type: YAML,
				},
			},
			shouldError: false,
		},
		{
			name: "Valid XML format",
			config: &Config{
				Path:     "/valid/path",
				MaxDepth: 1,
				Format: FormatCfg{
					Type: XML,
				},
			},
			shouldError: false,
		},
		{
			name: "Valid TXT format",
			config: &Config{
				Path:     "/valid/path",
				MaxDepth: 1,
				Format: FormatCfg{
					Type: TXT,
				},
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// TestConfigBuilder tests the ConfigBuilder methods
func TestConfigBuilder(t *testing.T) {
	tests := []struct {
		name     string
		build    func(*ConfigBuilder) *ConfigBuilder
		expected *Config
	}{
		{
			name: "Default config",
			build: func(b *ConfigBuilder) *ConfigBuilder {
				return b
			},
			expected: &Config{
				Path:         ".",
				MaxDepth:     1,
				ExcludePaths: []string{},
				ExcludeTypes: []string{},
				IncludeFiles: true,
				FollowLinks:  false,
				Format: FormatCfg{
					Type:              JSON,
					OutputPath:        "output-dir",
					Indent:            2,
					ExcludeNodeFields: []string{"size", "is_hidden", "type", "path"},
				},
			},
		},
		{
			name: "With custom path and depth",
			build: func(b *ConfigBuilder) *ConfigBuilder {
				return b.WithPath("/custom/path").WithMaxDepth(5)
			},
			expected: &Config{
				Path:         "/custom/path",
				MaxDepth:     5,
				ExcludePaths: []string{},
				ExcludeTypes: []string{},
				IncludeFiles: true,
				FollowLinks:  false,
				Format: FormatCfg{
					Type:              JSON,
					OutputPath:        "output-dir",
					Indent:            2,
					ExcludeNodeFields: []string{"size", "is_hidden", "type", "path"},
				},
			},
		},
		{
			name: "With exclude paths and types",
			build: func(b *ConfigBuilder) *ConfigBuilder {
				return b.WithExcludePaths([]string{".git", "node_modules"}).
					WithExcludeTypes([]string{".tmp", ".log"})
			},
			expected: &Config{
				Path:         ".",
				MaxDepth:     1,
				ExcludePaths: []string{".git", "node_modules"},
				ExcludeTypes: []string{".tmp", ".log"},
				IncludeFiles: true,
				FollowLinks:  false,
				Format: FormatCfg{
					Type:              JSON,
					OutputPath:        "output-dir",
					Indent:            2,
					ExcludeNodeFields: []string{"size", "is_hidden", "type", "path"},
				},
			},
		},
		{
			name: "With format settings",
			build: func(b *ConfigBuilder) *ConfigBuilder {
				return b.WithFormat(YAML).
					WithOutputPath("custom-output").
					WithIndent(4).
					WithExcludeNodeFields([]string{"size", "children"})
			},
			expected: &Config{
				Path:         ".",
				MaxDepth:     1,
				ExcludePaths: []string{},
				ExcludeTypes: []string{},
				IncludeFiles: true,
				FollowLinks:  false,
				Format: FormatCfg{
					Type:              YAML,
					OutputPath:        "custom-output",
					Indent:            4,
					ExcludeNodeFields: []string{"size", "children"},
				},
			},
		},
		{
			name: "With boolean flags",
			build: func(b *ConfigBuilder) *ConfigBuilder {
				return b.WithIncludeFiles(false).WithFollowLinks(true)
			},
			expected: &Config{
				Path:         ".",
				MaxDepth:     1,
				ExcludePaths: []string{},
				ExcludeTypes: []string{},
				IncludeFiles: false,
				FollowLinks:  true,
				Format: FormatCfg{
					Type:              JSON,
					OutputPath:        "output-dir",
					Indent:            2,
					ExcludeNodeFields: []string{"size", "is_hidden", "type", "path"},
				},
			},
		},
		{
			name: "Add exclusions incrementally",
			build: func(b *ConfigBuilder) *ConfigBuilder {
				return b.AddExcludePath(".git").
					AddExcludePath("node_modules").
					AddExcludeType(".tmp").
					AddExcludeType(".log").
					AddExcludeNodeField("size").
					AddExcludeNodeField("children")
			},
			expected: &Config{
				Path:         ".",
				MaxDepth:     1,
				ExcludePaths: []string{".git", "node_modules"},
				ExcludeTypes: []string{".tmp", ".log"},
				IncludeFiles: true,
				FollowLinks:  false,
				Format: FormatCfg{
					Type:              JSON,
					OutputPath:        "output-dir",
					Indent:            2,
					ExcludeNodeFields: []string{"size", "is_hidden", "type", "path", "size", "children"},
				},
			},
		},
		{
			name: "Chained methods",
			build: func(b *ConfigBuilder) *ConfigBuilder {
				return b.WithPath("/path").
					WithMaxDepth(3).
					WithIncludeFiles(false).
					WithFollowLinks(true).
					WithFormat(XML).
					WithOutputPath("xml-output")
			},
			expected: &Config{
				Path:         "/path",
				MaxDepth:     3,
				ExcludePaths: []string{},
				ExcludeTypes: []string{},
				IncludeFiles: false,
				FollowLinks:  true,
				Format: FormatCfg{
					Type:              XML,
					OutputPath:        "xml-output",
					Indent:            2,
					ExcludeNodeFields: []string{"size", "is_hidden", "type", "path"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := New()
			config := tt.build(builder).Build()

			compareConfig(t, config, tt.expected)
		})
	}
}

// TestConfigBuilderIsolation tests that built configs are isolated
func TestConfigBuilderIsolation(t *testing.T) {
	builder := New()

	// Build first config
	config1 := builder.WithPath("/path1").WithMaxDepth(5).Build()

	// Build second config with different values
	config2 := builder.WithPath("/path2").WithMaxDepth(10).Build()

	// Configs should be independent
	if config1.Path != "/path1" {
		t.Errorf("config1.Path = %s, want %s", config1.Path, "/path1")
	}
	if config1.MaxDepth != 5 {
		t.Errorf("config1.MaxDepth = %d, want %d", config1.MaxDepth, 5)
	}

	if config2.Path != "/path2" {
		t.Errorf("config2.Path = %s, want %s", config2.Path, "/path2")
	}
	if config2.MaxDepth != 10 {
		t.Errorf("config2.MaxDepth = %d, want %d", config2.MaxDepth, 10)
	}
}

// TestHasExtension tests the hasExtension function (via GetOutputPath)
func TestHasExtension(t *testing.T) {
	// We test hasExtension indirectly through GetOutputPath
	tests := []struct {
		path     string
		format   OutputFormat
		expected string
	}{
		{"file", JSON, "file.json"},
		{"file.json", JSON, "file.json"},
		{"file.JSON", JSON, "file.JSON"},
		{"file.txt", JSON, "file.txt.json"}, // This might be unexpected but that's current behavior
		{"dir/file", YAML, "dir/file.yaml"},
		{"dir/file.yaml", YAML, "dir/file.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			formatCfg := FormatCfg{
				Type:       tt.format,
				OutputPath: tt.path,
			}
			result := formatCfg.GetOutputPath()
			if result != tt.expected {
				t.Errorf("GetOutputPath(%s, %s) = %s, want %s",
					tt.path, tt.format, result, tt.expected)
			}
		})
	}
}

// Helper function to compare configs
func compareConfig(t *testing.T, actual, expected *Config) {
	t.Helper()

	if actual.Path != expected.Path {
		t.Errorf("Path = %s, want %s", actual.Path, expected.Path)
	}
	if actual.MaxDepth != expected.MaxDepth {
		t.Errorf("MaxDepth = %d, want %d", actual.MaxDepth, expected.MaxDepth)
	}
	if actual.IncludeFiles != expected.IncludeFiles {
		t.Errorf("IncludeFiles = %v, want %v", actual.IncludeFiles, expected.IncludeFiles)
	}
	if actual.FollowLinks != expected.FollowLinks {
		t.Errorf("FollowLinks = %v, want %v", actual.FollowLinks, expected.FollowLinks)
	}

	// Compare slices
	if !equalStringSlices(actual.ExcludePaths, expected.ExcludePaths) {
		t.Errorf("ExcludePaths = %v, want %v", actual.ExcludePaths, expected.ExcludePaths)
	}
	if !equalStringSlices(actual.ExcludeTypes, expected.ExcludeTypes) {
		t.Errorf("ExcludeTypes = %v, want %v", actual.ExcludeTypes, expected.ExcludeTypes)
	}

	// Compare format config
	if actual.Format.Type != expected.Format.Type {
		t.Errorf("Format.Type = %s, want %s", actual.Format.Type, expected.Format.Type)
	}
	if actual.Format.OutputPath != expected.Format.OutputPath {
		t.Errorf("Format.OutputPath = %s, want %s", actual.Format.OutputPath, expected.Format.OutputPath)
	}
	if actual.Format.Indent != expected.Format.Indent {
		t.Errorf("Format.Indent = %d, want %d", actual.Format.Indent, expected.Format.Indent)
	}
	if !equalStringSlices(actual.Format.ExcludeNodeFields, expected.Format.ExcludeNodeFields) {
		t.Errorf("Format.ExcludeNodeFields = %v, want %v",
			actual.Format.ExcludeNodeFields, expected.Format.ExcludeNodeFields)
	}
}

// Helper function to compare string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
