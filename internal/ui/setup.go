package ui

import (
	"fmt"
	"simplepdfcompress/internal/system"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Setup initializes the application UI based on system checks
func Setup(w fyne.Window, a fyne.App) {
	checks := system.PerformChecks()
	if !checks.IsReady {
		w.SetContent(createDependencyErrorScreen(checks))
		return
	}
	w.SetContent(createMainScreen(w))
}

func createDependencyErrorScreen(checks system.CheckResult) fyne.CanvasObject {
	icon := widget.NewIcon(nil) // Placeholder for warning icon

	title := widget.NewLabelWithStyle("Missing Dependencies", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	msg := widget.NewLabel(checks.Message)
	msg.Wrapping = fyne.TextWrapWord

	return container.NewVBox(
		icon,
		title,
		msg,
		widget.NewButton("Quit", func() {
			fyne.CurrentApp().Quit()
		}),
	)
}

func createMainScreen(w fyne.Window) fyne.CanvasObject {
	var tabs *container.AppTabs

	// State management callbacks
	// onProcessStart disables tab switching (except current)
	onProcessStart := func() {
		if tabs != nil {
			selected := tabs.Selected()
			for _, item := range tabs.Items {
				if item != selected {
					tabs.DisableItem(item)
				}
			}
		}
	}

	// onProcessEnd re-enables tab switching
	onProcessEnd := func() {
		if tabs != nil {
			for _, item := range tabs.Items {
				tabs.EnableItem(item)
			}
		}
	}

	single := createSingleFileTab(w, onProcessStart, onProcessEnd)
	batch := createBatchFileTab(w, onProcessStart, onProcessEnd)

	tabs = container.NewAppTabs(
		container.NewTabItem("Single File", single),
		container.NewTabItem("Batch Compression", batch),
		container.NewTabItem("About", createAboutTab()),
	)

	// Dynamic Resizing Logic
	resizeWindow := func(content fyne.CanvasObject) {
		padding := fyne.NewSize(40, 40)
		targetSize := content.MinSize().Add(padding)
		if targetSize.Width < 660 {
			targetSize.Width = 660
		}
		w.Resize(targetSize)
	}

	tabs.OnSelected = func(t *container.TabItem) {
		resizeWindow(t.Content)
	}

	// Set initial size
	resizeWindow(single)

	return tabs
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
