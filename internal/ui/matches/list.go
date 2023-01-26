package matches

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/internal/game"
)

type MatchListView struct {
	widget.BaseWidget

	list *widget.List

	// matches holds the match data used by our list.
	// Since the Fyne list widget can't work with a list, we use a map (see below)
	// to associate the match IDs.
	matches [][]*game.RoundInfo

	// matchIDs maps match IDs to their indexes in the list.
	matchIDs map[string]int

	mutexMatches sync.Mutex
}

// TODO: switch to scrollable Grid with 2 cols

func NewMatchListView() *MatchListView {
	l := &MatchListView{
		matches:      [][]*game.RoundInfo{},
		matchIDs:     map[string]int{},
		mutexMatches: sync.Mutex{},
	}

	l.ExtendBaseWidget(l)
	return l
}

func (l *MatchListView) Add(r *game.RoundInfo) {
	l.mutexMatches.Lock()
	if matchIndex, exists := l.matchIDs[r.MatchID]; exists {
		// match for this round already present in data
		l.matches[matchIndex] = append(l.matches[matchIndex], r)
	} else {
		// match not present yet
		l.matches = append(l.matches, []*game.RoundInfo{r})
		l.matchIDs[r.MatchID] = len(l.matches) - 1
	}
	l.list.Refresh()
	l.list.ScrollToBottom()
	l.mutexMatches.Unlock()
}

func (l *MatchListView) CreateRenderer() fyne.WidgetRenderer {
	l.list = widget.NewList(
		func() int {
			return len(l.matches)
		},
		func() fyne.CanvasObject {
			return newMatchItem()
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*matchItem).Load(l.matches[lii])
		},
	)
	return widget.NewSimpleRenderer(l.list)
}
