package root

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

type footer struct {
	widget.BaseWidget
	parent fyne.Window

	aboutDialog *aboutDialog
}

func newFooter(parent fyne.Window) *footer {
	f := &footer{
		parent:      parent,
		aboutDialog: newAboutDialog(parent),
	}
	f.ExtendBaseWidget(f)

	return f
}

func (f *footer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewHBox(
		layout.NewSpacer(),
		&canvas.Text{
			Text:     fmt.Sprintf("v%s - %s", constants.Version, constants.Commit),
			TextSize: theme.CaptionTextSize(),
			Color:    theme.DisabledColor(),
		},
		&widget.Button{
			Icon:       theme.InfoIcon(),
			Importance: widget.LowImportance,
			OnTapped:   f.aboutDialog.Show,
		},
	))
}
