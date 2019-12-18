//+build !headless

package gui

import (
	"github.com/p9c/gio-parallel/app"
	"github.com/p9c/pod/cmd/gui/duoui"
	"github.com/p9c/pod/pkg/log"
)

func WalletGUI(duo *duoui.DuoUI) (err error) {
	go func() {
		if err := duoui.DuoUIloop(duo); err != nil {
			log.FATAL(err)
		}
	}()
	app.Main()
	return
}
