package controller

import (
	"context"
	chain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/controller/broadcast"
	"github.com/p9c/pod/pkg/controller/gcm"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	"go.uber.org/atomic"
	"sync"
)

type Blocks struct {
	PrevBlock    *chainhash.Hash
	Bits         uint32
	Transactions []*wire.MsgTx
	Listeners    []string
}

var (
	WorkMagic = [4]byte{'w', 'o', 'r', 'k'}
)

// Blocks is a block broadcast message for miners to mine from
type Blocks struct {
	// New is a flag that distinguishes a newly accepted/connected block from a rebroadcast
	New bool
	// Payload is a map of bytes indexed by block version number
	Payload map[int32][]byte
}

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	var ctx context.Context
	var active atomic.Bool
	ctx, cancel = context.WithCancel(context.Background())
	if len(*cx.Config.RPCListeners) < 1 || *cx.Config.DisableRPC {
		log.WARN("not running controller without RPC enabled")
		cancel()
		return
	}
	if len(*cx.Config.Listeners) < 1 || *cx.Config.DisableListen {
		log.WARN("not running controller without p2p listener enabled")
		cancel()
		return
	}
	messageBase := GetMessageBase(cx)
	//log.SPEW(messageBase.CreateContainer(WorkMagic))
	go func() {
		// There is no unsubscribe but we can use an atomic to disable the
		// function instead - this also ensures that new work doesn't start
		// once the context is cancelled below
		active.Store(true)
		var subMx sync.Mutex
		log.DEBUG("miner controller starting")
		cx.RealNode.Chain.Subscribe(func(n *chain.Notification) {
			if active.Load() {
				// first to arrive locks out any others while processing
				switch n.Type {
				case chain.NTBlockAccepted:
					subMx.Lock()
					defer subMx.Unlock()
					log.DEBUG("received new chain notification")
					// construct work message
					//log.SPEW(n)
					mB, ok := n.Data.(*util.Block)
					if !ok {
						log.WARN("chain accepted notification is not a block")
						break
					}
					msg := Serializers{}
					msg = append(msg, messageBase...)
					h := NewHash()
					h.PutHash(mB.MsgBlock().Header.BlockHash())
					msg = append(msg, h)
					tip := cx.RealNode.Chain.BestChain.Tip()
					bM := map[int32]uint32{}
					bitsMap := &bM
					var err error
					tip.DiffMx.Lock()
					if tip.Diffs == nil ||
						len(*tip.Diffs) != len(fork.List[1].AlgoVers) {
						tip.DiffMx.Unlock()
						bitsMap, err = cx.RealNode.Chain.
							CalcNextRequiredDifficultyPlan9Controller(tip)
						if err != nil {
							log.ERROR(err)
							return
						}
					} else {
						bitsMap = tip.Diffs
						tip.DiffMx.Unlock()
					}
					bitses := NewBitses()
					bitses.PutBitses(*bitsMap)
					msg = append(msg, bitses)
					txs := mB.MsgBlock().Transactions
					for i := range txs {
						t := &Transaction{}
						t.PutTx(txs[i])
						msg = append(msg, t)
					}
					srs := msg.CreateContainer(WorkMagic)
					// send out srs.Data

					ip := NewIPs()
					ip.Decode(srs.Get(0))
					log.DEBUG(ip.GetIPs())
					listener := NewPort()
					listener.Decode(srs.Get(1))
					log.DEBUG(listener.GetUint16())
					rpcListener := NewPort()
					rpcListener.Decode(srs.Get(2))
					log.DEBUG(rpcListener.GetUint16())
					ctrlrListener := NewPort()
					ctrlrListener.Decode(srs.Get(3))
					log.DEBUG(ctrlrListener.GetUint16())
					prevH := NewHash()
					prevH.Decode(srs.Get(4))
					log.DEBUG(prevH.GetHash())
					bt := NewBitses()
					bt.Decode(srs.Get(5))
					log.DEBUG(bt.GetBitses())
					txn := NewTransaction()
					txn.Decode(srs.Get(6))
					log.SPEW(txn.GetTx())
				}
			}
		})
		select {
		case <-ctx.Done():
			log.DEBUG("miner controller shutting down")
			active.Store(false)
			break
		}
	}()
	return
}
