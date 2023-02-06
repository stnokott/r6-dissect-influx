package root

import (
	"log"
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
	"github.com/stnokott/r6-dissect-influx/internal/game"
	"github.com/stnokott/r6-dissect-influx/internal/ui/footer"
	"github.com/stnokott/r6-dissect-influx/internal/ui/matches"
	"github.com/stnokott/r6-dissect-influx/internal/utils"
)

type Window struct {
	fyne.Window

	influxClient *db.InfluxClient

	borderContainer *fyne.Container
	centerObject    fyne.CanvasObject
}

func NewWindow(a fyne.App) *Window {
	w := a.NewWindow(constants.WINDOW_TITLE)
	w.Resize(fyne.NewSize(800, 600))

	v := &Window{
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
		footer.New(v),
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

func (v *Window) replaceCenter(newCenter fyne.CanvasObject) {
	if v.centerObject != nil {
		v.borderContainer.Remove(v.centerObject)
	}
	v.borderContainer.Add(newCenter)
	v.centerObject = newCenter
	v.borderContainer.Refresh()
}

func (v *Window) validateConfig() {
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

func (v *Window) blockUntilConfigured() {
	v.replaceCenter(container.NewCenter(
		container.NewVBox(
			widget.NewLabel("Configuration required."),
			widget.NewButton("Setup", v.openSettings),
		),
	))
}

func (v *Window) loadMainView() {
	matchList := matches.NewMatchListView()

	v.replaceCenter(container.NewMax(matchList))

	reader, err := game.NewRoundsReader(config.Current.GameFolder)
	if err != nil {
		panic(err)
	}
	chRoundInfos, chErrors := reader.WatchAsync()
	go func() {
		for {
			select {
			case roundInfo, ok := <-chRoundInfos:
				if !ok {
					log.Println("match data channel closed")
					return
				}
				log.Println("got match info for ID:", roundInfo.MatchID)
				matchList.Add(&roundInfo)
			case err, ok := <-chErrors:
				if !ok {
					log.Println("match errors channel closed")
					return
				}
				if err != nil {
					log.Println("got error from match data channel:", err)
				}
			}
		}
	}()
}

func (v *Window) updateInfluxClient(c *db.InfluxClient) {
	if v.influxClient != nil {
		v.influxClient.Close()
	}
	v.influxClient = c
}

func (v *Window) openSettings() {
	//config.ShowDialog(v, v.onSettingsConfirmed)
}

func (v *Window) onSettingsConfirmed() {
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

func (v *Window) onClosed() {
	if v.influxClient != nil {
		v.influxClient.Close()
	}
}
