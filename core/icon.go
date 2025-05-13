package core

// GetIcon returns the appropriate icon string for a file or directory.
func GetIcon(isDir bool, showIcons bool) string {
	if !showIcons {
		return ""
	}
	if isDir {
		return "📁 "
	}
	return "📄 "
}
