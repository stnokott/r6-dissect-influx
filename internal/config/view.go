package config

import (
	"errors"
	"os"
	"regexp"
	"strconv"

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
	hintMatchReplayDir string = "Directory where your match replays are stored. Usually in your game folder named 'MatchReplays'"
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

func ShowDialog(parent fyne.Window) {
	// 1. can't use NewEntryWithData on global level
	// 2. calling .Bind() removes validators, so we need to set them here (instead of when constructing the entry, see https://github.com/fyne-io/fyne/issues/2542)
	inputMatchReplayDir.Bind(bindMatchReplayDir)
	inputInfluxHost.Bind(bindInfluxHost)
	inputInfluxPort.Bind(bindInfluxPortStr)
	inputInfluxOrg.Bind(bindInfluxOrg)
	inputInfluxBucket.Bind(bindInfluxBucket)
	inputInfluxToken.Bind(bindInfluxToken)

	inputMatchReplayDir.Validator = validation.NewAllStrings(requiredValidator, directoryValidator)
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
			utils.ShowErrDialog(err, parent)
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
		{Text: "Match replay directory", Widget: inputMatchReplayDir, HintText: hintMatchReplayDir},
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
				Write()
			}
		},
		parent,
	)
	d.Resize(d.MinSize().Add(fyne.NewDelta(100, 0)))
	d.Show()
}

// TODO: warn if no matches found or none could be parsed
func directoryValidator(s string) (err error) {
	stats, statErr := os.Stat(s)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			err = errors.New("Does not exist") //lint:ignore ST1005 will be displayed in UI
		} else {
			err = statErr
		}
	} else if !stats.IsDir() {
		err = errors.New("Not a directory") //lint:ignore ST1005 will be displayed in UI
	}
	return
}

var (
	regexIPv4     = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	regexHostname = regexp.MustCompile(`^(?:[0-9a-zA-Z]+\.)+[0-9a-zA-Z]{2,4}$`)
)

func hostAddressValidator(s string) (err error) {
	if regexIPv4.MatchString(s) {
		return
	} else if regexHostname.MatchString(s) {
		return
	} else {
		err = errors.New("Not a valid IPv4 address or URL") //lint:ignore ST1005 will be displayed in UI
	}
	return
}

func integerValidator(s string) (err error) {
	var port int
	port, err = strconv.Atoi(s)
	if err != nil {
		err = errors.New("Not a valid integer") //lint:ignore ST1005 will be displayed in UI
	} else if port <= 0 {
		err = errors.New("Must be greater than zero") //lint:ignore ST1005 will be displayed in UI
	}
	return
}

func requiredValidator(s string) (err error) {
	if s == "" {
		err = errors.New("Cannot be empty") //lint:ignore ST1005 will be displayed in UI
	}
	return
}
