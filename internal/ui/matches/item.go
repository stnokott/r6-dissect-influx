package matches

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/internal/game"
)

type matchItem struct {
	widget.BaseWidget

	rootCard *widget.Card

	lblMatchID   *widget.Label
	lblMapName   *widget.Label
	lblGameMode  *widget.Label
	lblNumRounds *widget.Label
}

func newMatchItem() *matchItem {
	mi := &matchItem{
		rootCard:     widget.NewCard("<placeholder>", "<placeholder>", layout.NewSpacer()),
		lblMatchID:   widget.NewLabel("<placeholder>"),
		lblMapName:   widget.NewLabel("<placeholder>"),
		lblGameMode:  widget.NewLabel("<placeholder>"),
		lblNumRounds: widget.NewLabel("<placeholder>"),
	}
	mi.ExtendBaseWidget(mi)
	return mi
}

func (mi *matchItem) CreateRenderer() fyne.WidgetRenderer {
	mi.rootCard.SetContent(
		widget.NewForm(
			widget.NewFormItem("Match ID:", mi.lblMatchID),
			widget.NewFormItem("Game Mode:", mi.lblGameMode),
			widget.NewFormItem("Rounds played:", mi.lblNumRounds),
		),
	)
	return widget.NewSimpleRenderer(mi.rootCard)
}

func (mi *matchItem) Load(match []*game.RoundInfo) {
	if len(match) == 0 {
		mi.rootCard.SetTitle("No data")
		mi.rootCard.SetSubTitle("No data")
		mi.lblMatchID.SetText("No data")
		mi.lblGameMode.SetText("No data")
		mi.lblNumRounds.SetText("No data")
	} else {
		firstRound := match[0]
		mi.rootCard.SetTitle(firstRound.MapName)
		mi.rootCard.SetSubTitle(firstRound.Time.Local().Format(time.RFC850))
		mi.lblMatchID.SetText(firstRound.MatchID)
		mi.lblGameMode.SetText(firstRound.MatchType + " - " + firstRound.GameMode)
		mi.lblNumRounds.SetText(strconv.Itoa(len(match)))
	}
}
