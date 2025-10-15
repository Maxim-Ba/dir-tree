// Package dirtree provides APIs for generating directory trees programmatically
package dirtree

import (
	"fmt"
	"os"

	"github.com/Maxim-Ba/dir-tree/configs"
	"github.com/Maxim-Ba/dir-tree/formatter"
	"github.com/Maxim-Ba/dir-tree/tree"
)

// Generate creates a directory tree based on the provided configuration
func Generate(cfg *configs.Config) ([]byte, error) {
	root, err := tree.BuildTree(
		tree.BuildOptions{Path: cfg.Path,
		MaxDepth:     cfg.MaxDepth,
		ExcludePaths: cfg.ExcludePaths,
		ExcludeTypes: cfg.ExcludePaths,
		IncludeFiles: cfg.IncludeFiles,
		FollowLinks:  cfg.FollowLinks,
	})
	if err != nil {
		return nil, err
	}
	return formatter.Format(root, &cfg.Format)
}

// GenerateToFile generates a directory tree and saves it to a file
func GenerateToFile(cfg *configs.Config) error {
	data, err := Generate(cfg)
	if err != nil {
		return err
	}

	outputPath := cfg.Format.GetOutputPath()
	if outputPath == "" {
		return fmt.Errorf("output path is required for file generation")
	}

	return os.WriteFile(outputPath, data, 0644)
}

// GenerateJSON quickly generates a JSON directory tree (convenience method)
func GenerateJSON(path string, maxDepth int) ([]byte, error) {
	cfg := configs.New().WithPath(path).Build()
	cfg.MaxDepth = maxDepth
	cfg.Format.Type = configs.JSON
	return Generate(cfg)
}

// GenerateASCII quickly generates a text directory tree (convenience method)
func GenerateASCII(path string, maxDepth int) (string, error) {
	cfg := configs.New().WithPath(path).Build()
	cfg.MaxDepth = maxDepth
	cfg.Format.Type = configs.TXT
	data, err := Generate(cfg)
	return string(data), err
}
