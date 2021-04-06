package pod

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/p9c/monorepo/pkg/chaincfg"
	"github.com/p9c/monorepo/pkg/fork"
	"github.com/p9c/monorepo/pkg/log"
	"github.com/p9c/monorepo/pkg/opts"
	"github.com/p9c/monorepo/pkg/pod"
	"os"
	
	"github.com/p9c/monorepo/pkg/interrupt"
	
	"github.com/p9c/monorepo/pkg/podconfig"
	
	"github.com/p9c/monorepo/cmd/kopach"
)

// kopachHandle runs the kopach miner
func KopachHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	log.AppColorizer = color.Bit24(255, 128, 128, false).Sprint
	log.App = "kopach"
	opts.I.Ln("starting up kopach standalone miner for parallelcoin")
	opts.D.Ln(os.Args)
	podconfig.Configure(cx, true)
	if cx.ActiveNet.Name == chaincfg.TestNet3Params.Name {
		fork.IsTestnet = true
	}
	defer cx.KillAll.Q()
	e = kopach.Run(cx)
	<-interrupt.HandlersDone
	opts.D.Ln("kopach main finished")
	return
}
