package root

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

type footer struct {
	widget.BaseWidget

	iconDisconnected *theme.ErrorThemedResource
	iconConnected    *theme.ThemedResource

	labelConnectionStatus *canvas.Text
	iconConnectionStatus  *widget.Icon
}

func newFooter() *footer {
	iconDisconnected := theme.NewErrorThemedResource(theme.CancelIcon())
	iconConnected := theme.NewThemedResource(theme.ConfirmIcon())
	iconConnected.ColorName = theme.ColorNameSuccess
	labelConnectionStatus := &canvas.Text{
		Text:     constants.STATUS_DISCONNECTED,
		Color:    color.Gray{Y: 128},
		TextSize: theme.CaptionTextSize(),
	}
	iconConnectionStatus := widget.NewIcon(iconDisconnected)
	f := &footer{
		iconDisconnected:      iconDisconnected,
		iconConnected:         iconConnected,
		labelConnectionStatus: labelConnectionStatus,
		iconConnectionStatus:  iconConnectionStatus,
	}
	f.ExtendBaseWidget(f)
	return f
}

func (f *footer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewHBox(
			layout.NewSpacer(),
			f.labelConnectionStatus,
			f.iconConnectionStatus,
		),
	)
}

func (f *footer) SetConnected(v bool) {
	if v {
		f.labelConnectionStatus.Text = constants.STATUS_CONNECTED
		f.iconConnectionStatus.SetResource(f.iconConnected)
	} else {
		f.labelConnectionStatus.Text = constants.STATUS_DISCONNECTED
		f.iconConnectionStatus.SetResource(f.iconDisconnected)
	}
}
