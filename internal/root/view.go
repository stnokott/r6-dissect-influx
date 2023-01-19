package root

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/db"
	"github.com/stnokott/r6-dissect-influx/internal/utils"
)

type View struct {
	fyne.Window

	influxClient *db.InfluxClient

	borderContainer *fyne.Container
	centerObject    fyne.CanvasObject
}

func NewView(a fyne.App) *View {
	w := a.NewWindow(constants.WINDOW_TITLE)
	w.Resize(fyne.NewSize(800, 600))

	v := &View{
		Window: w,
	}
	v.SetOnClosed(v.onClosed)

	v.borderContainer = container.NewBorder(
		widget.NewToolbar(
			widget.NewToolbarSpacer(),
			widget.NewToolbarAction(
				theme.SettingsIcon(),
				v.openSettings,
			),
		),
		newFooter(),
		nil,
		nil,
		layout.NewSpacer(), // placeholder
	)
	v.SetContent(v.borderContainer)

	if config.IsComplete() {
		go v.validateConfig()
	} else {
		v.blockUntilConfigured()
	}

	return v
}

func (v *View) replaceCenter(newCenter fyne.CanvasObject) {
	if v.centerObject != nil {
		v.borderContainer.Remove(v.centerObject)
	}
	v.borderContainer.Add(newCenter)
	v.centerObject = newCenter
	v.borderContainer.Refresh()
}

func (v *View) validateConfig() {
	v.replaceCenter(container.NewCenter(
		container.NewVBox(
			widget.NewProgressBarInfinite(),
			widget.NewLabel("Validating config..."),
		),
	))
	client := config.Current.NewInfluxClient()
	details := client.ValidateConn(10 * time.Second)
	if details.Err != nil {
		utils.ShowErrDialog(details.Err, v.openSettings, v.Window)
		v.blockUntilConfigured()
	} else {
		v.loadMainView()
	}
}

func (v *View) blockUntilConfigured() {
	v.replaceCenter(container.NewCenter(
		container.NewVBox(
			widget.NewLabel("Configuration required."),
			widget.NewButton("Setup", v.openSettings),
		),
	))
}

func (v *View) loadMainView() {
	v.replaceCenter(container.NewCenter(
		widget.NewLabel("PLACEHOLDER: connection validated"),
	))
}

func (v *View) updateInfluxClient(c *db.InfluxClient) {
	if v.influxClient != nil {
		v.influxClient.Close()
	}
	v.influxClient = c
}

func (v *View) openSettings() {
	config.ShowDialog(v, v.onSettingsConfirmed)
}

func (v *View) onSettingsConfirmed() {
	newClient := config.Current.NewInfluxClient()
	progressDialog := dialog.NewCustom(
		"Validating connection...",
		"Attempting to connect to InfluxDB...",
		widget.NewProgressBarInfinite(),
		v,
	)
	progressDialog.Show()
	details := newClient.ValidateConn(10 * time.Second)
	if details.Err != nil {
		progressDialog.Hide()
		utils.ShowErrDialog(details.Err, v.openSettings, v)
	} else {
		v.updateInfluxClient(newClient)
		progressDialog.Hide()
		dialog.NewCustom(
			"Success",
			"Yay",
			container.NewPadded(widget.NewRichTextFromMarkdown(
				"### Connection successful\n"+
					"URL: **"+newClient.URL+"**\n\n"+
					"Name: **"+details.Name+"**\n\n"+
					"Version: **"+details.Version+"**\n\n"+
					"Commit: **"+details.Commit+"**",
			)),
			v,
		).Show()
		v.loadMainView()
	}
}

func (v *View) onClosed() {
	if v.influxClient != nil {
		v.influxClient.Close()
	}
}
