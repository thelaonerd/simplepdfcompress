package compression

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CompressionOptions holds configuration for the compression job
type CompressionOptions struct {
	Quality string // e.g. /screen, /ebook, /printer, /prepress, /default
}

// CompressPDF compresses a single PDF file using ps2pdf
func CompressPDF(inputPath, outputPath string, opts CompressionOptions) (int64, int64, error) {
	// 1. Get initial file size
	info, err := os.Stat(inputPath)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to stat input file: %w", err)
	}
	initialSize := info.Size()

	// 2. Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return initialSize, 0, fmt.Errorf("failed to create output directory: %w", err)
	}

	// 3. Construct Ghostscript command
	// We call gs directly for better cross-platform support (windows differs from linux/mac)
	bin := GetGhostscriptCommand()
	args := []string{
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		fmt.Sprintf("-sOutputFile=%s", outputPath),
	}

	if opts.Quality != "" {
		// Ghostscript requires / prefix for string constants like /ebook
		args = append(args, fmt.Sprintf("-dPDFSETTINGS=/%s", opts.Quality))
	}
	args = append(args, inputPath)

	cmd := exec.Command(bin, args...)

	// 4. Execute blocking command
	if output, err := cmd.CombinedOutput(); err != nil {
		return initialSize, 0, fmt.Errorf("ps2pdf failed: %v, output: %s", err, string(output))
	}

	// 5. Get final file size
	info, err = os.Stat(outputPath)
	if err != nil {
		return initialSize, 0, fmt.Errorf("failed to stat output file: %w", err)
	}
	finalSize := info.Size()

	return initialSize, finalSize, nil
}
