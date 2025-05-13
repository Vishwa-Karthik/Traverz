package output

import (
	"fmt"
	"os"

	"github.com/Vishwa-Karthik/traverz/core"
)

// Writer handles writing output to terminal or file.
type Writer struct{}

// NewWriter creates a new Writer.
func NewWriter() *Writer {
	return &Writer{}
}

// Write outputs the content.
// If outputPath is empty, prints to terminal (with pagination if enabled).
// Otherwise, writes to the specified file.
func (w *Writer) Write(content string, outputPath string, paginate bool) error {
	if outputPath == "" {
		if paginate {
			return core.ShowWithPager(content)
		}
		_, err := fmt.Print(content)
		return err
	}

	err := os.WriteFile(outputPath, []byte(content), 0644)
	if err == nil {
		fmt.Printf("Output successfully written to %s\n", outputPath)
	}
	return err
}
