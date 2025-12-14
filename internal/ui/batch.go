package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"simplepdfcompress/internal/compression"
	"simplepdfcompress/internal/worker"

	"github.com/ncruces/zenity"
)

func createBatchFileTab(w fyne.Window, onStart, onEnd func()) fyne.CanvasObject {
	var inputFiles []string
	var outputFolderURI fyne.URI

	// UI Elements
	fileListLabel := widget.NewLabel("No files added")
	fileListLabel.Wrapping = fyne.TextWrapWord

	outputLabel := widget.NewLabel("Default output: ./compressed (relative to each file)")
	outputLabel.Wrapping = fyne.TextWrapBreak

	qualitySelect := createQualitySelect()

	maxThreads := float64(runtime.NumCPU())
	threadSlider := widget.NewSlider(1, maxThreads)
	threadSlider.Value = maxThreads
	threadLabel := widget.NewLabel(fmt.Sprintf("Threads: %d", int(maxThreads)))
	threadSlider.OnChanged = func(f float64) {
		threadLabel.SetText(fmt.Sprintf("Threads: %d", int(f)))
	}

	suffixEntry := createSuffixEntry()

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	statusLabel := widget.NewLabel("")

	logEntry := widget.NewMultiLineEntry()
	logEntry.Disable() // Read-only log
	logEntry.SetMinRowsVisible(8)
	logScroll := container.NewVScroll(logEntry)
	logScroll.SetMinSize(fyne.NewSize(0, 150))

	// Buttons
	var compressBtn *widget.Button

	addFilesBtn := widget.NewButton("Add File", func() {
		go func() {
			// Native File Picker (Multi)
			filenames, err := zenity.SelectFileMultiple(
				zenity.Title("Select PDF Files"),
				zenity.FileFilter{Name: "PDF Files", Patterns: []string{"*.pdf"}},
			)
			if err == nil {
				fyne.Do(func() {
					inputFiles = append(inputFiles, filenames...)
					updateFileListLabel(fileListLabel, inputFiles)
				})
				return
			}

			if err != zenity.ErrCanceled {
				// Fallback
				fyne.Do(func() {
					fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
						if err != nil || reader == nil {
							return
						}
						inputFiles = append(inputFiles, reader.URI().Path())
						updateFileListLabel(fileListLabel, inputFiles)
					}, w)
					fd.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
					fd.Show()
				})
			}
		}()
	})

	addFolderBtn := widget.NewButton("Add Folder", func() {
		go func() {
			directory, err := zenity.SelectFile(
				zenity.Title("Select Folder"),
				zenity.Directory(),
			)
			if err == nil {
				// Scan folder for PDFs
				var pdfs []string
				err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".pdf") {
						pdfs = append(pdfs, path)
					}
					return nil
				})

				if err != nil {
					fyne.Do(func() {
						dialog.ShowError(fmt.Errorf("failed to scan folder: %w", err), w)
					})
					return
				}

				if len(pdfs) > 0 {
					fyne.Do(func() {
						inputFiles = append(inputFiles, pdfs...)
						updateFileListLabel(fileListLabel, inputFiles)
						dialog.ShowInformation("Folder Added", fmt.Sprintf("Added %d PDF files from folder.", len(pdfs)), w)
					})
				} else {
					fyne.Do(func() {
						dialog.ShowInformation("No PDFs Found", "No PDF files were found in the selected folder.", w)
					})
				}
				return
			}

			if err != zenity.ErrCanceled {
				// Fallback to Fyne Folder Dialog
				fyne.Do(func() {
					fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
						if err != nil || uri == nil {
							return
						}
						// Scan folder (Fyne URI)
						// Note: Walking Fyne URIs recursively is cleaner with Fyne storage API or converting to path if local.
						// Since we target local mostly, let's try path.
						dirPath := uri.Path()
						var pdfs []string
						filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
							if err == nil && !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".pdf") {
								pdfs = append(pdfs, path)
							}
							return nil
						})

						if len(pdfs) > 0 {
							inputFiles = append(inputFiles, pdfs...)
							updateFileListLabel(fileListLabel, inputFiles)
						}
					}, w)
					fd.Show()
				})
			}
		}()
	})

	clearFilesBtn := widget.NewButton("Clear List", func() {
		inputFiles = []string{}
		updateFileListLabel(fileListLabel, inputFiles)
		logEntry.SetText("")
	})

	selectOutputBtn := widget.NewButton("Select Output Folder (Optional)", func() {
		selectFolder(w, "Select Output Folder", func(uri fyne.URI) {
			outputFolderURI = uri
			outputLabel.SetText(uri.Path())
		})
	})

	compressBtn = widget.NewButton("Compress All", func() {
		if len(inputFiles) == 0 {
			dialog.ShowInformation("Info", "Please add files first", w)
			return
		}

		// Disable interactions
		compressBtn.Disable()
		addFilesBtn.Disable()
		addFolderBtn.Disable()
		clearFilesBtn.Disable()
		selectOutputBtn.Disable()
		qualitySelect.Disable()
		suffixEntry.Disable()
		// threadSlider.SetValue(threadSlider.Value) // Hack to keep visual state? No, SetValue doesn't disable.
		// There is no Disable() on slider in older Fyne versions easily exposed?
		// Actually widget.Slider has Disable().
		// Let's assume it does.
		// threadSlider.Disable() // Assuming available
		onStart()

		progressBar.Show()
		progressBar.SetValue(0)
		statusLabel.SetText(fmt.Sprintf("Starting compression of %d files...", len(inputFiles)))
		logEntry.SetText("Starting batch compression...\n")

		go func() {
			defer fyne.Do(func() {
				compressBtn.Enable()
				addFilesBtn.Enable()
				addFolderBtn.Enable()
				clearFilesBtn.Enable()
				selectOutputBtn.Enable()
				qualitySelect.Enable()
				suffixEntry.Enable()
				// threadSlider.Enable()
				onEnd()
			})

			// 1. Prepare Jobs & Check Overwrites
			jobs := make([]worker.Job, 0, len(inputFiles))
			opts := compression.CompressionOptions{Quality: qualitySelect.Selected}

			// Map to check overwrites
			var overwriteCandidates []string

			// Get suffix
			suffix := suffixEntry.Text
			if suffix == "" {
				suffix = "_spc_compressed"
			}

			for _, file := range inputFiles {
				outDirPath := ""
				if outputFolderURI != nil {
					outDirPath = outputFolderURI.Path()
				}
				outFile := GenerateOutputPath(file, outDirPath, suffixEntry.Text)

				if _, err := os.Stat(outFile); err == nil {
					overwriteCandidates = append(overwriteCandidates, outFile)
				}

				jobs = append(jobs, worker.Job{
					InputPath:  file,
					OutputPath: outFile,
					Options:    opts,
				})
			}

			// Ask permission if files exist
			if len(overwriteCandidates) > 0 {
				// Blocking dialog needs to be on Main? No, Zenity is fine in goroutine usually,
				// but let's check if we want Fyne dialog.
				// Zenity question:
				err := zenity.Question(
					fmt.Sprintf("Found %d existing files that will be overwritten. Continue?", len(overwriteCandidates)),
					zenity.Title("Overwrite Confirmation"),
					zenity.OKLabel("Overwrite"),
					zenity.CancelLabel("Cancel"),
				)
				if err != nil { // err is not nil if cancelled or error
					// Cancelled
					fyne.Do(func() {
						statusLabel.SetText("Cancelled.")
						logEntry.SetText(logEntry.Text + "\nCancelled by user.")
						progressBar.SetValue(0)
					})
					return
				}
			}

			// 2. Run Pool
			startTime := time.Now()
			numWorkers := int(threadSlider.Value)
			results := worker.RunPool(jobs, numWorkers)

			completed := 0
			total := len(jobs)
			var successes, failures int
			var unoptimizedFiles []string // Files that got bigger or didn't shrink well (negative ratio?)
			// Actually User asked "If the file failed to compress tell ... original already optimised... delete?"

			for res := range results {
				completed++
				progVal := float64(completed) / float64(total)

				var logMsg string
				if res.Error != nil {
					failures++
					logMsg = fmt.Sprintf("[X] %s: Failed - %v\n", filepath.Base(res.Job.InputPath), res.Error)
				} else {
					successes++
					// Calculate Ratio
					// (1 - Compressed/Original) * 100
					// If Compressed > Original, Ratio is negative.
					ratio := CalculateRatio(res.OriginalSize, res.FinalSize)

					logMsg = fmt.Sprintf("[O] %s: Ratio: %.1f%% (%s -> %s)\n",
						filepath.Base(res.Job.InputPath), ratio,
						formatBytes(res.OriginalSize), formatBytes(res.FinalSize))

					if res.FinalSize >= res.OriginalSize {
						unoptimizedFiles = append(unoptimizedFiles, res.Job.OutputPath)
						logMsg += "    -> Larger/Same size. Marked as unoptimized.\n"
					}
				}

				fyne.Do(func() {
					progressBar.SetValue(progVal)
					statusLabel.SetText(fmt.Sprintf("Processed %d/%d", completed, total))
					logEntry.SetText(logEntry.Text + logMsg)
					// Scroll to bottom? Fyne currently doesn't have easy auto-scroll entry method
					// But we can leave it.
				})
			}

			duration := time.Since(startTime)

			// 3. Post-Process Unoptimized
			if len(unoptimizedFiles) > 0 {
				err := zenity.Question(
					fmt.Sprintf("%d files were already optimized (compression did not reduce size). Delete these output files?", len(unoptimizedFiles)),
					zenity.Title("Delete Unoptimized Files?"),
					zenity.OKLabel("Delete"),
					zenity.CancelLabel("Keep"),
				)

				if err == nil {
					// Delete
					deletedCount := 0
					for _, f := range unoptimizedFiles {
						if remErr := os.Remove(f); remErr == nil {
							deletedCount++
						}
					}
					fyne.Do(func() {
						logEntry.SetText(logEntry.Text + fmt.Sprintf("\nDeleted %d unoptimized files.", deletedCount))
					})
				}
			}

			fyne.Do(func() {
				statusLabel.SetText(fmt.Sprintf("Done in %s. Success: %d, Failures: %d", duration.Round(time.Millisecond), successes, failures))
				progressBar.SetValue(1)
				dialog.ShowInformation("Batch Complete", fmt.Sprintf("Processed %d files in %s.\nSee log for details.", total, duration.Round(time.Millisecond)), w)

				// Clear file list to reset session
				inputFiles = []string{}
				updateFileListLabel(fileListLabel, inputFiles)
			})
		}()
	})
	compressBtn.Importance = widget.HighImportance

	// 33% width constraint
	compressBtnLayout := container.NewGridWithColumns(3, layout.NewSpacer(), compressBtn, layout.NewSpacer())

	// Main Content
	content := container.NewVBox(
		widget.NewLabelWithStyle("Batch File Compression", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("Files", container.NewVBox(fileListLabel, container.NewHBox(addFilesBtn, addFolderBtn, clearFilesBtn))),
			widget.NewFormItem("Output Folder", container.NewVBox(outputLabel, selectOutputBtn)),
			widget.NewFormItem("Quality", qualitySelect),
			widget.NewFormItem("Filename Suffix", suffixEntry),
			widget.NewFormItem("Max Threads", container.NewVBox(threadLabel, threadSlider)),
		),
		layoutSpacer(),
		widget.NewSeparator(),
		layoutSpacer(),
		progressBar,
		statusLabel,
		widget.NewLabel("Log:"),
		logScroll,
		layoutSpacer(),
		compressBtnLayout, // Modified
	)

	return container.NewPadded(content)
}

func updateFileListLabel(l *widget.Label, files []string) {
	if len(files) == 0 {
		l.SetText("No files selected")
		return
	}
	// Show first few lines
	msg := fmt.Sprintf("%d files selected:\n", len(files))
	for i, f := range files {
		if i >= 3 {
			msg += "... and more"
			break
		}
		msg += filepath.Base(f) + "\n"
	}
	l.SetText(msg)
}
