package ui

import (
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
)

// Common UI Widgets

func createQualitySelect() *widget.Select {
	sel := widget.NewSelect([]string{"default", "screen", "ebook", "printer", "prepress"}, nil)
	sel.SetSelected("ebook")
	return sel
}

func createSuffixEntry() *widget.Entry {
	entry := widget.NewEntry()
	entry.SetText("_spc_compressed")
	entry.PlaceHolder = "_spc_compressed"
	return entry
}

func layoutSpacer() fyne.CanvasObject {
	return widget.NewLabel("")
}

// Dialog Helpers

func selectFolder(w fyne.Window, title string, onSelected func(fyne.URI)) {
	go func() {
		directory, err := zenity.SelectFile(
			zenity.Title(title),
			zenity.Directory(),
		)
		if err == nil {
			fyne.Do(func() {
				onSelected(storage.NewFileURI(directory))
			})
			return
		}

		if err != zenity.ErrCanceled {
			fyne.Do(func() {
				fd := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
					if err != nil || uri == nil {
						return
					}
					onSelected(uri)
				}, w)
				fd.Show()
			})
		}
	}()
}

func selectFile(w fyne.Window, title string, onSelected func(fyne.URI)) {
	go func() {
		filename, err := zenity.SelectFile(
			zenity.Title(title),
			zenity.FileFilter{Name: "PDF Files", Patterns: []string{"*.pdf"}},
		)
		if err == nil {
			fyne.Do(func() {
				onSelected(storage.NewFileURI(filename))
			})
			return
		}

		if err != zenity.ErrCanceled {
			fyne.Do(func() {
				fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
					if err != nil || reader == nil {
						return
					}
					onSelected(reader.URI())
				}, w)
				fd.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
				fd.Show()
			})
		}
	}()
}

// Logic Helpers

func GenerateOutputPath(inputFile, outputFolderURIPath, suffix string) string {
	outputDir := ""
	if outputFolderURIPath != "" {
		outputDir = outputFolderURIPath
	} else {
		outputDir = filepath.Join(filepath.Dir(inputFile), "compressed")
	}

	if suffix == "" {
		suffix = "_spc_compressed"
	}

	baseName := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	return filepath.Join(outputDir, fmt.Sprintf("%s%s.pdf", baseName, suffix))
}

func CalculateRatio(original, final int64) float64 {
	if original == 0 {
		return 0.0
	}
	return (1.0 - (float64(final) / float64(original))) * 100.0
}
