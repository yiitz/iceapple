package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/yiitz/iceapple/media"
)

type PlayListItem interface {
	GetName() string
	GetUri() string
}

type PlayList struct {
	items  []PlayListItem
	list   *tview.List
	player *media.Player
	app    *tview.Application
}

func NewPlayList(app *tview.Application, list *tview.List, player *media.Player) *PlayList {
	p := PlayList{list: list, player: player, app: app}

	list.SetSelectedFunc(func(i int, s string, s2 string, r rune) {
		item := p.items[i]
		p.player.Play(item.GetUri())
	})
	return &p
}

func (pl *PlayList) SetItems(items []PlayListItem) {
	pl.list.Clear()
	pl.items = items
	for i, v := range items {
		pl.list.AddItem(fmt.Sprintf("[%d] %s", i+1, v.GetName()), "", 0, nil)
	}
}
