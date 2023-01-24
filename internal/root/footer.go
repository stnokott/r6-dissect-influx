package root

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/update"
	"github.com/stnokott/r6-dissect-influx/internal/utils"
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

type aboutDialog struct {
	dialog.Dialog

	parent fyne.Window

	btnCheckForUpdates   *widget.Button
	latestRelease        *update.Release
	lblNoUpdateAvailable *widget.Label
	lblUpdateAvailable   *widget.RichText
	btnUpdate            *widget.Button

	err error
}

func newAboutDialog(parent fyne.Window) *aboutDialog {
	d := &aboutDialog{
		lblNoUpdateAvailable: widget.NewLabel("No update available."),
		lblUpdateAvailable:   widget.NewRichText(),
		parent:               parent,
	}
	d.lblNoUpdateAvailable.Hide()
	d.lblUpdateAvailable.Hide()

	d.btnCheckForUpdates = widget.NewButtonWithIcon(
		"Check for updates",
		theme.ViewRefreshIcon(),
		func() {
			go d.checkForUpdates()
		},
	)
	d.btnUpdate = widget.NewButtonWithIcon(
		"Update",
		theme.DownloadIcon(),
		d.performUpdate,
	)
	d.btnUpdate.Hide()

	d.Dialog = dialog.NewCustom(
		"About",
		"Close",
		container.NewPadded(
			container.NewVBox(
				widget.NewRichTextFromMarkdown(
					"## "+constants.WINDOW_TITLE+"\n"+
						"**Version:** "+constants.Version,
				),
				container.NewHBox(
					widget.NewHyperlink(
						"GitHub Repository",
						constants.GithubURL,
					),
					layout.NewSpacer(),
				),
				container.NewHBox(
					d.lblUpdateAvailable,
					d.btnUpdate,
				),
				layout.NewSpacer(),
				d.btnCheckForUpdates,
				d.lblNoUpdateAvailable,
			),
		),
		parent,
	)
	d.SetOnClosed(func() {
		d.err = nil
		d.lblNoUpdateAvailable.Hide()
	})
	d.Resize(fyne.NewSize(d.MinSize().Width*1.4, d.MinSize().Height))
	return d
}

func (d *aboutDialog) checkForUpdates() {
	d.btnCheckForUpdates.Disable()
	var err error
	defer func() {
		// TODO: show error
		if err != nil {
			log.Println(err)
			d.err = err
			d.latestRelease = nil
		}
		d.btnCheckForUpdates.Enable()
	}()
	d.latestRelease, err = update.GetLatestRelease()
	if err != nil {
		return
	}
	d.updateContent()
}

func (d *aboutDialog) performUpdate() {
	progress := dialog.NewCustom(
		"Update",
		"Applying update...",
		widget.NewProgressBarInfinite(),
		d.parent,
	)
	progress.Show()

	go func() {
		err := d.latestRelease.DownloadAndApply()
		if err != nil {
			progress.Hide()
			utils.ShowErrDialog(err, d.Hide, d.parent)
		} else {
			progress.Hide()
			info := dialog.NewInformation(
				"Update applied",
				"Update successfully applied.\n"+
					"Please restart the application to see the changes.\n\n"+
					"The application will close now.",
				d.parent,
			)
			info.SetOnClosed(d.parent.Close)
			info.Show()
		}
	}()
}

func (d *aboutDialog) updateContent() {
	if d.latestRelease == nil || !d.latestRelease.IsNewer() {
		d.lblNoUpdateAvailable.Show()
		d.lblUpdateAvailable.Hide()
		d.btnUpdate.Hide()
	} else {
		d.lblNoUpdateAvailable.Hide()
		d.lblUpdateAvailable.ParseMarkdown(
			"**Update available to " + d.latestRelease.SemVer.String() + "**",
		)
		d.lblUpdateAvailable.Show()
		d.btnUpdate.Show()
	}
	d.Dialog.Refresh()
}
