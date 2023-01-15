package config

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/utils"
)

const (
	WINDOW_TITLE       string = "Settings"
	hintMatchReplayDir string = "Directory where your match replays are stored. Usually in your game folder named 'MatchReplays'"
	hintInfluxHost     string = `Host of your InfluxDB instance (IP or hostname, no "http(s)" required)`
	hintInfluxPort     string = "Port of your InfluxDB instance"
	hintInfluxBucket   string = "InfluxDB bucket where the data should be pushed"
)

func ShowDialog(parent fyne.Window) {
	inputMatchReplayDir := &widget.Entry{Validator: validateMatchReplayDir}
	buttonAutodetectMatchReplayDir := widget.NewButton("Autodetect", func() {
		folder, err := matchReplayFolderFromRegistry()
		if err == nil {
			inputMatchReplayDir.SetText(folder)
		} else {
			dialog.NewInformation("Autodetection failed", utils.TitleErr(err, true), parent).Show()
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
	inputInfluxHost := &widget.Entry{Validator: validateInfluxHost}
	inputInfluxPort := &widget.Entry{Text: "8086", Validator: validateInfluxPort}
	inputInfluxBucket := &widget.Entry{Validator: validateInfluxBucket}

	formItems := []*widget.FormItem{
		{Text: "Match replay directory", Widget: inputMatchReplayDir, HintText: hintMatchReplayDir},
		{Text: "", Widget: container.NewGridWithColumns(2, buttonAutodetectMatchReplayDir, buttonSelectMatchReplayDir)},
		{Text: "InfluxDB host", Widget: inputInfluxHost, HintText: hintInfluxHost},
		{Text: "InfluxDB port", Widget: inputInfluxPort, HintText: hintInfluxPort},
		{Text: "InfluxDB bucket", Widget: inputInfluxBucket, HintText: hintInfluxBucket},
	}

	d := dialog.NewForm(
		WINDOW_TITLE,
		"Confirm",
		"Cancel",
		formItems,
		func(ok bool) {
			log.Println("TODO: handle")
		},
		parent,
	)
	d.Resize(d.MinSize().Add(fyne.NewDelta(200, 0)))
	d.Show()
}

// TODO: warn if no matches found or none could be parsed
func validateMatchReplayDir(s string) (err error) {
	if s == "" {
		err = errors.New("Cannot be empty") //lint:ignore ST1005 will be displayed in UI
	} else {
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
	}
	return
}

var (
	regexIPv4     = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
	regexHostname = regexp.MustCompile(`^(?:[0-9a-zA-Z]+\.)+[0-9a-zA-Z]{2,4}$`)
)

func validateInfluxHost(s string) (err error) {
	if s == "" {
		err = errors.New("Cannot be empty") //lint:ignore ST1005 will be displayed in UI
	} else if regexIPv4.MatchString(s) {
		return
	} else if regexHostname.MatchString(s) {
		return
	} else {
		err = errors.New("Not a valid IPv4 address or URL") //lint:ignore ST1005 will be displayed in UI
	}
	return
}

func validateInfluxPort(s string) (err error) {
	if s == "" {
		err = errors.New("Cannot be empty") //lint:ignore ST1005 will be displayed in UI
	} else {
		var port int
		port, err = strconv.Atoi(s)
		if err != nil {
			err = errors.New("Not a valid integer") //lint:ignore ST1005 will be displayed in UI
		} else if port <= 0 {
			err = errors.New("Must be greater than zero") //lint:ignore ST1005 will be displayed in UI
		}
	}
	return
}

func validateInfluxBucket(s string) (err error) {
	if s == "" {
		err = errors.New("Cannot be empty") //lint:ignore ST1005 will be displayed in UI
	}
	return
}
