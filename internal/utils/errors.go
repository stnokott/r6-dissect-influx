package utils

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

func TitleErr(err error, trailingDot bool) error {
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
	return errors.New(s)
}

func ShowErrDialog(err error, parent fyne.Window) {
	dialog.NewError(TitleErr(err, true), parent).Show()
}
