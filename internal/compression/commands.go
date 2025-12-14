package compression

import (
	"os/exec"
	"runtime"
)

// GetGhostscriptCommand returns the executable name for Ghostscript
// depending on the operating system.
func GetGhostscriptCommand() string {
	if runtime.GOOS == "windows" {
		// Try gswin64c.exe first, then 32, then gs
		if _, err := exec.LookPath("gswin64c"); err == nil {
			return "gswin64c"
		}
		if _, err := exec.LookPath("gswin32c"); err == nil {
			return "gswin32c"
		}
		return "gs" // Fallback
	}
	return "gs"
}
