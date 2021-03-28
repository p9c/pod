package podcfg

import (
	"fmt"
	"github.com/p9c/pod/pkg/appdata"
	"github.com/p9c/pod/pkg/chaincfg"
	uberatomic "go.uber.org/atomic"
	"math/rand"
	"net"
	"path/filepath"
	"sync/atomic"
	"time"
)

func (c *Config) getAllOptionStrings() (s map[string][]string, e error) {
	s = make(map[string][]string)
	if c.ForEach(func(ifc Option) bool {
		md := ifc.GetMetadata()
		if _, ok := s[ifc.Name()]; ok {
			e = fmt.Errorf("conflicting option names: %v %v", ifc.GetAllOptionStrings(), s[ifc.Name()])
			return false
		}
		s[ifc.Name()] = md.GetAllOptionStrings()
		return true
	},
	) {
	}
	s["commandslist"] = c.Commands.GetAllCommands()
	// I.S(s["commandslist"])
	return
}

func findConflictingItems(valOpts map[string][]string) (o []string, e error) {
	var ss, ls string
	for i := range valOpts {
		for j := range valOpts {
			// W.Ln(s[i], s[j], i==j, s[i]==s[j])
			if i == j {
				continue
			}
			a := valOpts[i]
			b := valOpts[j]
			for ii := range a {
				for jj := range b {
					if ii == jj {
						continue
					}
					// W.Ln(i == j, s[i] == s[j])
					// I.Ln(s[i], s[j])
					ss, ls = shortestString(a[ii], b[jj])
					// I.Ln("these should not be the same string", ss, ls)
					if ss == ls[:len(ss)] {
						E.F("conflict between %s and %s, ", ss, ls)
						o = append(o, ss, ls)
					}
				}
			}
		}
	}
	if len(o) > 0 {
		panic(fmt.Sprintf("conflicts found: %v", o))
	}
	return
}

func shortestString(a, b string) (s, l string) {
	switch {
	case len(a) > len(b):
		s, l = b, a
	default:
		s, l = a, b
	}
	return
}

// GetDefaultConfig returns a Config struct pristine factory freshaoeu
func GetDefaultConfig() (c *Config) {
	network := "mainnet"
	rand.Seed(time.Now().Unix())
	var datadir = &atomic.Value{}
	datadir.Store([]byte(appdata.Dir(Name, false)))
	c = &Config{
		Commands: Commands{
			{Name: "gui", Description:
			"ParallelCoin GUI Wallet/Miner/Explorer",
				Entrypoint: func(c *Config) error { return nil },
			},
			{Name: "version", Description:
			"print version and exit",
				Entrypoint: func(c *Config) error { return nil },
			},
			{Name: "ctl", Description:
			"command line wallet and chain RPC client",
				Entrypoint: func(c *Config) error { return nil },
			},
			{Name: "node", Description:
			"ParallelCoin blockchain node",
				Entrypoint: func(c *Config) error { return nil },
				Commands: []Command{
					{Name: "dropaddrindex", Description:
					"drop the address database index",
						Entrypoint: func(c *Config) error { return nil },
					},
					{Name: "droptxindex", Description:
					"drop the transaction database index",
						Entrypoint: func(c *Config) error { return nil },
					},
					{Name: "dropcfindex", Description:
					"drop the cfilter database index",
						Entrypoint: func(c *Config) error { return nil },
					},
					{Name: "dropindexes", Description:
					"drop all of the indexes",
						Entrypoint: func(c *Config) error { return nil },
					},
					{Name: "resetchain", Description:
					"deletes the current blockchain cache to force redownload",
						Entrypoint: func(c *Config) error { return nil },
					},
				},
			},
			{Name: "wallet", Description:
			"run the wallet server (requires a chain node to function)",
				Entrypoint: func(c *Config) error { return nil },
				Commands: []Command{
					{Name: "drophistory", Description:
					"reset the wallet transaction history",
						Entrypoint: func(c *Config) error { return nil },
					},
				},
			},
			{Name: "kopach", Description:
			"standalone multicast miner for easy mining farm deployment",
				Entrypoint: func(c *Config) error { return nil },
			},
			{Name: "worker", Description:
			"single thread worker process, normally started by kopach",
				Entrypoint: func(c *Config) error { return nil },
			},
		},
		AddCheckpoints: NewStrings(Metadata{
			Option:  "addcheckpoint",
			Aliases: []string{"ac"},
			Group:   "debug",
			Label:   "Add Checkpoints",
			Description:
			"add custom checkpoints",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		AddPeers: NewStrings(Metadata{
			Option:  "addpeer",
			Aliases: []string{"ap"},
			Group:   "node",
			Label:   "Add Peers",
			Description:
			"manually adds addresses to try to connect to",
			Type:   "ipaddress",
			Widget: "multi",
			// Hook:        "addpeer",
			OmitEmpty: true,
		},
			[]string{},
			// []string{"127.0.0.1:12345", "127.0.0.1:12345", "127.0.0.1:12345", "127.0.0.1:12344"},
		),
		AddrIndex: NewBool(Metadata{
			Option:  "addrindex",
			Aliases: []string{"ai"},
			Group:   "node",
			Label:   "Address Index",
			Description:
			"maintain a full address-based transaction index which makes the searchrawtransactions RPC available",
			Widget: "toggle",
			// Hook:        "dropaddrindex",
			OmitEmpty: true,
		},
			true,
		),
		AutoPorts: NewBool(Metadata{
			Option: "autoports",
			Group:  "debug",
			Label:  "Automatic Ports",
			Description:
			"RPC and controller ports are randomized, use with controller for automatic peer discovery",
			Widget: "toggle",
			// Hook: "restart",
			OmitEmpty: true,
		},
			false,
		),
		AutoListen: NewBool(Metadata{
			Option:  "autolisten",
			Aliases: []string{"al"},
			Group:   "node",
			Label:   "Manual Listeners",
			Description:
			"automatically update inbound addresses dynamically according to discovered network interfaces",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		BanDuration: NewDuration(Metadata{
			Option:  "banduration",
			Aliases: []string{"bd"},
			Group:   "debug",
			Label:   "Ban Duration",
			Description:
			"how long a ban of a misbehaving peer lasts",
			Widget: "duration",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			time.Hour*24,
		),
		BanThreshold: NewInt(Metadata{
			Option:  "banthreshold",
			Aliases: []string{"bt"},
			Group:   "debug",
			Label:   "Ban Threshold",
			Description:
			"ban score that triggers a ban (default 100)",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultBanThreshold,
		),
		BlockMaxSize: NewInt(Metadata{
			Option:  "blockmaxsize",
			Aliases: []string{"bmxs"},
			Group:   "mining",
			Label:   "Block Max Size",
			Description:
			"maximum block size in bytes to be used when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			BlockMaxSizeMax,
		),
		BlockMaxWeight: NewInt(Metadata{
			Option:  "blockmaxweight",
			Aliases: []string{"bmxw"},
			Group:   "mining",
			Label:   "Block Max Weight",
			Description:
			"maximum block weight to be used when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			BlockMaxWeightMax,
		),
		BlockMinSize: NewInt(Metadata{
			Option:  "blockminsize",
			Aliases: []string{"bms"},
			Group:   "mining",
			Label:   "Block Min Size",
			Description:
			"minimum block size in bytes to be used when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			BlockMaxSizeMin,
		),
		BlockMinWeight: NewInt(Metadata{
			Option:  "blockminweight",
			Aliases: []string{"bmw"},
			Group:   "mining",
			Label:   "Block Min Weight",
			Description:
			"minimum block weight to be used when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			BlockMaxWeightMin,
		),
		BlockPrioritySize: NewInt(Metadata{
			Option:  "blockprioritysize",
			Aliases: []string{"bps"},
			Group:   "mining",
			Label:   "Block Priority Size",
			Description:
			"size in bytes for high-priority/low-fee transactions when creating a block",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultBlockPrioritySize,
		),
		BlocksOnly: NewBool(Metadata{
			Option:  "blocksonly",
			Aliases: []string{"bo"},
			Group:   "node",
			Label:   "Blocks Only",
			Description:
			"do not accept transactions from remote peers",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		CAFile: NewString(Metadata{
			Option:  "cafile",
			Aliases: []string{"ca"},
			Group:   "tls",
			Label:   "Certificate Authority File",
			Description:
			"certificate authority file for TLS certificate validation",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "ca.cert"),
		),
		ConfigFile: NewString(Metadata{
			Option:  "configfile",
			Aliases: []string{"cf"},
			Label:   "Configuration File",
			Description:
			"location of configuration file, cannot actually be changed",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), PodConfigFilename),
		),
		ConnectPeers: NewStrings(Metadata{
			Option: "connect",
			// Aliases: []string{"cp"},
			Group: "node",
			Label: "Connect Peers",
			Description:
			"connect ONLY to these addresses (disables inbound connections)",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		Controller: NewBool(Metadata{
			Option:  "controller",
			Aliases: []string{"ctrl"},
			Group:   "node",
			Label:   "Enable Controller",
			Description:
			"delivers mining jobs over multicast",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		CPUProfile: NewString(Metadata{
			Option:  "cpuprofile",
			Aliases: []string{"cprof"},
			Group:   "debug",
			Label:   "CPU Profile",
			Description:
			"write cpu profile to this file",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		DarkTheme: NewBool(Metadata{
			Option:  "darktheme",
			Aliases: []string{"dt"},
			Group:   "config",
			Label:   "Dark Theme",
			Description:
			"sets dark theme for GUI",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DataDir: &String{
			value: datadir,
			Metadata: Metadata{
				Option:  "datadir",
				Aliases: []string{"dd", "D"},
				Label:   "Data Directory",
				Description:
				"root folder where application data is stored",
				Type:      "directory",
				Widget:    "string",
				OmitEmpty: true,
			},
			def: appdata.Dir(Name, false),
		},
		DbType: NewString(Metadata{
			Option:  "dbtype",
			Aliases: []string{"dt"},
			Group:   "debug",
			Label:   "Database Type",
			Description:
			"type of database storage engine to use (only one right now, ffldb)",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultDbType,
		),
		DisableBanning: NewBool(Metadata{
			Option:  "nobanning",
			Aliases: []string{"nb"},
			Group:   "debug",
			Label:   "Disable Banning",
			Description:
			"disables banning of misbehaving peers",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DisableCheckpoints: NewBool(Metadata{
			Option:  "nocheckpoints",
			Aliases: []string{"nc"},
			Group:   "debug",
			Label:   "Disable Checkpoints",
			Description:
			"disables all checkpoints",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DisableDNSSeed: NewBool(Metadata{
			Option:  "nodnsseed",
			Aliases: []string{"nds"},
			Group:   "node",
			Label:   "Disable DNS Seed",
			Description:
			"disable seeding of addresses to peers",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DisableListen: NewBool(Metadata{
			Option:  "nolisten",
			Aliases: []string{"nl"},
			Group:   "node",
			Label:   "Disable Listen",
			Description:
			"disables inbound connections for the peer to peer network",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		DisableRPC: NewBool(Metadata{
			Option:  "norpc",
			Aliases: []string{"nr"},
			Group:   "rpc",
			Label:   "Disable RPC",
			Description:
			"disable rpc servers, as well as kopach controller",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		Discovery: NewBool(Metadata{
			Option:  "discover",
			Aliases: []string{"di"},
			Group:   "node",
			Label:   "Disovery",
			Description:
			"enable LAN peer discovery in GUI",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		ExternalIPs: NewStrings(Metadata{
			Option:  "externalip",
			Aliases: []string{"ei"},
			Group:   "node",
			Label:   "External IP Addresses",
			Description:
			"extra addresses to tell peers they can connect to",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		FreeTxRelayLimit: NewFloat(Metadata{
			Option:  "limitfreerelay",
			Aliases: []string{"lfr"},
			Group:   "policy",
			Label:   "Free Tx Relay Limit",
			Description:
			"limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute",
			Widget: "float",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultFreeTxRelayLimit,
		),
		Generate: NewBool(Metadata{
			Option: "generate",
			Group:  "mining",
			Label:  "Generate Blocks",
			Description:
			"turn on Kopach CPU miner",
			Widget: "toggle",
			// Hook:        "generate",
			OmitEmpty: true,
		},
			false,
		),
		GenThreads: NewInt(Metadata{
			Option:  "genthreads",
			Aliases: []string{"G"},
			Group:   "mining",
			Label:   "Generate Threads",
			Description:
			"number of threads to mine with",
			Widget: "integer",
			// Hook:        "genthreads",
			OmitEmpty: true,
		},
			-1,
		),
		Hilite: NewStrings(Metadata{
			Option:  "highlight",
			Aliases: []string{"hl"},
			Group:   "debug",
			Label:   "Hilite",
			Description:
			"list of packages that will print with attention getters",
			Type:   "string",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		LAN: NewBool(Metadata{
			Option: "lan",
			Group:  "debug",
			Label:  "LAN Testnet Mode",
			Description:
			"run without any connection to nodes on the internet (does not apply on mainnet)",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		Language: NewString(Metadata{
			Option:  "language",
			Aliases: []string{"L"},
			Group:   "config",
			Label:   "Language",
			Description:
			"user interface language i18 localization",
			Widget: "string",
			// Hook:        "language",
			OmitEmpty: true,
		},
			"en",
		),
		LimitPass: NewString(Metadata{
			Option:  "limitpass",
			Aliases: []string{"lp"},
			Group:   "rpc",
			Label:   "Limit Password",
			Description:
			"limited user password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			genPassword(),
		),
		LimitUser: NewString(Metadata{
			Option:  "limituser",
			Aliases: []string{"lu"},
			Group:   "rpc",
			Label:   "Limit Username",
			Description:
			"limited user name",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"limit",
		),
		LogDir: NewString(Metadata{
			Option:  "logdir",
			Aliases: []string{"ld"},
			Group:   "config",
			Label:   "Log Directory",
			Description:
			"folder where log files are written",
			Type:   "directory",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			string(datadir.Load().([]byte)),
		),
		LogFilter: NewStrings(Metadata{
			Option:  "logfilter",
			Aliases: []string{"lf"},
			Group:   "debug",
			Label:   "Log Filter",
			Description:
			"list of packages that will not print logs",
			Type:   "string",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		LogLevel: NewString(Metadata{
			Option:  "loglevel",
			Aliases: []string{"L", "ll"},
			Group:   "config",
			Label:   "Log Level",
			Description:
			"maximum log level to output\n(fatal error check warning info debug trace - what is selected includes all items to the left of the one in that list)",
			Widget: "radio",
			Options: []string{"off",
				"fatal",
				"error",
				"info",
				"check",
				"debug",
				"trace",
			},
			// Hook:        "loglevel",
			OmitEmpty: true,
		},
			"info",
		),
		MaxOrphanTxs: NewInt(Metadata{
			Option:  "maxorphantx",
			Aliases: []string{"mt"},
			Group:   "policy",
			Label:   "Max Orphan Txs",
			Description:
			"max number of orphan transactions to keep in memory",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxOrphanTransactions,
		),
		MaxPeers: NewInt(Metadata{
			Option:  "maxpeers",
			Aliases: []string{"mp"},
			Group:   "node",
			Label:   "Max Peers",
			Description:
			"maximum number of peers to hold connections with",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxPeers,
		),
		MulticastPass: NewString(Metadata{
			Option:  "minerpass",
			Aliases: []string{"M"},
			Group:   "config",
			Label:   "Multicast Pass",
			Description:
			"password that encrypts the connection to the mining controller",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"pa55word",
		),
		MiningAddrs: NewStrings(Metadata{
			Option: "miningaddrs",
			// Aliases: []string{"ma"},
			Label: "Mining Addresses",
			Description:
			"addresses to pay block rewards to (not in use)",
			Type:   "base58",
			Widget: "multi",
			// Hook:        "miningaddr",
			OmitEmpty: true,
		},
			[]string{},
		),
		MinRelayTxFee: NewFloat(Metadata{
			Option:  "minrelaytxfee",
			Aliases: []string{"mrf"},
			Group:   "policy",
			Label:   "Min Relay Transaction Fee",
			Description:
			"the minimum transaction fee in DUO/kB to be considered a non-zero fee",
			Widget: "float",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMinRelayTxFee.ToDUO(),
		),
		Network: NewString(Metadata{
			Option:  "network",
			Aliases: []string{"nw"},
			Group:   "node",
			Label:   "Network",
			Description:
			"connect to this network: (mainnet, testnet)",
			Widget: "radio",
			Options: []string{"mainnet",
				"testnet",
				"regtestnet",
				"simnet",
			},
			// Hook:        "restart",
			OmitEmpty: true,
		},
			network,
		),
		NoCFilters: NewBool(Metadata{
			Option:  "nocfilters",
			Aliases: []string{"ncf"},
			Group:   "node",
			Label:   "No CFilters",
			Description:
			"disable committed filtering (CF) support",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		NodeOff: NewBool(Metadata{
			Option:  "nodeoff",
			Aliases: []string{"nn"},
			Group:   "debug",
			Label:   "Node Off",
			Description:
			"turn off the node backend",
			Widget: "toggle",
			// Hook:        "node",
			OmitEmpty: true,
		},
			false,
		),
		NoInitialLoad: NewBool(Metadata{
			Option:  "noinitialload",
			Aliases: []string{"nil"},
			Label:   "No Initial Load",
			Description:
			"do not load a wallet at startup",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		NoPeerBloomFilters: NewBool(Metadata{
			Option:  "nopeerbloomfilters",
			Aliases: []string{"nbf"},
			Group:   "node",
			Label:   "No Peer Bloom Filters",
			Description:
			"disable bloom filtering support",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		NoRelayPriority: NewBool(Metadata{
			Option:  "norelaypriority",
			Aliases: []string{"nrp"},
			Group:   "policy",
			Label:   "No Relay Priority",
			Description:
			"do not require free or low-fee transactions to have high priority for relaying",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		OneTimeTLSKey: NewBool(Metadata{
			Option:  "onetimetlskey",
			Aliases: []string{"otk"},
			Group:   "wallet",
			Label:   "One Time TLS Key",
			Description:
			"generate a new TLS certificate pair at startup, but only write the certificate to disk",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		Onion: NewBool(Metadata{
			Option:  "onion",
			Aliases: []string{"O"},
			Group:   "proxy",
			Label:   "Onion Enabled",
			Description:
			"enable tor proxy",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		OnionProxy: NewString(Metadata{
			Option:  "onionproxy",
			Aliases: []string{"ox"},
			Group:   "proxy",
			Label:   "Onion Proxy Address",
			Description:
			"address of tor proxy you want to connect to",
			Type:   "address",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		OnionProxyPass: NewString(Metadata{
			Option:  "onionproxypass",
			Aliases: []string{"op"},
			Group:   "proxy",
			Label:   "Onion Proxy Password",
			Description:
			"password for tor proxy",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		OnionProxyUser: NewString(Metadata{
			Option:  "onionproxyuser",
			Aliases: []string{"ou"},
			Group:   "proxy",
			Label:   "Onion Proxy Username",
			Description:
			"tor proxy username",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		P2PConnect: NewStrings(Metadata{
			Option:  "p2pconnect",
			Aliases: []string{"p2c"},
			Group:   "node",
			Label:   "P2P Connect",
			Description:
			"list of addresses reachable from connected networks",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		P2PListeners: NewStrings(Metadata{
			Option:  "listen",
			Aliases: []string{"L"},
			Group:   "node",
			Label:   "P2PListeners",
			Description:
			"list of addresses to bind the node listener to",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{net.JoinHostPort("0.0.0.0",
				chaincfg.MainNetParams.DefaultPort,
			),
			},
		),
		Password: NewString(Metadata{
			Option:  "password",
			Aliases: []string{"P"},
			Group:   "rpc",
			Label:   "Password",
			Description:
			"password for client RPC connections",
			Type:   "password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			genPassword(),
		),
		PipeLog: NewBool(Metadata{
			Option:  "pipelog",
			Aliases: []string{"pl"},
			Label:   "Pipe Logger",
			Description:
			"enable pipe based logger IPC",
			Widget: "toggle",
			// Hook:        "",
			OmitEmpty: true,
		},
			false,
		),
		Profile: NewString(Metadata{
			Option: "profile",
			// Aliases: []string{"pr"},
			Group: "debug",
			Label: "Profile",
			Description:
			"http profiling on given port (1024-40000)",
			// Type:        "",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		Proxy: NewString(Metadata{
			Option:  "proxy",
			Aliases: []string{"P"},
			Group:   "proxy",
			Label:   "Proxy",
			Description:
			"address of proxy to connect to for outbound connections",
			Type:   "url",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		ProxyPass: NewString(Metadata{
			Option:  "proxypass",
			Aliases: []string{"pp"},
			Group:   "proxy",
			Label:   "Proxy Pass",
			Description:
			"proxy password, if required",
			Type:   "password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			genPassword(),
		),
		ProxyUser: NewString(Metadata{
			Option:  "proxyuser",
			Aliases: []string{"pu"},
			Group:   "proxy",
			Label:   "ProxyUser",
			Description:
			"proxy username, if required",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"proxyuser",
		),
		RejectNonStd: NewBool(Metadata{
			Option:  "rejectnonstd",
			Aliases: []string{"rn"},
			Group:   "node",
			Label:   "Reject Non Std",
			Description:
			"reject non-standard transactions regardless of the default settings for the active network",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		RelayNonStd: NewBool(Metadata{
			Option:  "relaynonstd",
			Aliases: []string{"R"},
			Group:   "node",
			Label:   "Relay Nonstandard Transactions",
			Description:
			"relay non-standard transactions regardless of the default settings for the active network",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		RPCCert: NewString(Metadata{
			Option:  "rpccert",
			Aliases: []string{"rc"},
			Group:   "rpc",
			Label:   "RPC Cert",
			Description:
			"location of RPC TLS certificate",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "rpc.cert"),
		),
		RPCConnect: NewString(Metadata{
			Option:  "rpcconnect",
			Aliases: []string{"R"},
			Group:   "wallet",
			Label:   "RPC Connect",
			Description:
			"full node RPC for wallet",
			Type:   "address",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			net.JoinHostPort("127.0.0.1", chaincfg.MainNetParams.DefaultPort),
		
		),
		RPCKey: NewString(Metadata{
			Option:  "rpckey",
			Aliases: []string{"rk"},
			Group:   "rpc",
			Label:   "RPC Key",
			Description:
			"location of rpc TLS key",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "rpc.key"),
		),
		RPCListeners: NewStrings(Metadata{
			Option:  "rpclisten",
			Aliases: []string{"rl"},
			Group:   "rpc",
			Label:   "RPC Listeners",
			Description:
			"addresses to listen for RPC connections",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{net.JoinHostPort("127.0.0.1",
				chaincfg.MainNetParams.DefaultPort,
			),
			},
		),
		RPCMaxClients: NewInt(Metadata{
			Option:  "rpcmaxclients",
			Aliases: []string{"rmc"},
			Group:   "rpc",
			Label:   "Maximum RPC Clients",
			Description:
			"maximum number of clients for regular RPC",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxRPCClients,
		),
		RPCMaxConcurrentReqs: NewInt(Metadata{
			Option:  "rpcmaxconcurrentreqs",
			Aliases: []string{"rr"},
			Group:   "rpc",
			Label:   "Maximum RPC Concurrent Reqs",
			Description:
			"maximum number of requests to process concurrently",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxRPCConcurrentReqs,
		),
		RPCMaxWebsockets: NewInt(Metadata{
			Option:  "rpcmaxwebsockets",
			Aliases: []string{"rw"},
			Group:   "rpc",
			Label:   "Maximum RPC Websockets",
			Description:
			"maximum number of websocket clients to allow",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultMaxRPCWebsockets,
		),
		RPCQuirks: NewBool(Metadata{
			Option:  "rpcquirks",
			Aliases: []string{"rq"},
			Group:   "rpc",
			Label:   "RPC Quirks",
			Description:
			"enable bugs that replicate bitcoin core RPC's JSON",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		RunAsService: NewBool(Metadata{
			Option:  "runasservice",
			Aliases: []string{"raas"},
			Label:   "Run As Service",
			Description:
			"shuts down on lock timeout",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		ServerPass: NewString(Metadata{
			Option:  "serverpass",
			Aliases: []string{"sp"},
			Group:   "rpc",
			Label:   "Server Pass",
			Description:
			"password for server connections",
			Type:   "password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			genPassword(),
		),
		ServerTLS: NewBool(Metadata{
			Option:  "servertls",
			Aliases: []string{"st"},
			Group:   "wallet",
			Label:   "Server TLS",
			Description:
			"enable TLS for the wallet connection to node RPC server",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		ServerUser: NewString(Metadata{
			Option:  "serveruser",
			Aliases: []string{"su"},
			Group:   "rpc",
			Label:   "Server User",
			Description:
			"username for chain server connections",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"client",
		),
		SigCacheMaxSize: NewInt(Metadata{
			Option:  "sigcachemaxsize",
			Aliases: []string{"scm"},
			Group:   "node",
			Label:   "Signature Cache Max Size",
			Description:
			"the maximum number of entries in the signature verification cache",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultSigCacheMaxSize,
		),
		Solo: NewBool(Metadata{
			Option: "solo",
			Group:  "mining",
			Label:  "Solo Generate",
			Description:
			"mine even if not connected to a network",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		TLS: NewBool(Metadata{
			Option:  "clienttls",
			Aliases: []string{"ct"},
			Group:   "tls",
			Label:   "TLS",
			Description:
			"enable TLS for RPC client connections",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		TLSSkipVerify: NewBool(Metadata{
			Option:  "tlsskipverify",
			Aliases: []string{"sv"},
			Group:   "tls",
			Label:   "TLS Skip Verify",
			Description:
			"skip TLS certificate verification (ignore CA errors)",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			false,
		),
		TorIsolation: NewBool(Metadata{
			Option:  "torisolation",
			Aliases: []string{"T"},
			Group:   "proxy",
			Label:   "Tor Isolation",
			Description:
			"makes a separate proxy connection for each connection",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		TrickleInterval: NewDuration(Metadata{
			Option:  "trickleinterval",
			Aliases: []string{"tt"},
			Group:   "policy",
			Label:   "Trickle Interval",
			Description:
			"minimum time between attempts to send new inventory to a connected peer",
			Widget: "duration",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultTrickleInterval,
		),
		TxIndex: NewBool(Metadata{
			Option:  "txindex",
			Aliases: []string{"ti"},
			Group:   "node",
			Label:   "Tx Index",
			Description:
			"maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC",
			Widget: "toggle",
			// Hook:        "droptxindex",
			OmitEmpty: true,
		},
			true,
		),
		UPNP: NewBool(Metadata{
			Option: "upnp",
			Group:  "node",
			Label:  "UPNP",
			Description:
			"enable UPNP for NAT traversal",
			Widget: "toggle",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			true,
		),
		UserAgentComments: NewStrings(Metadata{
			Option:  "uacomment",
			Aliases: []string{"ua"},
			Group:   "policy",
			Label:   "User Agent Comments",
			Description:
			"comment to add to the user agent -- See BIP 14 for more information",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
		Username: NewString(Metadata{
			Option:  "username",
			Aliases: []string{"U"},
			Group:   "rpc",
			Label:   "Username",
			Description:
			"password for client RPC connections",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"username",
		),
		UUID: &Int{Metadata: Metadata{
			Option: "uuid",
			Label:  "UUID",
			Description:
			"instance unique id (64bit random value)",
			Widget:    "string",
			OmitEmpty: true,
		},
			value: uberatomic.NewInt64(rand.Int63()),
		},
		Wallet: NewBool(Metadata{
			Option:  "walletconnect",
			Aliases: []string{"wc"},
			Group:   "debug",
			Label:   "Connect to Wallet",
			Description:
			"set ctl to connect to wallet instead of chain server",
			Widget:    "toggle",
			OmitEmpty: true,
		},
			false,
		),
		WalletFile: NewString(Metadata{
			Option:  "walletfile",
			Aliases: []string{"wf"},
			Group:   "config",
			Label:   "Wallet File",
			Description:
			"wallet database file",
			Type:   "path",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			filepath.Join(string(datadir.Load().([]byte)), "mainnet", DbName),
		),
		WalletOff: NewBool(Metadata{
			Option:  "walletoff",
			Aliases: []string{"wo"},
			Group:   "debug",
			Label:   "Wallet Off",
			Description:
			"turn off the wallet backend",
			Widget: "toggle",
			// Hook:        "wallet",
			OmitEmpty: true,
		},
			false,
		),
		WalletPass: NewString(Metadata{
			Option:  "walletpass",
			Aliases: []string{"wp"},
			Label:   "Wallet Pass",
			Description:
			"password encrypting public data in wallet - hash is stored so give on command line",
			Type:   "password",
			Widget: "password",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			"",
		),
		WalletRPCListeners: NewStrings(Metadata{
			Option:  "walletrpclisten",
			Aliases: []string{"wr"},
			Group:   "wallet",
			Label:   "Wallet RPC Listeners",
			Description:
			"addresses for wallet RPC server to listen on",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{net.JoinHostPort("0.0.0.0",
				chaincfg.MainNetParams.WalletRPCServerPort,
			),
			},
		),
		WalletRPCMaxClients: NewInt(Metadata{
			Option:  "walletrpcmaxclients",
			Aliases: []string{"wmc"},
			Group:   "wallet",
			Label:   "Legacy RPC Max Clients",
			Description:
			"maximum number of RPC clients allowed for wallet RPC",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultRPCMaxClients,
		),
		WalletRPCMaxWebsockets: NewInt(Metadata{
			Option:  "walletrpcmaxwebsockets",
			Aliases: []string{"wrm"},
			Group:   "wallet",
			Label:   "Legacy RPC Max Websockets",
			Description:
			"maximum number of websocket clients allowed for wallet RPC",
			Widget: "integer",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			DefaultRPCMaxWebsockets,
		),
		WalletServer: NewString(Metadata{
			Option:  "walletserver",
			Aliases: []string{"ws"},
			Group:   "wallet",
			Label:   "Wallet Server",
			Description:
			"node address to connect wallet server to",
			Type:   "address",
			Widget: "string",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			net.JoinHostPort("127.0.0.1",
				chaincfg.MainNetParams.WalletRPCServerPort,
			),
		),
		Whitelists: NewStrings(Metadata{
			Option:  "whitelists",
			Aliases: []string{"wl"},
			Group:   "debug",
			Label:   "Whitelists",
			Description:
			"peers that you don't want to ever ban",
			Type:   "address",
			Widget: "multi",
			// Hook:        "restart",
			OmitEmpty: true,
		},
			[]string{},
		),
	}
	// check sanity of configuration
	// I.S(c.getAllOptionStrings())
	return
}
