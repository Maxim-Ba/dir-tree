package tree

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsHiddenFile tests the isHiddenFile function
func TestIsHiddenFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Starts with dot", ".hidden", true},
		{"Does not start with dot", "visible", false},
		{"Empty string", "", false},
		{"Only dot", ".", true},
		{"Multiple dots", "..", true},
		{"Dot in middle", "file.txt", false},
		{"Dot at start", ".file.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isHiddenFile(tt.input)
			if result != tt.expected {
				t.Errorf("isHiddenFile(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsExcludedType tests the isExcludedType function
func TestIsExcludedType(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		excludeTypes  []string
		expected      bool
	}{
		{
			name:          "Excluded extension",
			path:          "file.go",
			excludeTypes:  []string{".go", ".txt"},
			expected:      true,
		},
		{
			name:          "Not excluded extension",
			path:          "file.py",
			excludeTypes:  []string{".go", ".txt"},
			expected:      false,
		},
		{
			name:          "Case insensitive",
			path:          "file.GO",
			excludeTypes:  []string{".go"},
			expected:      true,
		},
		{
			name:          "No extension",
			path:          "README",
			excludeTypes:  []string{".go"},
			expected:      false,
		},
		{
			name:          "Empty exclude list",
			path:          "file.go",
			excludeTypes:  []string{},
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isExcludedType(tt.path, tt.excludeTypes)
			if result != tt.expected {
				t.Errorf("isExcludedType(%q, %v) = %v, want %v", 
					tt.path, tt.excludeTypes, result, tt.expected)
			}
		})
	}
}

// TestIsExcludedPath tests the isExcludedPath function
func TestIsExcludedPath(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		excludePatterns []string
		expected        bool
	}{
		{
			name:            "Matched pattern",
			path:            "/home/user/test",
			excludePatterns: []string{".*test.*"},
			expected:        true,
		},
		{
			name:            "No pattern match",
			path:            "/home/user/docs",
			excludePatterns: []string{".*test.*"},
			expected:        false,
		},
		{
			name:            "Multiple patterns",
			path:            "/home/user/temp",
			excludePatterns: []string{".*test.*", ".*temp.*"},
			expected:        true,
		},
		{
			name:            "Invalid pattern",
			path:            "/home/user/test",
			excludePatterns: []string{"["}, // Invalid regex
			expected:        false,
		},
		{
			name:            "Empty patterns",
			path:            "/home/user/test",
			excludePatterns: []string{},
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isExcludedPath(tt.path, tt.excludePatterns)
			if result != tt.expected {
				t.Errorf("isExcludedPath(%q, %v) = %v, want %v", 
					tt.path, tt.excludePatterns, result, tt.expected)
			}
		})
	}
}

// TestBuildTree tests the BuildTree function
func TestBuildTree(t *testing.T) {
	// Create temporary test directory structure
	tmpDir := t.TempDir()
	
	// Create test files and directories
	dirs := []string{
		"dir1",
		"dir1/subdir1",
		"dir2",
	}
	
	files := []string{
		"file1.txt",
		"file2.go",
		"dir1/file3.txt",
		"dir1/subdir1/file4.go",
		"dir2/file5.txt",
	}
	
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}
	
	for _, file := range files {
		f, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		f.Close()
	}

	tests := []struct {
		name        string
		opts        BuildOptions
		shouldError bool
		validate    func(*testing.T, *Node)
	}{
		{
			name: "Basic structure",
			opts: BuildOptions{
				Path:         tmpDir,
				MaxDepth:     -1,
				IncludeFiles: true,
			},
			shouldError: false,
			validate: func(t *testing.T, root *Node) {
				if root == nil {
					t.Error("Root node should not be nil")
				}
				if root.Type != Directory {
					t.Error("Root should be a directory")
				}
				// Should have 2 directories and 2 files at root level
				if len(root.Children) != 4 {
					t.Errorf("Expected 4 children at root, got %d", len(root.Children))
				}
			},
		},
		{
			name: "With max depth 1",
			opts: BuildOptions{
				Path:         tmpDir,
				MaxDepth:     1,
				IncludeFiles: true,
			},
			shouldError: false,
			validate: func(t *testing.T, root *Node) {
				// Should not see subdir1 contents
				for _, child := range root.Children {
					if child.Name == "dir1" {
						for _, subchild := range child.Children {
							if subchild.Name == "subdir1" {
								if len(subchild.Children) > 0 {
									t.Error("subdir1 should have no children due to max depth")
								}
							}
						}
					}
				}
			},
		},
		{
			name: "Exclude go files",
			opts: BuildOptions{
				Path:         tmpDir,
				MaxDepth:     -1,
				ExcludeTypes: []string{".go"},
				IncludeFiles: true,
			},
			shouldError: false,
			validate: func(t *testing.T, root *Node) {
				var checkForGoFiles func(*Node)
				checkForGoFiles = func(node *Node) {
					if node.Type == File && filepath.Ext(node.Name) == ".go" {
						t.Errorf("Found go file that should be excluded: %s", node.Path)
					}
					for _, child := range node.Children {
						checkForGoFiles(child)
					}
				}
				checkForGoFiles(root)
			},
		},
		{
			name: "Exclude dir1 path",
			opts: BuildOptions{
				Path:         tmpDir,
				MaxDepth:     -1,
				ExcludePaths: []string{".*dir1.*"},
				IncludeFiles: true,
			},
			shouldError: false,
			validate: func(t *testing.T, root *Node) {
				for _, child := range root.Children {
					if child.Name == "dir1" {
						t.Error("dir1 should be excluded")
					}
				}
			},
		},
		{
			name: "Files not included",
			opts: BuildOptions{
				Path:         tmpDir,
				MaxDepth:     -1,
				IncludeFiles: false,
			},
			shouldError: false,
			validate: func(t *testing.T, root *Node) {
				var checkForFiles func(*Node)
				checkForFiles = func(node *Node) {
					if node.Type == File {
						t.Errorf("Found file when IncludeFiles is false: %s", node.Path)
					}
					for _, child := range node.Children {
						checkForFiles(child)
					}
				}
				checkForFiles(root)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildTree(tt.opts)
			
			if tt.shouldError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.shouldError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestBuildTreeError tests error cases for BuildTree
func TestBuildTreeError(t *testing.T) {
	tests := []struct {
		name string
		opts BuildOptions
	}{
		{
			name: "Non-existent path",
			opts: BuildOptions{
				Path: "/non/existent/path",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := BuildTree(tt.opts)
			if err == nil {
				t.Error("Expected error but got none")
			}
			if result != nil {
				t.Error("Expected nil result when error occurs")
			}
		})
	}
}
