package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"simplepdfcompress/internal/compression"

	"github.com/ncruces/zenity"
)

func createSingleFileTab(w fyne.Window, onStart, onEnd func()) fyne.CanvasObject {
	var selectedFileURI fyne.URI
	var outputFolderURI fyne.URI

	// UI Elements
	fileLabel := widget.NewLabel("No file selected")
	fileLabel.Truncation = fyne.TextTruncateEllipsis

	outputLabel := widget.NewLabel("Default output: ./compressed")
	outputLabel.Truncation = fyne.TextTruncateEllipsis

	qualitySelect := createQualitySelect()

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	statusLabel := widget.NewLabel("")

	logEntry := widget.NewMultiLineEntry()
	logEntry.Disable() // Read-only log
	logEntry.SetMinRowsVisible(8)
	logScroll := container.NewVScroll(logEntry)
	logScroll.SetMinSize(fyne.NewSize(0, 150))

	// Buttons
	selectFileBtn := widget.NewButton("Select PDF File", func() {
		selectFile(w, "Select PDF File", func(uri fyne.URI) {
			selectedFileURI = uri
			fileLabel.SetText(uri.Path())
			logEntry.SetText("") // Clear log on new file
		})
	})

	selectOutputBtn := widget.NewButton("Select Output Folder (Optional)", func() {
		selectFolder(w, "Select Output Folder", func(uri fyne.URI) {
			outputFolderURI = uri
			outputLabel.SetText(uri.Path())
		})
	})

	suffixEntry := createSuffixEntry()

	var compressBtn *widget.Button
	compressBtn = widget.NewButton("Compress", func() {
		if selectedFileURI == nil {
			dialog.ShowError(errors.New("please select a file first"), w)
			return
		}

		// Disable interactions
		compressBtn.Disable()
		selectFileBtn.Disable() // Good practice to disable inputs too
		selectOutputBtn.Disable()
		qualitySelect.Disable()
		suffixEntry.Disable()
		onStart()

		progressBar.Show()
		progressBar.SetValue(0)
		statusLabel.SetText("Compressing...")
		logEntry.SetText("Starting compression...\n")

		// Run in background
		go func() {
			defer fyne.Do(func() {
				compressBtn.Enable()
				selectFileBtn.Enable()
				selectOutputBtn.Enable()
				qualitySelect.Enable()
				suffixEntry.Enable()
				onEnd()
			})

			// 1. Determine Output Path & Check Overwrite
			// 1. Determine Output Path
			inputFile := selectedFileURI.Path()
			outDirPath := ""
			if outputFolderURI != nil {
				outDirPath = outputFolderURI.Path()
			}
			outputFile := GenerateOutputPath(inputFile, outDirPath, suffixEntry.Text)

			logEntryAppend := func(s string) {
				fyne.Do(func() {
					logEntry.SetText(logEntry.Text + s)
				})
			}

			// Overwrite Check
			if _, err := os.Stat(outputFile); err == nil {
				// File exists
				err := zenity.Question(
					fmt.Sprintf("File already exists:\n%s\nOutput will be overwritten. Continue?", filepath.Base(outputFile)),
					zenity.Title("Overwrite Confirmation"),
					zenity.OKLabel("Overwrite"),
					zenity.CancelLabel("Cancel"),
				)
				if err != nil {
					// Cancelled
					fyne.Do(func() {
						statusLabel.SetText("Cancelled by user.")
						progressBar.SetValue(0)
						logEntry.SetText(logEntry.Text + "Cancelled: File exists and user chose not to overwrite.\n")
					})
					return
				}
			}

			// 2. Compress
			startTime := time.Now()
			opts := compression.CompressionOptions{
				Quality: qualitySelect.Selected,
			}

			initial, final, err := compression.CompressPDF(inputFile, outputFile, opts)

			if err != nil {
				fyne.Do(func() {
					dialog.ShowError(err, w)
					progressBar.SetValue(0)
					statusLabel.SetText("Error: " + err.Error())
					logEntry.SetText(logEntry.Text + fmt.Sprintf("Error: %s\n", err.Error()))
				})
				return
			}

			duration := time.Since(startTime)

			// 3. Ratio & Unoptimized Logic
			// 3. Ratio & Unoptimized Logic
			ratio := CalculateRatio(initial, final)

			// Msg for success
			msg := fmt.Sprintf("Success! (Time: %s)\nRatio: %.1f%%\n\nOriginal: %s\nCompressed: %s",
				duration.Round(time.Millisecond), ratio, formatBytes(initial), formatBytes(final))

			logEntryAppend(fmt.Sprintf("Success! Ratio: %.1f%% (%s -> %s) in %s\n",
				ratio, formatBytes(initial), formatBytes(final), duration.Round(time.Millisecond)))

			// Unoptimized check
			if final >= initial {
				msg += "\n\nWarning: File did not shrink (already optimized)."
				logEntryAppend("Warning: File did not shrink.\n")

				err := zenity.Question(
					"Compression did not reduce file size (it is larger or same size).\nDelete the new output file?",
					zenity.Title("Unoptimized Result"),
					zenity.OKLabel("Delete"),
					zenity.CancelLabel("Keep"),
				)

				if err == nil {
					// Delete
					if remErr := os.Remove(outputFile); remErr == nil {
						fyne.Do(func() {
							dialog.ShowInformation("Deleted", "Unoptimized output file was deleted.", w)
							statusLabel.SetText("Cancelled (Unoptimized deleted).")
							progressBar.SetValue(0)
							logEntry.SetText(logEntry.Text + "Deleted unoptimized output file.\n")
						})
						return
					} else {
						msg += "\n(Failed to delete unoptimized file)"
						logEntryAppend("Failed to delete unoptimized file.\n")
					}
				}
			}

			fyne.Do(func() {
				dialog.ShowInformation("Compression Complete", msg, w)
				progressBar.SetValue(1)
				// avoiding layout change (Hide)
				statusLabel.SetText(fmt.Sprintf("Saved to %s (%s)", filepath.Base(outputFile), duration.Round(time.Millisecond)))
				logEntry.SetText(logEntry.Text + fmt.Sprintf("Saved to: %s\n", outputFile))
			})
		}()
	})
	compressBtn.Importance = widget.HighImportance

	// Compress button layout: 33% width (middle column of 3)
	compressBtnLayout := container.NewGridWithColumns(3, layout.NewSpacer(), compressBtn, layout.NewSpacer())

	// Main Content
	content := container.NewVBox(
		widget.NewLabelWithStyle("Single File Compression", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("Input File", container.NewVBox(fileLabel, selectFileBtn)),
			widget.NewFormItem("Output Folder", container.NewVBox(outputLabel, selectOutputBtn)),
			widget.NewFormItem("Quality", qualitySelect),
			widget.NewFormItem("Filename Suffix", suffixEntry),
		),
		layoutSpacer(),
		widget.NewSeparator(),
		layoutSpacer(),
		progressBar,
		statusLabel,
		widget.NewLabel("Log:"),
		logScroll,
		layoutSpacer(),
		compressBtnLayout, // Modified layout
	)

	return container.NewPadded(content)
}
