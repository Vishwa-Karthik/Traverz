package renderer

import (
	"fmt"
	"strings"

	"github.com/Vishwa-Karthik/traverz/core"
)

// MarkdownRenderer renders the node tree into a Markdown string.
type MarkdownRenderer struct{}

// NewMarkdownRenderer creates a new MarkdownRenderer.
func NewMarkdownRenderer() *MarkdownRenderer {
	return &MarkdownRenderer{}
}

// Render converts the tree of Nodes into a Markdown string.
func (r *MarkdownRenderer) Render(rootNode *core.Node, showIcons bool) string {
	if rootNode == nil {
		return ""
	}
	var sb strings.Builder

	// Print root node itself
	rootIcon := ""
	if rootNode.IsDir { // Root is typically a directory for traversal
		rootIcon = core.GetIcon(true, showIcons)
	}
	// The example output suggests the root directory ends with a slash.
	sb.WriteString(fmt.Sprintf("%s%s/\n", rootIcon, rootNode.Name))

	// Then print its children
	r.renderChildren(&sb, rootNode.Children, "", showIcons)
	return sb.String()
}

func (r *MarkdownRenderer) renderChildren(sb *strings.Builder, children []*core.Node, prefix string, showIcons bool) {
	for i, child := range children {
		isLast := (i == len(children)-1)

		connector := "├── "
		if isLast {
			connector = "└── "
		}

		iconStr := core.GetIcon(child.IsDir, showIcons)
		// Add comment if present (future feature)
		// comment := ""
		// if child.Comment != "" { comment = " # " + child.Comment }
		// sb.WriteString(fmt.Sprintf("%s%s%s%s%s\n", prefix, connector, iconStr, child.Name, comment))
		sb.WriteString(fmt.Sprintf("%s%s%s%s\n", prefix, connector, iconStr, child.Name))

		if child.IsDir && len(child.Children) > 0 {
			newPrefix := prefix
			if isLast {
				newPrefix += "    " // Four spaces for alignment
			} else {
				newPrefix += "│   " // Vertical bar and three spaces
			}
			r.renderChildren(sb, child.Children, newPrefix, showIcons)
		}
	}
}
