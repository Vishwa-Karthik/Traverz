package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Node represents a file or directory in the tree
type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Children []*Node
	// Level int // Not strictly needed if traversal handles depth
}

// Traverser handles the logic of walking the directory tree
type Traverser struct {
	excludePatterns []*regexp.Regexp
	showIcons       bool // This might be better handled by the renderer
	maxDepth        int
}

// NewTraverser creates a new Traverser instance
func NewTraverser(excludePatterns []*regexp.Regexp, showIcons bool, maxDepth int) *Traverser {
	return &Traverser{
		excludePatterns: excludePatterns,
		showIcons:       showIcons, // Kept for now, though renderer is primary user
		maxDepth:        maxDepth,
	}
}

// Traverse recursively scans the given path and builds a tree of Nodes
func (t *Traverser) Traverse(rootPath string, currentDepth int) (*Node, error) {
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %w", rootPath, err)
	}

	// Check if the root path itself should be excluded (by name)
	if IsExcluded(filepath.Base(absPath), t.excludePatterns) {
		return nil, nil // Root itself is excluded
	}

	info, err := os.Stat(absPath)
	if err != nil {
		// If the root path itself is unreadable, we can't proceed for this path
		fmt.Fprintf(os.Stderr, "Warning: Could not stat root path %s: %v. Skipping.\n", absPath, err)
		return nil, nil // Or return err if this should halt operation
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", absPath)
	}

	rootNode := &Node{
		Name:  filepath.Base(absPath),
		Path:  absPath,
		IsDir: true,
	}

	if currentDepth >= t.maxDepth {
		return rootNode, nil // Max depth reached, return the directory node without children
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		// Log warning but continue if possible, this directory's children won't be listed
		fmt.Fprintf(os.Stderr, "Warning: Could not read directory %s: %v. Skipping contents.\n", absPath, err)
		return rootNode, nil // Return the directory node itself, but empty
	}

	// Sort entries: directories first, then files, all alphabetically
	sort.SliceStable(entries, func(i, j int) bool {
		infoI, errI := entries[i].Info()
		infoJ, errJ := entries[j].Info()

		// Handle errors during Info() call, e.g., due to symlink issues or permissions
		if errI != nil && errJ == nil {
			return false
		} // Place entries with errors last or handle as per preference
		if errI == nil && errJ != nil {
			return true
		}
		if errI != nil && errJ != nil {
			return entries[i].Name() < entries[j].Name()
		} // Sort by name if both error

		isDirI := infoI.IsDir()
		isDirJ := infoJ.IsDir()

		if isDirI && !isDirJ {
			return true // Directories first
		}
		if !isDirI && isDirJ {
			return false // Files after directories
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name()) // Alphabetical for same type
	})

	for _, entry := range entries {
		entryName := entry.Name()
		entryPath := filepath.Join(absPath, entryName)

		if IsExcluded(entryName, t.excludePatterns) {
			continue
		}

		var entryInfo fs.FileInfo
		entryInfo, err = entry.Info() // Get FileInfo
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not get info for %s: %v. Skipping.\n", entryPath, err)
			continue
		}

		if entryInfo.IsDir() {
			childDirNode, err := t.Traverse(entryPath, currentDepth+1)
			if err != nil {
				// Error already logged by the recursive call or stat
				// fmt.Fprintf(os.Stderr, "Warning: Error traversing subdirectory %s: %v. Skipping.\n", entryPath, err)
				continue
			}
			if childDirNode != nil { // It might be nil if excluded or unreadable
				rootNode.Children = append(rootNode.Children, childDirNode)
			}
		} else {
			fileNode := &Node{
				Name:  entryName,
				Path:  entryPath,
				IsDir: false,
			}
			rootNode.Children = append(rootNode.Children, fileNode)
		}
	}
	return rootNode, nil
}
