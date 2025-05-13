package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mattn/go-isatty"
)

// ShowWithPager displays content using a pager (like 'less') if stdout is a TTY.
// If not a TTY or pager fails, it prints directly to stdout.
func ShowWithPager(content string) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		// Not a TTY, just print
		fmt.Print(content)
		return nil
	}

	// Try to use 'less' as a pager
	pagerCmd := os.Getenv("PAGER")
	if pagerCmd == "" {
		pagerCmd = "less"
	}

	cmd := exec.Command(pagerCmd, "-R") // -R for ANSI color/style passthrough
	stdin, err := cmd.StdinPipe()
	if err != nil {
		// Fallback to direct print if pipe fails
		fmt.Print(content)
		return fmt.Errorf("failed to create stdin pipe for pager: %w, falling back to direct print", err)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr // So user sees errors from less

	if err := cmd.Start(); err != nil {
		// Fallback to direct print if command start fails
		fmt.Print(content)
		return fmt.Errorf("failed to start pager command '%s': %w, falling back to direct print", pagerCmd, err)
	}

	_, err = io.WriteString(stdin, content)
	if err != nil {
		// Best effort to close stdin, then let cmd.Wait() handle cleanup
		_ = stdin.Close()
		_ = cmd.Wait()     // Wait for command to finish to avoid zombie processes
		fmt.Print(content) // Fallback
		return fmt.Errorf("failed to write to pager stdin: %w, falling back to direct print", err)
	}

	err = stdin.Close() // Must close stdin to signal EOF to 'less'
	if err != nil {
		_ = cmd.Wait()
		fmt.Print(content) // Fallback
		return fmt.Errorf("failed to close pager stdin: %w, falling back to direct print", err)
	}

	return cmd.Wait() // Wait for the pager to exit
}
