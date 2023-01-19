package root

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

func newFooter() fyne.CanvasObject {
	return container.NewHBox(
		layout.NewSpacer(),
		&canvas.Text{
			Text:     fmt.Sprintf("v%s - %s", constants.Version, constants.Commit),
			TextSize: theme.CaptionTextSize(),
			Color:    theme.DisabledColor(),
		},
	)
}
