package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/yiitz/iceapple/media"
	"github.com/yiitz/iceapple/entity"
)

type PlayList struct {
	items      []*entity.Song
	list       *tview.List
	player     *media.Player
	app        *tview.Application
	Selectable bool
}

func NewPlayList(app *tview.Application, list *tview.List, player *media.Player) *PlayList {
	p := PlayList{list: list, player: player, app: app}

	list.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		if p.Selectable {
			item := p.items[i]
			p.player.Play(item)
		}
	})
	return &p
}

func (pl *PlayList) SetItems(items []*entity.Song) {
	pl.app.Lock()
	pl.list.Clear()
	pl.items = items
	for i, v := range items {
		pl.list.AddItem(fmt.Sprintf("[%d] %s - %s - %s", i+1, v.Name, v.Artist, v.Album), "", 0, nil)
	}
	pl.app.Unlock()
	pl.app.Draw()
}
