# Traverz

`traverz` is a developer-centric, cross-platform CLI tool that recursively scans any directory and produces a Markdown-rendered tree layout, much like the Unix `tree` command, but with optional file/folder icons and direct Markdown output suitable for documentation or sharing. The output can be displayed in the terminal or saved to a file.

Designed with modern development workflows in mind, `traverz` aims to be an intuitive and useful utility for visualizing directory structures.

## Key Features

*   **Recursive Directory Traversal:** Scans directories and their subdirectories.
*   **Markdown Tree Output:** Generates a clean, hierarchical tree structure formatted in Markdown.
*   **Optional Icons:** Enhances visual distinction with 📁 (folder) and 📄 (file) icons, which can be toggled off for plain output.
*   **Flexible Output:**
    *   Print directly to the terminal.
    *   Save the output to a specified file.
*   **Exclusion Filters:** Uses regular expressions to exclude specific files or folders from the scan (e.g., `.git`, `node_modules`, `*.log`). Supports multiple patterns.
*   **Depth Control:** Limits the maximum depth of directory traversal.
*   **Terminal Pagination:** For large directory structures displayed in the terminal, `traverz` can paginate the output using your system's default pager (like `less`).
*   **Cross-Platform Compatibility:** Works seamlessly on Windows, macOS, and Linux.
*   **Graceful Error Handling:** Skips unreadable directories with warnings instead of crashing.

## Example Output

**Default (with Icons):**

```markdown
📁 traverz/
├── 📁 cmd/
│   └── 📄 root.go
├── 📁 core/
│   ├── 📄 traverser.go
│   ├── 📄 filter.go
│   ├── 📄 icon.go
│   └── 📄 paginator.go
├── 📁 renderer/
│   └── 📄 markdown.go
├── 📁 output/
│   └── 📄 writer.go
├── 📁 utils/
│   └── 📄 fsutils.go
├── 📄 main.go
└── 📄 go.mod
```

**(without Icons):**
```markdown
traverz/
├── cmd/
│   └── root.go
├── core/
│   ├── traverser.go
│   ├── filter.go
│   ├── icon.go
│   └── paginator.go
├── renderer/
│   └── markdown.go
├── output/
│   └── writer.go
├── utils/
│   └── fsutils.go
├── main.go
└── go.mod
```

## Installation
+ There are a couple of ways to install traverz:

1. Option 1:
Using `go install` (Recommended if you have Go) If you have Go (version 1.20 or newer recommended) installed and your $GOPATH/bin (or $HOME/go/bin) is in your system's PATH:

```bash
go install github.com/Vishwa-Karthik/traverz@latest
```

Use code with caution.
Bash
This command will download the source, compile it, and place the traverz executable in your Go binary directory.
2. Option 2: Pre-compiled Binaries (For all users)
You can download pre-compiled binaries for your specific operating system and architecture from the GitHub Releases Page. (Note: You'll need to create releases on GitHub for this link to be active).

Download the appropriate archive (.zip for Windows, .tar.gz for Linux/macOS).
Extract the traverz (or traverz.exe for Windows) executable.

3. Place the executable in a directory that is part of your system's PATH environment variable (e.g., /usr/local/bin or ~/bin on Linux/macOS, or a custom folder added to PATH on Windows).

4. ensure the binary is executable (e.g., chmod +x traverz on Linux/macOS).

## 🛠️ Usage
The basic command structure for traverz is:
```bash
traverz [flags]
```

`If no path is specified, it defaults to the current directory (.).`


## CLI Arguments

The following flags and options are available:

| Flag / Option       | Shorthand | Description                                                                                                                                                                 | Default Value         |
| :------------------ | :-------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :-------------------- |
| `--path`            | `-p`      | Specifies the root directory to traverse.                                                                                                                                   | `.` (current directory) |
| `--output`          | `-o`      | Specifies the file path to write the Markdown output to. If omitted, output is printed to the terminal.                                                                   | (none)                |
| `--no-icons`        |           | Disables the display of 📁 and 📄 icons, resulting in a plain text tree.                                                                                                       | `false`               |
| `--exclude`         | `-e`      | A regular expression pattern to exclude files or folders. This flag can be used multiple times for multiple patterns, or patterns can be comma-separated (e.g., `-e ".git,pubspec.lock,node_modules"`). | (none)                |
| `--depth`           | `-d`      | Sets the maximum depth of directory traversal. A value of `0` means only the root directory's immediate contents (files/folders at level 1). `-1` means unlimited depth.   | `-1`                  |
| `--paginate`        |           | If outputting to the terminal, this flag will attempt to paginate large outputs using the system's pager (e.g., `less`).                                                     | `false`               |
| `--style`           |           | Defines the output style. Currently, only `tree` is supported, which produces the Markdown tree layout.                                                                       | `tree`                |
| `--help`            | `-h`      | Displays this help message listing all available arguments and exits.                                                                                                       |                       |

## Practical Examples
1. Scan the current directory and print to terminal (default behavior):
```bash 
traverz
```

2. Scan a specific project directory:
```bash 
traverz --path ./my-awesome-project
```

or using the shorthand:
```bash 
traverz -p ./my-awesome-project
```

3. Scan the current directory and save the output to project_structure.md:
```bash 
traverz --output project_structure.md
```

or
```bash 
traverz -o project_structure.md
```

4. Scan with no icons and limit traversal to a depth of 2:
```bash 
traverz --path /var/log --no-icons --depth 2
```

5. Exclude .git directory and all node_modules folders:
```bash 
traverz --exclude "\.git" --exclude "node_modules"
```

(Note: Shells might interpret . in regex, so escaping \. can be safer for literal dots. Quotes help prevent shell expansion of special characters.)

6. Alternatively, using comma-separated values:
```bash 
traverz --exclude "\.git,node_modules,target"
```

7. Scan a large directory and paginate the output in the terminal:
```bash 
traverz --path /usr/local/lib --paginate --depth 4
```

## Platform Compatibility
`traverz` is designed to be fully cross-platform and has been tested on:
+ `Windows` (CMD, PowerShell, Windows Terminal)
+ `macOS` (Terminal, iTerm2)
+ `Linux` (various distributions using bash, zsh, fish, etc.)
It handles platform-specific path separators correctly using Go's `path/filepath` package.

## Building from Source (For Developers)
If you wish to build traverz from source, for example, to contribute or to use unreleased features:
1. Clone the repository:
```bash
git clone https://github.com/Vishwa-Karthik/traverz.git
cd traverz
```
2. Ensure `Go` is installed: You'll need Go (version 1.20+ recommended).
3. Build the executable:
```bash
go build -o traverz .
```

This will create a traverz (or `traverz.exe` on Windows) executable in the current directory.

4. (Optional) Cross-compilation:
Go makes cross-compilation straightforward. Here are a few examples:
+ For Windows (64-bit)
```bash 
GOOS=windows GOARCH=amd64 go build -o traverz.exe .
```

+ For Linux (64-bit)
```bash
GOOS=linux GOARCH=amd64 go build -o traverz .
```

+ For macOS (Apple Silicon ARM64)
```bash
GOOS=darwin GOARCH=arm64 go build -o traverz .
```

+ For macOS (Intel AMD64)
```bash
GOOS=darwin GOARCH=amd64 go build -o traverz .
```

## Contributing

Contributions are highly welcome! Whether it's reporting a bug, suggesting a feature, improving documentation, or submitting a pull request, your input is valued.
Please feel free to:

+ Open an Issue to report bugs or discuss new features.
+ Fork the repository, make your changes, and submit a Pull Request.
+ When contributing code, please try to follow existing code style and ensure your changes are well-tested.