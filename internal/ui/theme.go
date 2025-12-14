package ui

import (
	"image/color"
	// Will create this or load dynamically
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type SystemTheme struct {
	fontResource fyne.Resource
}

func NewSystemTheme(fontPath string) *SystemTheme {
	res, err := fyne.LoadResourceFromPath(fontPath)
	if err != nil {
		return nil
	}
	return &SystemTheme{fontResource: res}
}

func (t *SystemTheme) Font(s fyne.TextStyle) fyne.Resource {
	// For simplicity, return same font for all styles, or fallback if needed.
	// Valid system font usually covers regular. Bold/Italic might need separate loading if perfect match desired,
	// but using Regular for all is a start if we only have one path.
	if t.fontResource != nil {
		return t.fontResource
	}
	return theme.DefaultTheme().Font(s)
}

func (t *SystemTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (t *SystemTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t *SystemTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
