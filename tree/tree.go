package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileType represents the type of a file system node
type FileType string

const (
	Directory FileType = "directory"
	File      FileType = "file"
	Symlink   FileType = "symlink"
)

// Node represents a file system node in the directory tree
type Node struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Type     FileType `json:"type"`
	Size     int64    `json:"size,omitempty"`
	Children []*Node  `json:"children,omitempty"`
	IsHidden bool     `json:"is_hidden,omitempty"`
}
type BuildOptions struct {
	Path         string
	MaxDepth     int
	ExcludePaths []string
	ExcludeTypes []string
	IncludeFiles bool
	FollowLinks  bool
}

// BuildTree constructs a directory tree from the given options
func BuildTree(opts BuildOptions) (*Node, error) {
	info, err := os.Stat(opts.Path)
	if err != nil {
		return nil, fmt.Errorf("error accessing path %s: %w", opts.Path, err)
	}

	return buildTreeRecursive(opts.Path, info, &opts, 0)
}

// buildTreeRecursive recursively builds the directory tree
func buildTreeRecursive(currentPath string, info os.FileInfo, opts *BuildOptions, currentDepth int) (*Node, error) {
	// Check depth limit
	if opts.MaxDepth != -1 && currentDepth > opts.MaxDepth {
		return nil, nil
	}

	// Check path exclusions
	if isExcludedPath(currentPath, opts.ExcludePaths) {
		return nil, nil
	}

	node := &Node{
		Name: info.Name(),
		Path: currentPath,
	}

	// Determine node type and set size
	if info.IsDir() {
		node.Type = Directory
		node.Size = 0 // Directories have size 0 or could calculate total size
	} else if info.Mode()&os.ModeSymlink != 0 {
		node.Type = Symlink
		node.Size = info.Size()

		// Handle symlinks if following is enabled
		if opts.FollowLinks {
			targetPath, err := filepath.EvalSymlinks(currentPath)
			if err == nil {
				targetInfo, err := os.Stat(targetPath)
				if err == nil {
					if targetInfo.IsDir() {
						node.Type = Directory
						node.Size = 0
					} else {
						node.Type = File
						node.Size = targetInfo.Size()
					}
				}
			}
		}
	} else {
		node.Type = File
		node.Size = info.Size()
	}

	// Check type exclusions
	if node.Type == File && isExcludedType(currentPath, opts.ExcludeTypes) {
		return nil, nil
	}

	// Check if file is hidden
	node.IsHidden = isHiddenFile(info.Name())

	// If directory (or symlink to directory with followLinks), process children
	if node.Type == Directory {
		var entries []os.DirEntry
		var err error

		// Get directory entries
		if opts.FollowLinks && node.Type == Symlink {
			// For symlinks, get contents of target directory
			targetPath, err := filepath.EvalSymlinks(currentPath)
			if err == nil {
				entries, err = os.ReadDir(targetPath)
				if err != nil {
					return nil, fmt.Errorf("error reading directory %s: %w", targetPath, err)
				}
			}
		} else {
			entries, err = os.ReadDir(currentPath)
		}

		if err != nil {
			return nil, fmt.Errorf("error reading directory %s: %w", currentPath, err)
		}

		for _, entry := range entries {
			entryInfo, err := entry.Info()
			if err != nil {
				continue // Skip problematic entries
			}

			fullPath := filepath.Join(currentPath, entryInfo.Name())

			// Skip files if not included
			if !opts.IncludeFiles && !entryInfo.IsDir() {
				continue
			}

			child, err := buildTreeRecursive(fullPath, entryInfo, opts, currentDepth+1)
			if err != nil {
				return nil, err
			}
			if child != nil {
				node.Children = append(node.Children, child)
			}
		}
	}

	return node, nil
}

// isExcludedPath checks if a path matches any exclusion patterns
func isExcludedPath(path string, excludePatterns []string) bool {
	for _, pattern := range excludePatterns {
		matched, err := regexp.MatchString(pattern, path)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// isExcludedType checks if a file type should be excluded
func isExcludedType(path string, excludeTypes []string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, excludedExt := range excludeTypes {
		if strings.ToLower(excludedExt) == ext {
			return true
		}
	}
	return false
}

// isHiddenFile checks if a file is hidden (starts with dot)
func isHiddenFile(name string) bool {
	return strings.HasPrefix(name, ".")
}
