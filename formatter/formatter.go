package formatter

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/Maxim-Ba/dir-tree/configs"
	"github.com/Maxim-Ba/dir-tree/tree"
	"gopkg.in/yaml.v2"
)

// Format converts a tree node to the specified output format
func Format(tree *tree.Node, cfg *configs.FormatCfg) ([]byte, error) {
	switch cfg.Type {
	case configs.JSON:
		return formatJSON(tree, cfg)
	case configs.YAML:
		return formatYAML(tree, cfg)
	case configs.XML:
		return formatXML(tree, cfg)
	case configs.TXT:
		return formatTXT(tree, 0, cfg), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", cfg.Type)
	}
}

// filteredNode represents a node with filtered fields for output
type filteredNode struct {
	Name     string          `json:"name,omitempty" yaml:"name,omitempty" xml:"name,omitempty"`
	Path     string          `json:"path,omitempty" yaml:"path,omitempty" xml:"path,omitempty"`
	Type     tree.FileType   `json:"type,omitempty" yaml:"type,omitempty" xml:"type,omitempty"`
	Size     int64           `json:"size,omitempty" yaml:"size,omitempty" xml:"size,omitempty"`
	Children []*filteredNode `json:"children,omitempty" yaml:"children,omitempty" xml:"children>node,omitempty"`
	IsHidden bool            `json:"is_hidden,omitempty" yaml:"is_hidden,omitempty" xml:"is_hidden,omitempty"`
}

// createFilteredNode creates a filtered node with excluded fields removed
func createFilteredNode(node *tree.Node, excludeFields []string) *filteredNode {
	if node == nil {
		return nil
	}

	filtered := &filteredNode{}

	// Copy only non-excluded fields
	if !contains(excludeFields, "name") {
		filtered.Name = node.Name
	}
	if !contains(excludeFields, "path") {
		filtered.Path = node.Path
	}
	if !contains(excludeFields, "type") {
		filtered.Type = node.Type
	}
	if !contains(excludeFields, "size") {
		filtered.Size = node.Size
	}
	if !contains(excludeFields, "is_hidden") {
		filtered.IsHidden = node.IsHidden
	}

	// Recursively process children (if children field is not excluded)
	if !contains(excludeFields, "children") && node.Children != nil {
		filtered.Children = make([]*filteredNode, 0, len(node.Children))
		for _, child := range node.Children {
			filteredChild := createFilteredNode(child, excludeFields)
			if filteredChild != nil {
				filtered.Children = append(filtered.Children, filteredChild)
			}
		}
	}

	return filtered
}

// contains checks if a string slice contains a specific item
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// formatJSON formats the tree as JSON
func formatJSON(node *tree.Node, cfg *configs.FormatCfg) ([]byte, error) {
	var data interface{} = node

	// Apply field filtering if needed
	if len(cfg.ExcludeNodeFields) > 0 {
		data = createFilteredNode(node, cfg.ExcludeNodeFields)
	}

	if cfg.Indent > 0 {
		return json.MarshalIndent(data, "", strings.Repeat(" ", cfg.Indent))
	}
	return json.Marshal(data)
}

// formatYAML formats the tree as YAML
func formatYAML(node *tree.Node, cfg *configs.FormatCfg) ([]byte, error) {
	var data interface{} = node

	// Apply field filtering if needed
	if len(cfg.ExcludeNodeFields) > 0 {
		data = createFilteredNode(node, cfg.ExcludeNodeFields)
	}

	return yaml.Marshal(data)
}

// formatXML formats the tree as XML
func formatXML(node *tree.Node, cfg *configs.FormatCfg) ([]byte, error) {
	var data interface{} = node

	// Apply field filtering if needed
	if len(cfg.ExcludeNodeFields) > 0 {
		data = createFilteredNode(node, cfg.ExcludeNodeFields)
	}

	return xml.MarshalIndent(data, "", "  ")
}

// formatTXT formats the tree as plain text with visual indicators
func formatTXT(node *tree.Node, level int, cfg *configs.FormatCfg) []byte {
	var result strings.Builder
	indent := strings.Repeat("  ", level)

	// Build line parts based on included fields
	parts := []string{}

	// Add prefix for type indication
	prefix := "📁 " // directory
	if node.Type == tree.File {
		prefix = "📄 " // file
	} else if node.Type == tree.Symlink {
		prefix = "🔗 " // symlink
	}
	parts = append(parts, prefix)

	// Add name (if not excluded)
	if !contains(cfg.ExcludeNodeFields, "name") {
		parts = append(parts, node.Name)
	}

	// Add size (if not excluded and if file with size > 0)
	if !contains(cfg.ExcludeNodeFields, "size") && node.Type == tree.File && node.Size > 0 {
		parts = append(parts, fmt.Sprintf("(%d bytes)", node.Size))
	}

	// Add hidden status (if not excluded and file is hidden)
	if !contains(cfg.ExcludeNodeFields, "is_hidden") && node.IsHidden {
		parts = append(parts, "[hidden]")
	}

	result.WriteString(fmt.Sprintf("%s%s\n", indent, strings.Join(parts, " ")))

	// Recursively process children (if children field is not excluded)
	if !contains(cfg.ExcludeNodeFields, "children") {
		for _, child := range node.Children {
			result.Write(formatTXT(child, level+1, cfg))
		}
	}

	return []byte(result.String())
}
