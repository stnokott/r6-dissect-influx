package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

func TitleErr(err error, trailingDot bool) (out string) {
	defer func() {
		if trailingDot {
			// add trailing dot
			if out[len(out)-1] != '.' {
				out += "."
			}
		}
	}()

	s := err.Error()
	parts := strings.SplitN(s, " ", 2)
	firstTitle := titleCaser.String(parts[0])
	if len(parts) == 1 {
		out = firstTitle
	} else {
		out = firstTitle + " " + parts[1]
	}
	return
}
