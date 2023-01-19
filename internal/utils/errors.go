package utils

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

func TitleErr(err error, trailingDot bool) string {
	s := err.Error()
	parts := strings.SplitN(s, " ", 2)
	firstTitle := titleCaser.String(parts[0])
	if len(parts) == 1 {
		s = firstTitle
	} else {
		s = firstTitle + " " + parts[1]
	}
	if trailingDot && s[len(s)-1] != '.' {
		s += "."
	}
	return s
}

func ShowErrDialog(err error, onClose func(), parent fyne.Window) {
	text := TitleErr(err, true)
	label := &widget.Label{
		Text:     text,
		Wrapping: fyne.TextWrapWord,
	}
	d := dialog.NewCustom(
		"Error",
		"Dismiss",
		label,
		parent,
	)
	textSize := fyne.MeasureText(text, theme.TextSize(), label.TextStyle)
	d.Resize(textSize.Min(parent.Canvas().Size()))
	if onClose != nil {
		d.SetOnClosed(onClose)
	}
	d.Show()
}
