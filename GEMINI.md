# SimplePDFCompress Design Specification

## 1. Application Launch & Environment Check
**Objective:** Identify the runtime environment availability of dependencies.

### Operating System Identification
*   On application launch, the system must identify the underlying OS.
*   **Linux Specifics:** If the OS is Linux, determine the specific distribution (e.g., Ubuntu, Fedora, Arch) and the associated package manager (e.g., `apt`, `dnf`, `pacman`).

### Dependency Verification
*   Check for the presence of the `gs` (Ghostscript) and `ps2pdf` executable packages.
    *   **Dependencies Satisfied:** Display a confirmation message establishing that the application is ready to run.
    *   **Dependencies Missing:** Display precise installation steps tailored to the identified platform/distribution (e.g., "Run `sudo apt install ghostscript`").

## 2. Compression Modes
The application must expose all `ps2pdf` options and operate in one of two modes:

### A. Single File Compression
*   Allows the user to select and compress a single PDF document.

### B. Batch File Compression
*   Allows the user to select multiple files or directories for mass compression.

## 3. Output Configuration
*   **Output Folder Selection:** Users can choose a specific destination directory for the compressed files.
*   **Default Behavior:** If no output folder is specified, the application must create a new folder named `compressed` within the current working directory and save files there.

## 4. Execution & Processing
*   **Method:** File compression is performed by recursively invoking the `ps2pdf` command.
*   **Execution Model:** The command execution is blocking (synchronous) per file context to ensure data integrity.
*   **Post-Processing Visualization:** Upon completion, the application must calculate and display the **Compression Ratio** to the user (e.g., "Reduced file size by 45%").

## 5. Performance & Concurrency
*   **Multiprocessing/Multi-threading:** The application must utilize concurrency for batch operations.
*   **User Control:** Provide a UI element allowing users to select the number of worker threads, bounded by the number of available CPU cores/threads.

## 6. Implementation Tasks
- [ ] **Project Initialization**: Initialize Go module and install Fyne.
- [ ] **System Checks**: Implement OS detection and correct `gs`/`ps2pdf` dependency checks.
- [ ] **Core Logic**: Create `ps2pdf` wrapper with blocking execution and error handling.
- [ ] **Concurrency**: Implement worker pool for batch processing.
- [ ] **UI - Main**: Build main layout with File/Different modes.
- [ ] **UI - Feedback**: Add progress bars and compression ratio display.
- [ ] **Output Logic**: Handle default folder creation.

## 7. References
*   **Fyne.io Documentation**: [https://docs.fyne.io/](https://docs.fyne.io/)

