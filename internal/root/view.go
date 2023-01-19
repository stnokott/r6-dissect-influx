package root

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/db"
)

const windowTitle string = "R6 Match InfluxDB Exporter"

type View struct {
	fyne.Window

	influxClient  *db.InfluxClient
	chConnUpdates <-chan db.ConnectionUpdate

	footer *footer
}

func NewView(a fyne.App) *View {
	w := a.NewWindow(windowTitle)
	w.Resize(fyne.NewSize(800, 600))

	v := &View{
		Window: w,
		footer: newFooter(),
	}

	v.resetClient()

	toolbar := widget.NewToolbar(
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(
			theme.SettingsIcon(),
			v.openSettings,
		),
	)

	w.SetContent(container.NewBorder(
		toolbar,
		v.footer,
		nil,
		nil,
		container.NewCenter(
			widget.NewLabel("<Placeholder>"),
		),
	))

	v.SetOnClosed(func() {
		if v.influxClient != nil {
			log.Println("closing connection")
			v.influxClient.Close()
		}
	})

	return v
}

func (v *View) resetClient() {
	if v.influxClient != nil {
		v.influxClient.Close()
	}
	v.footer.SetConnected(false)
	v.influxClient, v.chConnUpdates = db.NewInfluxClient(
		db.ConnectOpts{
			Host:            config.Current.InfluxDBHost,
			Port:            config.Current.InfluxDBPort,
			Token:           config.Current.InfluxDBToken,
			Org:             config.Current.InfluxDBOrg,
			Bucket:          config.Current.InfluxDBBucket,
			RefreshInterval: constants.INFLUX_PING_INTERVAL,
		},
	)
	// TODO: dont start if not configured
	v.influxClient.Start()
	go v.handlePingResult()
}

func (v *View) openSettings() {
	config.ShowDialog(v, v.onSettingsConfirmed)
}

func (v *View) onSettingsConfirmed() {
	// TODO: show infinite progress, wait for update
	v.resetClient()
}

func (v *View) handlePingResult() {
	for {
		update, ok := <-v.chConnUpdates
		if !ok {
			log.Println("stopping handlePingResult")
			// channel closed
			return
		}
		if update.Err == nil {
			v.footer.SetConnected(true)
		} else {
			v.footer.SetConnected(false)
		}
	}
}
