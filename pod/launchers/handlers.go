package launchers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/p9c/matrjoska/cmd/ctl"
	"github.com/p9c/matrjoska/cmd/node/node"
	"github.com/p9c/matrjoska/pkg/constant"
	"github.com/p9c/matrjoska/pkg/pod"

	"github.com/p9c/qu"

	"github.com/p9c/matrjoska/pkg/apputil"
	"github.com/p9c/matrjoska/walletmain"
)

// NodeHandle runs the ParallelCoin blockchain node
func NodeHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	I.Ln("running node handler")
	cx.NodeReady = qu.T()
	cx.Node.Store(false)
	// // serviceOptions defines the configuration options for the daemon as a service on Windows.
	// type serviceOptions struct {
	// 	ServiceCommand string `short:"s" long:"service" description:"Service command {install, remove, start, stop}"`
	// }
	// // runServiceCommand is only set to a real function on Windows. It is used to parse and execute service commands
	// // specified via the -s flag.
	// runServiceCommand := func(string) (e error) { return nil }
	// // Service options which are only added on Windows.
	// serviceOpts := serviceOptions{}
	// // Perform service command and exit if specified. Invalid service commands show an appropriate error. Only runs
	// // on Windows since the runServiceCommand function will be nil when not on Windows.
	// if serviceOpts.ServiceCommand != "" && runServiceCommand != nil {
	// 	if e = runServiceCommand(serviceOpts.ServiceCommand); E.Chk(e) {
	// 		return e
	// 	}
	// 	return nil
	// }
	go func() {
		if e := node.Main(cx); E.Chk(e) {
			E.Ln("error starting node ", e)
		}
	}()
	I.Ln("starting node")
	if cx.Config.DisableRPC.False() {
		cx.RPCServer = <-cx.NodeChan
		cx.NodeReady.Q()
		cx.Node.Store(true)
		I.Ln("node started")
	}
	// }
	cx.WaitWait()
	I.Ln("node is now fully shut down")
	cx.WaitGroup.Wait()
	<-cx.KillAll
	return nil
}

// walletHandle runs the wallet server
func walletHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	cx.Config.WalletFile.Set(filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, constant.DbName))
	// dbFilename := *cx.Config.DataDir + slash + cx.ActiveNet.
	// 	Params.Name + slash + wallet.WalletDbName
	if !apputil.FileExists(cx.Config.WalletFile.V()) && !cx.IsGUI {
		// D.Ln(cx.ActiveNet.Name, *cx.Config.WalletFile)
		if e = walletmain.CreateWallet(cx.ActiveNet, cx.Config); E.Chk(e) {
			E.Ln("failed to create wallet", e)
			return e
		}
		fmt.Println("restart to complete initial setup")
		os.Exit(0)
	}
	// for security with apps launching the wallet, the public password can be set with a file that is deleted after
	walletPassPath := filepath.Join(cx.Config.DataDir.V(), cx.ActiveNet.Name, "wp.txt")
	D.Ln("reading password from", walletPassPath)
	if apputil.FileExists(walletPassPath) {
		var b []byte
		if b, e = ioutil.ReadFile(walletPassPath); !E.Chk(e) {
			cx.Config.WalletPass.SetBytes(b)
			D.Ln("read password '" + string(b) + "'")
			for i := range b {
				b[i] = 0
			}
			if e = ioutil.WriteFile(walletPassPath, b, 0700); E.Chk(e) {
			}
			if e = os.Remove(walletPassPath); E.Chk(e) {
			}
			D.Ln("wallet cookie deleted", *cx.Config.WalletPass)
		}
	}
	cx.WalletKill = qu.T()
	if e = walletmain.Main(cx); E.Chk(e) {
		E.Ln("failed to start up wallet", e)
	}
	// if !*cx.Config.DisableRPC {
	// 	cx.WalletServer = <-cx.WalletChan
	// }
	// cx.WaitGroup.Wait()
	cx.WaitWait()
	return
}

func CtlHandleList(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	_ = cx
	ctl.ListCommands()
	return nil
}

func CtlHandle(ifc interface{}) (e error) {
	var cx *pod.State
	var ok bool
	if cx, ok = ifc.(*pod.State); !ok {
		return fmt.Errorf("cannot run without a state")
	}
	cx.Config.LogLevel.Set("off")
	ctl.Main(cx)
	return nil
}
