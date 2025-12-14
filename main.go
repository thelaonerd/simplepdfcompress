package main

import (
	_ "embed"
	"fmt"
	"os/exec"
	"runtime"
	"simplepdfcompress/internal/ui"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

//go:embed icon.png
var iconData []byte

func main() {
	a := app.NewWithID("com.simplepdfcompress.app")

	// Attempt to set system font (Linux/fontconfig)
	if runtime.GOOS == "linux" {
		if out, err := exec.Command("fc-match", "-f", "%{file}", "sans-serif").Output(); err == nil {
			fontPath := strings.TrimSpace(string(out))
			if fontPath != "" {
				fmt.Println("Detected system font:", fontPath)
				if sysTheme := ui.NewSystemTheme(fontPath); sysTheme != nil {
					a.Settings().SetTheme(sysTheme)
					fmt.Println("Applied system font theme.")
				}
			}
		}
	}

	w := a.NewWindow("SimplePDFCompress")
	w.Resize(fyne.NewSize(600, 400))

	// Set Icon (App & Window)
	resource := fyne.NewStaticResource("icon.png", iconData)
	a.SetIcon(resource)
	w.SetIcon(resource)
	fmt.Println("Icon loaded from embedded data.")

	ui.Setup(w, a)

	w.ShowAndRun()
}
