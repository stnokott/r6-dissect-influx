package root

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/update"
	"github.com/stnokott/r6-dissect-influx/internal/utils"
)

type aboutDialog struct {
	dialog.Dialog

	parent fyne.Window

	btnCheckForUpdates *widget.Button
	latestRelease      *update.Release

	lblNoUpdateAvailable *widget.Label

	containerUpdateAvailable *fyne.Container
	lblUpdateAvailable       *widget.RichText
	lblUpdateReleaseNotes    *widget.RichText
	btnOpenRelease           *widget.Button
	btnUpdate                *widget.Button

	err error
}

func newAboutDialog(parent fyne.Window) *aboutDialog {
	d := &aboutDialog{
		lblNoUpdateAvailable:  widget.NewLabel("No update available."),
		lblUpdateAvailable:    widget.NewRichText(),
		lblUpdateReleaseNotes: widget.NewRichText(),
		parent:                parent,
	}
	d.lblNoUpdateAvailable.Hide()

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
	d.btnOpenRelease = widget.NewButtonWithIcon(
		"Show on GitHub",
		theme.InfoIcon(),
		func() {
			utils.OpenURL(d.latestRelease.URL)
		},
	)

	d.containerUpdateAvailable = container.NewVBox(
		container.NewHBox(
			d.lblUpdateAvailable,
			d.btnOpenRelease,
			layout.NewSpacer(),
			d.btnUpdate,
		),
		widget.NewAccordion(widget.NewAccordionItem(
			"Release notes",
			d.lblUpdateReleaseNotes,
		)),
	)
	d.containerUpdateAvailable.Hide()

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
				d.containerUpdateAvailable,
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
		"Downloading update...",
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
				"Download complete",
				"Update successfully downloaded and applied.\n\n"+
					"Please restart the application now.",
				d.parent,
			)
			info.SetDismissText("Restart")
			info.SetOnClosed(d.onUpdateComplete)
			info.Show()
		}
	}()
}

func (d *aboutDialog) onUpdateComplete() {
	if err := utils.RestartApp(); err != nil {
		d.Hide()
		e := dialog.NewError(
			fmt.Errorf( //lint:ignore ST1005 will be displayed in UI
				"Failed to restart application after applying update:\n"+
					"%w\nPlease restart the application yourself.\n\n"+
					"The app will close now.",
				err,
			),
			d.parent,
		)
		e.SetDismissText("Exit")
		e.SetOnClosed(d.parent.Close)
	}
}

func (d *aboutDialog) updateContent() {
	if d.latestRelease == nil || !d.latestRelease.IsNewer() {
		d.lblNoUpdateAvailable.Show()
		d.containerUpdateAvailable.Hide()
	} else {
		d.lblNoUpdateAvailable.Hide()
		d.lblUpdateAvailable.ParseMarkdown(
			"New version available: **" + d.latestRelease.SemVer.String() + "**",
		)
		d.lblUpdateReleaseNotes.ParseMarkdown(d.latestRelease.Body)
		d.containerUpdateAvailable.Show()
	}
	d.Dialog.Refresh()
}
