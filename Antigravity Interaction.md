# Antigravity Interaction Log: SimplePDFCompress

## 1. Project Overview
**SimplePDFCompress** is a desktop utility designed to compress PDF files efficiently using Ghostscript. The project transitioned from a basic Go/Fyne prototype to a production-ready, cross-platform application capable of running on Linux, Windows, and macOS.

**Developer**: thelaonerd
**AI Assistant**: Antigravity (Google Deepmind)
**Date**: December 2025

## 2. Chronological Development Log

### Phase 1: Foundation & Logic
*   **Initial Setup**: Initialized Go module (`com.simplepdfcompress.app`) and installed Fyne dependencies.
*   **Core Logic**: Implemented the `compression` package to wrap `ps2pdf` (Ghostscript).
*   **Concurrency**: Built a worker pool system to handle batch compression without freezing the UI.

### Phase 2: UI Implementation
*   **Single File Tab**: Created a user-friendly interface for single file processing.
*   **Batch Tab**: Implemented a multi-file interface with "Add File" and "Add Folder" capabilities.
*   **Native Experience**:
    *   Replaced Fyne's default file pickers with `zenity` and `sqweek/dialog` to use the operating system's native file dialogs.
    *   Integrated `fc-match` (Linux) to detect and apply the system's default font for a seamless look on KDE/Gnome.

### Phase 3: Polish & Refactoring
*   **Code Cleanup**:
    *   Identified duplication between `single.go` and `batch.go` regarding path generation and ratio calculation.
    *   Extracted this logic into reusable helper functions (`GenerateOutputPath`, `CalculateRatio`) in `internal/ui/common.go`.
*   **Application Icon**: Embedded a high-quality icon (`icon.png`) into the binary using Go's `//go:embed`.
*   **About Tab**: Added an info tab with developer details, GitHub links, and license statements.

### Phase 4: Cross-Platform & Documentation
*   **Build System**:
    *   Created `build.sh` with optimization flags (`-ldflags="-s -w"`).
    *   Added cross-compilation checks to gracefully handle missing Windows/Mac toolchains.
*   **Documentation**:
    *   Authored a comprehensive `README.md` covering features, dependencies (Ghostscript installation), and build steps.
    *   Corrected discrepancies (e.g., removing unimplemented "drag-and-drop" claims).
    *   Updated the License section to clarify the distinction between the Apache 2.0 app code and the AGPL Ghostscript runtime dependency.

### Phase 5: Final Polish
*   **Window Tuning**: Increased the initial application window width by 10% (to 660px) for better layout spacing.
*   **Bug Fixes**: Resolved syntax errors introduced during refactoring (`undefined` helpers, missing imports) to ensure a stable build.

## 3. Technical Deep Dive

### Libraries & Tools
| Library | Purpose | Reason for Choice |
| :--- | :--- | :--- |
| **fyne.io/fyne/v2** | GUI Toolkit | Go-native, lightweight, compiles to single binary. |
| **github.com/ncruces/zenity** | File Dialogs | Provides native OS file pickers, improving UX over Fyne's emulated dialogs. |
| **Ghostscript** | Compression Engine | Industry-standard PDF processing. Used via CLI execution to avoid complex CGO compilation and AGPL licensing issues. |

### Architecture
*   **Wrapper Pattern**: The app acts as a GUI wrapper around the `gs` command line tool.
*   **Worker Pool**: Batch processing uses a buffered channel of jobs and a fixed number of worker goroutines.
*   **Code Reusability**: Common logic (`GenerateOutputPath`, `CalculateRatio`) resides in `common.go` to prevent logic drift between Single and Batch modes.

## 4. Challenges & Solutions

| Challenge | Solution |
| :--- | :--- |
| **Tab Jumping** | The UI would switch tabs unexpectedly during compression. Fixed by implementing state callbacks (`onProcessStart`) to lock non-active tabs. |
| **Cross-Compilation** | Building for Windows from Linux requires CGO. Solution: Added checks for `mingw-w64` in `build.sh` to gracefully skip Windows builds if the compiler is missing. |
| **Refactoring Regressions** | Moving logic to `common.go` caused build failures (`undefined` functions, missing imports). Verified build logs, corrected `common.go` package structure, and fixed variable shadowing in `single.go`. |
| **Layout Spacing** | Default window size was slightly too narrow. Solution: Adjusted `resizeWindow` logic in `setup.go` to enforce a wider minimum width (660px). |

## 5. Artifacts Created
*   `README.md`: User guide and build instructions.
*   `Antigravity Interaction.md`: This log.
*   `build.sh`: Automated build script.
*   `.gitignore`: Build artifact exclusion.
*   `task.md` & `walkthrough.md`: Internal planning documents.
