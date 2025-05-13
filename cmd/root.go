package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Vishwa-Karthik/traverz/core"
	"github.com/Vishwa-Karthik/traverz/output"
	renderer "github.com/Vishwa-Karthik/traverz/render"
	"github.com/Vishwa-Karthik/traverz/utils"
	"github.com/spf13/cobra"
)

var (
	rootPath    string
	outputPath  string
	noIcons     bool
	excludeRaw  []string // Use StringSliceVar for multiple or comma-separated
	maxDepth    int
	paginate    bool
	outputStyle string // Though only "tree" is supported for now
)

var rootCmd = &cobra.Command{
	Use:   "traverz",
	Short: "traverz recursively scans a directory and produces a Markdown-rendered tree layout.",
	Long: `traverz is a developer-centric CLI tool that recursively scans any directory 
and produces a Markdown-rendered tree layout, similar to the Unix tree command. 
It can optionally include file/folder icons and output to the terminal or a file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate and normalize rootPath
		absPath, err := filepath.Abs(rootPath)
		if err != nil {
			return fmt.Errorf("error getting absolute path for %s: %w", rootPath, err)
		}
		if exists, isDir := utils.PathExists(absPath); !exists {
			return fmt.Errorf("path does not exist: %s", absPath)
		} else if !isDir {
			return fmt.Errorf("path is not a directory: %s", absPath)
		}

		// Compile exclusion patterns
		var excludePatterns []*regexp.Regexp
		if len(excludeRaw) > 0 {
			// Handle comma-separated values within each string slice element
			var actualExcludes []string
			for _, e := range excludeRaw {
				actualExcludes = append(actualExcludes, strings.Split(e, ",")...)
			}

			for _, pattern := range actualExcludes {
				trimmedPattern := strings.TrimSpace(pattern)
				if trimmedPattern == "" {
					continue
				}
				re, err := regexp.Compile(trimmedPattern)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Invalid regex pattern '%s', skipping: %v\n", pattern, err)
					continue
				}
				excludePatterns = append(excludePatterns, re)
			}
		}

		// Determine actual max depth (-1 means effectively unlimited)
		currentDepth := 0
		if maxDepth < 0 {
			maxDepth = 1<<31 - 1 // Effectively infinite for practical purposes
		}

		// Core logic
		traverser := core.NewTraverser(excludePatterns, !noIcons, maxDepth)
		rootNode, err := traverser.Traverse(absPath, currentDepth)
		if err != nil {
			return fmt.Errorf("error traversing directory: %w", err)
		}

		if rootNode == nil {
			// This might happen if the root directory itself was excluded or unreadable completely
			fmt.Println("No files or directories to display (root might be excluded or unreadable).")
			return nil
		}

		// Rendering
		markdownRenderer := renderer.NewMarkdownRenderer()
		outputString := markdownRenderer.Render(rootNode, !noIcons)

		// Output
		writer := output.NewWriter()
		err = writer.Write(outputString, outputPath, paginate)
		if err != nil {
			return fmt.Errorf("error writing output: %w", err)
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootPath, "path", "p", ".", "Root directory to traverse")
	rootCmd.PersistentFlags().StringVarP(&outputPath, "output", "o", "", "File path to write the output to (if omitted, print to terminal)")
	rootCmd.PersistentFlags().BoolVar(&noIcons, "no-icons", false, "Disable 📁 and 📄 icons")
	// Cobra's StringSliceVar handles multiple uses. It also handles comma-separated values if a single flag is provided.
	rootCmd.PersistentFlags().StringSliceVarP(&excludeRaw, "exclude", "e", []string{}, "Regex pattern(s) to exclude files/folders (e.g., '\\.git', 'node_modules', or '.git,node_modules')")
	rootCmd.PersistentFlags().IntVarP(&maxDepth, "depth", "d", -1, "Max depth of directory traversal (-1 for unlimited)")
	rootCmd.PersistentFlags().BoolVar(&paginate, "paginate", false, "Paginate large outputs if displayed in terminal")
	rootCmd.PersistentFlags().StringVar(&outputStyle, "style", "tree", "Output style (currently only 'tree' is supported)")
	// Help flag is automatically added by Cobra
}
