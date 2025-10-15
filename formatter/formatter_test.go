package formatter

import (
	"testing"

	"github.com/Maxim-Ba/dir-tree/configs"
	"github.com/Maxim-Ba/dir-tree/tree"
)

// TestCreateFilteredNode tests field exclusion functionality
func TestCreateFilteredNode(t *testing.T) {
	testNode := &tree.Node{
		Name:     "test",
		Path:     "/test",
		Type:     tree.File,
		Size:     100,
		IsHidden: true,
		Children: []*tree.Node{
			{
				Name: "child",
				Path: "/test/child",
				Type: tree.File,
				Size: 50,
			},
		},
	}

	tests := []struct {
		name          string
		excludeFields []string
		checkFields   map[string]bool // field -> should exist
	}{
		{
			name:          "Exclude size field",
			excludeFields: []string{"size"},
			checkFields: map[string]bool{
				"name":     true,
				"path":     true,
				"type":     true,
				"size":     false,
				"is_hidden": true,
				"children": true,
			},
		},
		{
			name:          "Exclude multiple fields",
			excludeFields: []string{"size", "is_hidden", "path"},
			checkFields: map[string]bool{
				"name":     true,
				"path":     false,
				"type":     true,
				"size":     false,
				"is_hidden": false,
				"children": true,
			},
		},
		{
			name:          "Exclude children field",
			excludeFields: []string{"children"},
			checkFields: map[string]bool{
				"name":     true,
				"path":     true,
				"type":     true,
				"size":     true,
				"is_hidden": true,
				"children": false,
			},
		},
		{
			name:          "Exclude all fields",
			excludeFields: []string{"name", "path", "type", "size", "is_hidden", "children"},
			checkFields: map[string]bool{
				"name":     false,
				"path":     false,
				"type":     false,
				"size":     false,
				"is_hidden": false,
				"children": false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := createFilteredNode(testNode, tt.excludeFields)
			
			if filtered == nil {
				t.Fatal("createFilteredNode() returned nil")
			}

			// Check field presence
			if tt.checkFields["name"] && filtered.Name == "" {
				t.Error("Name field should not be empty")
			}
			if tt.checkFields["path"] && filtered.Path == "" {
				t.Error("Path field should not be empty")
			}
			if tt.checkFields["type"] && filtered.Type == "" {
				t.Error("Type field should not be empty")
			}
			if tt.checkFields["size"] && filtered.Size == 0 {
				t.Error("Size field should not be zero")
			}
			if tt.checkFields["is_hidden"] && !filtered.IsHidden {
				t.Error("IsHidden field should be true")
			}
			if tt.checkFields["children"] && len(filtered.Children) == 0 {
				t.Error("Children field should not be empty")
			}

			// Check field absence
			if !tt.checkFields["name"] && filtered.Name != "" {
				t.Error("Name field should be empty")
			}
			if !tt.checkFields["path"] && filtered.Path != "" {
				t.Error("Path field should be empty")
			}
			if !tt.checkFields["type"] && filtered.Type != "" {
				t.Error("Type field should be empty")
			}
			if !tt.checkFields["size"] && filtered.Size != 0 {
				t.Error("Size field should be zero")
			}
			if !tt.checkFields["is_hidden"] && filtered.IsHidden {
				t.Error("IsHidden field should be false")
			}
			if !tt.checkFields["children"] && len(filtered.Children) != 0 {
				t.Error("Children field should be empty")
			}
		})
	}
}


// TestContains tests the helper function
func TestContains(t *testing.T) {
	tests := []struct {
		name string
		slice []string
		item  string
		want  bool
	}{
		{
			name:  "Item exists in slice",
			slice: []string{"a", "b", "c"},
			item:  "b",
			want:  true,
		},
		{
			name:  "Item does not exist in slice",
			slice: []string{"a", "b", "c"},
			item:  "d",
			want:  false,
		},
		{
			name:  "Empty slice",
			slice: []string{},
			item:  "a",
			want:  false,
		},
		{
			name:  "Case sensitive match",
			slice: []string{"Hello", "World"},
			item:  "hello",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.slice, tt.item); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEdgeCases tests edge cases and error conditions
func TestFormat_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		node    *tree.Node
		cfg     *configs.FormatCfg
		wantErr bool
	}{
		{
			name:    "Nil node",
			node:    nil,
			cfg:     &configs.FormatCfg{Type: configs.JSON},
			wantErr: false, // Should handle nil gracefully
		},
		{
			name: "Empty node",
			node: &tree.Node{},
			cfg:  &configs.FormatCfg{Type: configs.JSON},
			wantErr: false,
		},
		{
			name: "Node with nil children",
			node: &tree.Node{
				Name:     "test",
				Children: nil,
			},
			cfg:     &configs.FormatCfg{Type: configs.JSON},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Format(tt.node, tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
