package component

import (
	"fmt"
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/model"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/cmd/gui/rcd"
)

var (
	transList = &layout.List{
		Axis: layout.Vertical,
	}
	allTxs      = new(controller.CheckBox)
	mintedTxs   = new(controller.CheckBox)
	immatureTxs = new(controller.CheckBox)
	sentTxs     = new(controller.CheckBox)
	receivedTxs = new(controller.CheckBox)
)

func TransactionsList(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		transList.Layout(gtx, len(rc.Status.Wallet.Transactions.Txs), func(i int) {
			t := rc.Status.Wallet.Transactions.Txs[i]
			HorizontalLine(gtx, 1, th.Color.Hint)()
			layout.Flex{
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				layout.Rigid(txsDetails(gtx, th, i, &t)),
				layout.Rigid(Label(gtx, th, th.Font.Mono, 12, th.Color.Dark, fmt.Sprintf("%0.8f", t.Amount))))
		})
	}
}

func txsDetails(gtx *layout.Context, th *theme.DuoUItheme, i int, t *model.DuoUItx) func() {
	return func() {
		layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(Label(gtx, th, th.Font.Primary, 12, th.Color.Dark, fmt.Sprint(i))),
			layout.Rigid(Label(gtx, th, th.Font.Primary, 12, th.Color.Dark, t.TxID)),
			layout.Rigid(Label(gtx, th, th.Font.Primary, 12, th.Color.Dark, fmt.Sprintf("%0.8f", t.Amount))),
			layout.Rigid(Label(gtx, th, th.Font.Primary, 12, th.Color.Dark, t.Category)),
			layout.Rigid(Label(gtx, th, th.Font.Primary, 12, th.Color.Dark, t.Time)),
		)
	}
}

func TransactionsFilter(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme) func() {
	return func() {
		layout.Flex{}.Layout(gtx,
			layout.Rigid(txsFilterItem(gtx, th, "ALL", allTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "MINTED", mintedTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "IMATURE", immatureTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "SENT", sentTxs)),
			layout.Rigid(txsFilterItem(gtx, th, "RECEIVED", receivedTxs)))
	}
}

func txsFilterItem(gtx *layout.Context, th *theme.DuoUItheme, id string, c *controller.CheckBox) func() {
	return func() {
		th.DuoUIcheckBox(id, th.Color.Light, th.Color.Light).Layout(gtx, c)
	}
}
