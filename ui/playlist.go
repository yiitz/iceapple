package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/yiitz/iceapple/media"
)

type PlayListItem interface {
	GetUri() string
	GetName() string
	GetArtist() string
	GetAlbum() string
}

type PlayList struct {
	items      []PlayListItem
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
			p.player.Play(item.GetUri())
		}
	})
	return &p
}

func (pl *PlayList) SetItems(items []PlayListItem) {
	pl.list.Clear()
	pl.items = items
	for i, v := range items {
		pl.list.AddItem(fmt.Sprintf("[%d] %s - %s - %s", i+1, v.GetName(), v.GetArtist(), v.GetAlbum()), "", 0, nil)
	}
}
