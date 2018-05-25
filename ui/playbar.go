package ui

import (
	"fmt"
	"github.com/alex023/clock"
	"github.com/rivo/tview"
	"github.com/yiitz/iceapple/media"
	"github.com/yiitz/iceapple/timer"
	"math"
	"strings"
	"time"
	"github.com/yiitz/iceapple/log"
)

type PlayBar struct {
	progress       *tview.TextView
	status         *tview.TextView
	app            *tview.Application
	player         *media.Player
	job            clock.Job
	OnSongFinished func()
}

func NewPlayBar(app *tview.Application, progress *tview.TextView, status *tview.TextView, player *media.Player) *PlayBar {
	pb := PlayBar{
		app:      app,
		progress: progress,
		status:   status,
		player:   player,
	}

	player.OnPlayStart = func() {
		pb.CancelTimer()
		pb.job, _ = timer.Timer.AddJobRepeat(time.Second, math.MaxUint64, func() {
			pb.Draw(false)
		})
	}
	return &pb
}

func (pb *PlayBar) Draw(force bool) {
	state := pb.player.GetState()
	log.LoggerRoot.Debugf("player state: %d",state)

	finished := false
	if state == media.GstStatePlaying || force || pb.job != nil {
		_, _, w, _ := pb.progress.GetInnerRect()
		w -= 3

		position, duration := pb.player.GetProgress()
		progress := int(float64(position) / float64(duration) * float64(w))

		if state == media.GstStateNull {
			pb.CancelTimer()
			pb.progress.SetText("[" + strings.Repeat("-", w+1) + "]")
			finished = true
		} else {
			pb.progress.SetText(strings.Repeat("=", progress) + ">" + strings.Repeat("-", w-progress))
		}

		pb.status.SetText(
			fmt.Sprintf("[time %02d:%02d:%02d/%02d:%02d:%02d]\t[volume %d%%]",
				int(position.Hours()), int(position.Minutes())%60, int(position.Seconds())%60,
				int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60, pb.player.GetVolume()))

		pb.app.Draw()
	}

	if finished && pb.OnSongFinished != nil {
		pb.OnSongFinished()
	}
}

func (pb *PlayBar) CancelTimer() {
	if pb.job != nil {
		pb.job.Cancel()
		pb.job = nil
	}
}
