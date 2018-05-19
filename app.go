package main

import (
	"github.com/yiitz/iceapple/player"
	"github.com/gizak/termui"
)

func main()  {
	player.Init()
	player.Start()
	p := player.NewPlayer()

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	strs := []string{
		"[1] 后会无期",
		"[2] sintel trailer"}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = 7
	ls.Width = 25
	ls.Y = 0

	termui.Render(ls)
	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		p.Stop()
		player.Stop()
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/<space>", func(termui.Event) {
		p.TriggerPlay()
	})
	termui.Handle("/sys/kbd/1", func(termui.Event) {
		p.Play("file:///home/yiitz/Downloads/hhwq.mp3")
	})
	termui.Handle("/sys/kbd/2", func(termui.Event) {
		p.Play("https://www.freedesktop.org/software/gstreamer-sdk/data/media/sintel_trailer-480p.webm")
	})
	termui.Loop()
}