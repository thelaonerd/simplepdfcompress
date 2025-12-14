# SimplePDFCompress

SimplePDFCompress is a fast, efficient, and user-friendly desktop application for compressing PDF files. Built with Go and Fyne, it provides a native experience across Linux, Windows, and macOS.

## Features

*   **Single File Compression**: Quickly reduce the size of individual PDF documents.
*   **Batch Compression**: Process folders or multiple files at once.
*   **Adjustable Quality**: Choose from multiple compression presets (Screen, Ebook, Printer, Prepress) to balance quality and file size.
*   **System Integration**:
    *   Uses your system's default font (Linux) for a native look.
    *   Native file dialogs for a familiar experience.
*   **Multi-threaded**: Utilizes your CPU cores responsibly for fast batch processing.
*   **Cross-Platform**: Runs on Linux [verified], Windows [verified in RDP with opengl32.dll in the same folder], and macOS [unverified].

## How to Use

### Single File Mode
1.  Open the **Single File Compression** tab.
2.  Select a PDF file using the "Select PDF File" button.
3.  (Optional) Choose an output folder. By default, a `compressed` folder is created next to your file.
4.  Select your desired **Quality** (default is "Ebook").
5.  Click **Compress**.

### Batch Mode
1.  Open the **Batch File Compression** tab.
2.  Add files individually or add entire folders containing PDFs.
3.  Adjust the **Max Threads** slider to control performance.
4.  Click **Compress All**.

---

## Runtime Dependencies

This application uses **Ghostscript** as its compression engine. You must have it installed for the application to work.

### 🐧 Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install ghostscript
```

### 🐧 Linux (Fedora)
```bash
sudo dnf install ghostscript
```

### 🍎 macOS
We recommend using [Homebrew](https://brew.sh/):
```bash
brew install ghostscript
```

### 🪟 Windows
1.  Download the **Ghostscript AGPL Release** (files named `gs...exe`) from the [official website](https://ghostscript.com/releases/gsdnld.html).
2.  Install the version matching your system (usually 64-bit).
3.  **Important**: Ensure the installer adds Ghostscript to your system `PATH` environment variable so the application can find it.
4.  For windows RDP sessions, running required opengl32.dll and can be downloaded from https://downloads.fdossena.com/geth.php?r=mesa64-latest. This is needed to execute the application in RDP sessions.  

---

## How to Compile / Build

To build SimplePDFCompress from source, you need Go and a C compiler (for the Fyne GUI toolkit).

### Prerequisites (Build Tools)

#### 🐧 Linux (Ubuntu/Debian)
```bash
sudo apt install golang gcc libgl1-mesa-dev xorg-dev
```

#### 🐧 Linux (Fedora)
```bash
sudo dnf install golang gcc libX11-devel libXcursor-devel libXrandr-devel libXinerama-devel libXi-devel libXext-devel libXfixes-devel
```

#### 🍎 macOS
1.  Install Xcode Command Line Tools: `xcode-select --install`
2.  Install Go: `brew install go`

#### 🪟 Windows
1.  Install **Go** from [go.dev](https://go.dev/dl/).
2.  Install a C compiler like **TDM-GCC** or **MinGW-w64**.

### Building the Application

1.  Clone the repository:
    ```bash
    git clone https://github.com/thelaonerd/simplepdfcompress.git
    cd simplepdfcompress
    ```

2.  Run the build script (Linux/macOS):
    ```bash
    ./build.sh
    ```
    *This script produces optimized binaries in the `bin/` directory.*

    **Or manual build:**
    ```bash
    go build -ldflags="-s -w" .
    ```

---

## License

**SimplePDFCompress** is licensed under the **Apache License 2.0**.

> **Note on Ghostscript**: This application relies on Ghostscript, which is licensed under the **GNU AGPL**. SimplePDFCompress calls Ghostscript as an external command-line process and does not link against it statically or dynamically. Users must comply with Ghostscript's licensing terms when using it.

---

⭐ **If you found this application useful, please give it a star on GitHub!**
