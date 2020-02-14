package duoui

import (
	"errors"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gui/io/system"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func DuoUImainLoop(d *model.DuoUI, r *rcd.RcVar) error {
	ui := new(DuoUI)
	ui = &DuoUI{
		ly: d,
		rc: r,
	}
	//ui.ly.Pages = ui.LoadPages()
	for {
		select {
		case <-ui.ly.Ready:
			updateTrigger := make(chan struct{}, 1)
			go func() {
			quitTrigger:
				for {
					select {
					case <-updateTrigger:
						log.DEBUG("repaint forced")
						ui.ly.Window.Invalidate()
					case <-ui.ly.Quit:
						break quitTrigger
					}
				}
			}()
			ui.rc.ListenInit(updateTrigger)
			ui.ly.IsReady = true
		case <-ui.ly.Quit:
			log.DEBUG("quit signal received")
			interrupt.Request()
			// This case is for handling when some external application is controlling the GUI and to gracefully
			// handle the back-end servers being shut down by the interrupt library receiving an interrupt signal
			// Probably nothing needs to be run between starting it and shutting down
			<-interrupt.HandlersDone
			log.DEBUG("closing GUI from interrupt/quit signal")
			return errors.New("shutdown triggered from back end")
		case e := <-ui.ly.Window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				log.DEBUG("destroy event received")
				interrupt.Request()
				// Here do cleanup like are you sure (optional) modal or shutting down indefinite spinner
				<-interrupt.HandlersDone
				return e.Err
			case system.FrameEvent:
				ui.ly.Context.Reset(e.Config, e.Size)
				ui.ly.Pages = ui.LoadPages()
				//if rc.Boot.IsBoot {
				//d.DuoUImainScreen()
				//e.Frame(d.mod.Context.Ops)
				//} else {
				//	d.mod.Context.Reset(e.Config, e.Size)
				//	if rc.Boot.IsFirstRun {
				//		//DuoUIloaderCreateWallet(duo.m, cx, rc)
				//	} else {
				ui.DuoUImainScreen()
				if ui.rc.Dialog.Show {
					ui.DuoUIdialog()
				}
				//d.DuoUItoastSys()
				//
				//	}
				e.Frame(ui.ly.Context.Ops)
				ui.ly.Context.Reset(e.Config, e.Size)
				//}
			}
		}
	}
}
