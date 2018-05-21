package app

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/yiitz/iceapple/entity"
	"github.com/yiitz/iceapple/log"
	"github.com/yiitz/iceapple/media"
	"github.com/yiitz/iceapple/ui"
)

func Run() {

	player := media.NewPlayer()

	//ui
	app := tview.NewApplication()

	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	mainContent := tview.NewPages()

	list := tview.NewList().ShowSecondaryText(false)
	pl := ui.NewPlayList(app, list, player)
	pl.SetItems([]ui.PlayListItem{&entity.Song{Name: "后会无期", Uri: "file:///home/yiitz/Downloads/hhwq.mp3"},
		&entity.Song{Name: "example", Uri: "https://www.freedesktop.org/software/gstreamer-sdk/data/media/sintel_trailer-480p.webm"}})
	mainContent.AddPage("main", list, true, true)

	flex.AddItem(mainContent, 0, 1, true)

	progress := tview.NewTextView()
	status := tview.NewTextView()
	pb := ui.NewPlayBar(app, progress, status, player)
	flex.AddItem(progress, 1, 0, false)
	flex.AddItem(status, 1, 0, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k, r := event.Key(), event.Rune()
		log.LoggerRoot.Debugf("key event:%d,%d", k, r)
		if k <= tcell.KeyRune {
			switch r {
			case ' ':
				player.TriggerPlay()
				return nil
			}
		} else {
			switch k {
			case 256:
			case tcell.KeyUp:
				if event.Modifiers()&tcell.ModCtrl != 0 {
					player.VolumeUp()
					pb.Draw(true)
					return nil
				}
			case tcell.KeyDown:
				if event.Modifiers()&tcell.ModCtrl != 0 {
					player.VolumeDown()
					pb.Draw(true)
					return nil
				}
			}
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

	pb.CancelTimer()
	player.Stop()
	media.ClosePlayer(player)
}
