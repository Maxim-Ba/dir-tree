package configs

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ParseConfig parses configuration from command line flags and/or config file
func ParseConfig() (*Config, error) {
	var configPath string
	var path string
	var outputFormat string
	var outputPath string
	var excludePaths string
	var maxDepth int
	var includeFiles bool
	var followLinks bool
	var excludeTypes string
	var excludeNodeFields string
	
	// Command line flags
	flag.StringVar(&configPath, "c", "", "Path to config file")
	flag.StringVar(&path, "p", ".", "Target directory path")
	flag.StringVar(&outputFormat, "f", "json", "Output format (json, yaml, xml, txt)")
	flag.StringVar(&outputPath, "o", "output-dir", "Output file path")
	flag.BoolVar(&includeFiles, "if", true, "Include files in output")
	flag.BoolVar(&followLinks, "fl", false, "Follow symbolic links")
	flag.StringVar(&excludePaths, "ep", ".git", "Exclude paths (regex patterns, comma separated)")
	flag.StringVar(&excludeTypes, "et", "", "Exclude types (file extensions, comma separated)")
	flag.IntVar(&maxDepth, "d", 1, "Maximum tree depth")
	flag.StringVar(&excludeNodeFields, "enf", "size,is_hidden,type,path", "Exclude node fields from output (comma separated)")
	flag.Parse()

	// Parse comma-separated strings into slices
	excludePathsSlice := parseCommaSeparated(excludePaths)
	excludeTypesSlice := parseCommaSeparated(excludeTypes)
	excludeNodeFieldsSlice := parseCommaSeparated(excludeNodeFields)

	cfg := &Config{
		Path:         path,
		MaxDepth:     maxDepth,
		ExcludePaths: excludePathsSlice,
		ExcludeTypes: excludeTypesSlice,
		IncludeFiles: includeFiles,
		FollowLinks:  followLinks,
		Format: FormatCfg{
			Type:             OutputFormat(outputFormat),
			OutputPath:       outputPath,
			Indent:           2,
			ExcludeNodeFields: excludeNodeFieldsSlice,
		},
	}

	// Load from config file if specified
	if configPath != "" {
		if err := loadConfigFromFile(configPath, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// parseCommaSeparated parses a comma-separated string into a string slice
func parseCommaSeparated(input string) []string {
	if input == "" {
		return []string{}
	}
	
	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// loadConfigFromFile loads configuration from a file using Viper
func loadConfigFromFile(path string, cfg *Config) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %w", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}

	return nil
}
