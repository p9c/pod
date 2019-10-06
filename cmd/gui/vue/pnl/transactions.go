package pnl

import "github.com/p9c/pod/cmd/gui/mod"

func TransactionsExcerpts() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    true,
		Name:     "Transactions",
		ID:       "paneltxsex",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "transactions",
		Js: `
	data () { return { 
		duOSys,
		pageSettings: { pageSize: 5 },
	}},
  methods: {
     
},

		`,
		Template: `<div id="paneltxsex" class="Txs"><div class="rwrap"><ejs-grid :dataSource="this.duOSys.txsEx.txs" height="100%" :pageSettings="pageSettings"><e-columns><e-column field="category" headerText="Category" textAlign="Right" width=90></e-column><e-column field="time" headerText="Time" format="auto"  textAlign="Right" width=90></e-column><e-column field="amount" headerText="Amount" textAlign="Right" width=90></e-column></e-columns></ejs-grid></div></div>`,
		Css: `
		
		`,
	}
}

func Transactions() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    true,
		Name:     "Transactions",
		ID:       "paneltransactions",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "transactions",
		Js: `
	data () { return { 
		duOSys,
		pageSettings: { pageSize: 10, pageSizes: [10,20,50,100], pageCount: 5 },
      	ddldata: ['All', 'generated', 'sent', 'received', 'immature']
	}},
  methods: {
     
},

		`,
		Template: `<div class="rwrap"><div class="select-wrap"><ejs-dropdownlist id="ddlelement" :dataSource="ddldata" placeholder="Select category to filter"></ejs-dropdownlist></div><ejs-grid :dataSource="this.duOSys.txsEx.txs" height="100%" :allowPaging="true" :pageSettings="pageSettings"><e-columns><e-column field="category" headerText="Category" textAlign="Right" width=90></e-column><e-column field="time" headerText="Time" format="unix"  textAlign="Right" width=90></e-column><e-column field="txid" headerText="TxID" textAlign="Right" width=240></e-column><e-column field="amount" headerText="Amount" textAlign="Right" width=90></e-column></e-columns></ejs-grid></div>`,
		Css: `
		
		`,
	}
}

func TimeBalance() mod.DuOScomp {
	return mod.DuOScomp{
		IsApp:    true,
		Name:     "Balance",
		ID:       "paneltimebalance",
		Version:  "0.0.1",
		CompType: "panel",
		SubType:  "transactions",
		Js: `
  data: function() {
    return {
		duOSys,
      	theme: "Material",
      primaryXAxis: {
		valueType: "DateTime",
        intervalType: "Auto",
        edgeLabelPlacement: "Shift",
        majorGridLines: { width: 0 }
      },
      primaryYAxis: {
        labelFormat: "{value} DUO",
        rangePadding: "None",
        minimum: 0,
        maximum: system.data.d.txsex.balanceheight,
        interval: 20,
        lineStyle: { width: 0 },
        majorTickLines: { width: 0 },
        minorTickLines: { width: 0 }
      },
      chartArea: {
        border: {
          width: 0
        }
      },
      marker: {
        visible: true,
        height: 10,
        width: 10
      },
      tooltip: {
        enable: true
      },
      title: "Inflation - Consumer Price"
    	}
  }

		`,
		Template: `<div class="rwrap">
  <div class="control-section">
    <div align="center">
        <ejs-chart style="display:block" :theme="this.theme" align="center" id="chartcontainer" :title="this.title" :primaryXAxis="this.primaryXAxis" :primaryYAxis="this.primaryYAxis'
            :tooltip="this.tooltip" :chartArea="this.chartArea" width="100%" height="100%">
            <e-series-collection>
                <e-series :dataSource="system.data.d.txsex.txs" type="Line" xName="time" yName="balance" name="Germany" width=2 :marker="this.marker"> </e-series>
            </e-series-collection>
        </ejs-chart>
    </div>
</div>
</div>`,
		Css: `
		
		`,
	}
}
