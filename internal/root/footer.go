package root

import (
	"fmt"
	"time"

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

	btnOpenAbout       *widget.Button
	lblUpdateAvailable *widget.RichText
	aboutDialog        *aboutDialog
}

func newFooter(parent fyne.Window) *footer {
	f := &footer{
		parent:             parent,
		lblUpdateAvailable: widget.NewRichTextFromMarkdown("**Update available**"),
		aboutDialog:        newAboutDialog(parent),
	}
	f.lblUpdateAvailable.Hide()
	f.btnOpenAbout = &widget.Button{
		Icon:       theme.InfoIcon(),
		Importance: widget.LowImportance,
		OnTapped:   f.aboutDialog.Show,
	}
	f.ExtendBaseWidget(f)

	return f
}

func (f *footer) CreateRenderer() fyne.WidgetRenderer {
	go f.updateChecker()

	return widget.NewSimpleRenderer(container.NewHBox(
		layout.NewSpacer(),
		&canvas.Text{
			Text:     fmt.Sprintf("v%s - %s", constants.Version, constants.Commit),
			TextSize: theme.CaptionTextSize(),
			Color:    theme.DisabledColor(),
		},
		f.lblUpdateAvailable,
		f.btnOpenAbout,
	))
}

func (f *footer) updateChecker() {
	ticker := time.NewTicker(constants.UpdateCheckInterval)
	for {
		updateAvailable := f.aboutDialog.CheckForUpdates()
		if updateAvailable {
			f.lblUpdateAvailable.Show()
			f.btnOpenAbout.SetIcon(theme.DownloadIcon())
		} else {
			f.lblUpdateAvailable.Hide()
			f.btnOpenAbout.SetIcon(theme.InfoIcon())
		}
		<-ticker.C
	}
}
