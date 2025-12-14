package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// CheckResult holds the result of the system and dependency checks
type CheckResult struct {
	OS             string
	Distro         string // Only for Linux
	PackageManager string // Suggested package manager command
	HasGS          bool
	IsReady        bool
	Message        string
}

// PerformChecks identifies the OS and checks for required dependencies
func PerformChecks() CheckResult {
	result := CheckResult{
		OS: runtime.GOOS,
	}

	// 1. Identify OS & Distro
	if result.OS == "linux" {
		result.Distro, result.PackageManager = getLinuxDistroInfo()
	}

	// 2. Check Dependencies
	bin := getGSBinaryName()
	result.HasGS = checkCommand(bin)

	// 3. Formulate Message & Readiness
	if result.HasGS {
		result.IsReady = true
		result.Message = "All dependencies are satisfied. Application is ready."
	} else {
		result.IsReady = false
		result.Message = buildMissingDependencyMessage(result)
	}

	return result
}

func checkCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func getLinuxDistroInfo() (distro string, pkgMgr string) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return "unknown", "unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			distro = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		}
	}

	switch distro {
	case "ubuntu", "debian":
		pkgMgr = "sudo apt install"
	case "fedora":
		pkgMgr = "sudo dnf install"
	case "arch":
		pkgMgr = "sudo pacman -S"
	default:
		pkgMgr = "install" // Generic fallback
	}

	return distro, pkgMgr
}

func buildMissingDependencyMessage(r CheckResult) string {
	var missing []string
	if !r.HasGS {
		missing = append(missing, "Ghostscript")
	}

	msg := fmt.Sprintf("Missing dependencies: %s.\n", strings.Join(missing, ", "))

	if r.OS == "linux" {
		msg += fmt.Sprintf("Please run: %s ghostscript", r.PackageManager)
	} else if r.OS == "darwin" {
		msg += "Please run: brew install ghostscript"
	} else if r.OS == "windows" {
		msg += "Please download and install Ghostscript (gswin64c) from the official website."
	} else {
		msg += "Please install Ghostscript for your operating system."
	}

	return msg
}

func getGSBinaryName() string {
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("gswin64c"); err == nil {
			return "gswin64c"
		}
		if _, err := exec.LookPath("gswin32c"); err == nil {
			return "gswin32c"
		}
	}
	return "gs"
}
