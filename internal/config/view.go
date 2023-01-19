package config

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/internal/utils"
)

const (
	WINDOW_TITLE       string = "Settings"
	hintMatchReplayDir string = "Directory where your match replays are stored. Usually a subfolder game folder named 'MatchReplays'"
	hintInfluxHost     string = `Host of your InfluxDB instance (IP or hostname, no "http(s)" required)`
	hintInfluxPort     string = "Port of your InfluxDB instance"
	hintInfluxOrg      string = "InfluxDB organization to connect to"
	hintInfluxBucket   string = "InfluxDB bucket where the data should be pushed"
	hintInfluxToken    string = "InfluxDB token for the bucket"
)

var (
	inputMatchReplayDir = &widget.Entry{}
	inputInfluxHost     = &widget.Entry{}
	inputInfluxPort     = &widget.Entry{}
	inputInfluxOrg      = &widget.Entry{}
	inputInfluxBucket   = &widget.Entry{}
	inputInfluxToken    = &widget.Entry{Password: true}
)

func ShowDialog(parent fyne.Window, onConfirm func()) {
	// 1. can't use NewEntryWithData on global level
	// 2. calling .Bind() removes validators, so we need to set them here (instead of when constructing the entry, see https://github.com/fyne-io/fyne/issues/2542)
	inputMatchReplayDir.Bind(bindMatchReplayDir)
	inputInfluxHost.Bind(bindInfluxHost)
	inputInfluxPort.Bind(bindInfluxPortStr)
	inputInfluxOrg.Bind(bindInfluxOrg)
	inputInfluxBucket.Bind(bindInfluxBucket)
	inputInfluxToken.Bind(bindInfluxToken)

	inputMatchReplayDir.Validator = validation.NewAllStrings(requiredValidator, gameDirectoryValidator)
	inputInfluxHost.Validator = validation.NewAllStrings(requiredValidator, hostAddressValidator)
	inputInfluxPort.Validator = validation.NewAllStrings(requiredValidator, integerValidator)
	inputInfluxOrg.Validator = requiredValidator
	inputInfluxBucket.Validator = requiredValidator
	inputInfluxToken.Validator = requiredValidator

	buttonAutodetectMatchReplayDir := widget.NewButton("Autodetect", func() {
		folder, err := matchReplayFolderFromRegistry()
		if err == nil {
			inputMatchReplayDir.SetText(folder)
		} else {
			utils.ShowErrDialog(err, nil, parent)
		}
	})
	dialogMatchReplayDir := dialog.NewFolderOpen(
		func(uc fyne.ListableURI, err error) {
			if uc != nil && err == nil {
				inputMatchReplayDir.SetText(uc.Path())
			}
		},
		parent,
	)
	buttonSelectMatchReplayDir := widget.NewButtonWithIcon("Open", theme.FolderIcon(), dialogMatchReplayDir.Show)

	formItems := []*widget.FormItem{
		{Text: "Game directory", Widget: inputMatchReplayDir, HintText: hintMatchReplayDir},
		{Text: "", Widget: container.NewGridWithColumns(2, buttonAutodetectMatchReplayDir, buttonSelectMatchReplayDir)},
		{Text: "InfluxDB host", Widget: inputInfluxHost, HintText: hintInfluxHost},
		{Text: "InfluxDB port", Widget: inputInfluxPort, HintText: hintInfluxPort},
		{Text: "InfluxDB org", Widget: inputInfluxOrg, HintText: hintInfluxOrg},
		{Text: "InfluxDB bucket", Widget: inputInfluxBucket, HintText: hintInfluxBucket},
		{Text: "InfluxDB token", Widget: inputInfluxToken, HintText: hintInfluxToken},
	}

	d := dialog.NewForm(
		WINDOW_TITLE,
		"Save",
		"Cancel",
		formItems,
		func(confirmed bool) {
			if confirmed {
				write()
				onConfirm()
			}
		},
		parent,
	)
	d.Show()
}
