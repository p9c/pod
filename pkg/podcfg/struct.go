package podcfg

// Config defines the configuration items used by pod along with the various components included in the suite
type Config struct {
	// ShowAll is a flag to make the json encoder explicitly define all fields and not just the ones different to the
	// defaults
	ShowAll bool
	// Map is the same data but addressible using its name as found inside the various configuration types, the key is
	// the same as the .Name field field in the various data types
	Map            map[string]interface{}
	Commands       Commands
	RunningCommand *Command
	// These are just the definitions, the things put in them are more useful than doc comments
	// todo: these will be generated
	AddCheckpoints         *Strings
	AddPeers               *Strings
	AddrIndex              *Bool
	AutoListen             *Bool
	AutoPorts              *Bool
	BanDuration            *Duration
	BanThreshold           *Int
	BlockMaxSize           *Int
	BlockMaxWeight         *Int
	BlockMinSize           *Int
	BlockMinWeight         *Int
	BlockPrioritySize      *Int
	BlocksOnly             *Bool
	CAFile                 *String
	ConfigFile             *String
	ConnectPeers           *Strings
	Controller             *Bool
	CPUProfile             *String
	DarkTheme              *Bool
	DataDir                *String
	DbType                 *String
	DisableBanning         *Bool
	DisableCheckpoints     *Bool
	DisableDNSSeed         *Bool
	DisableListen          *Bool
	DisableRPC             *Bool
	Discovery              *Bool
	ExternalIPs            *Strings
	FreeTxRelayLimit       *Float
	Generate               *Bool
	GenThreads             *Int
	Hilite                 *Strings
	LAN                    *Bool
	Language               *String
	LimitPass              *String
	LimitUser              *String
	LogDir                 *String
	LogFilter              *Strings
	LogLevel               *String
	MaxOrphanTxs           *Int
	MaxPeers               *Int
	MulticastPass          *String
	MiningAddrs            *Strings
	MinRelayTxFee          *Float
	Network                *String
	NoCFilters             *Bool
	NodeOff                *Bool
	NoInitialLoad          *Bool
	NoPeerBloomFilters     *Bool
	NoRelayPriority        *Bool
	OneTimeTLSKey          *Bool
	Onion                  *Bool
	OnionProxy             *String
	OnionProxyPass         *String
	OnionProxyUser         *String
	P2PConnect             *Strings
	P2PListeners           *Strings
	Password               *String
	PipeLog                *Bool
	Profile                *String
	Proxy                  *String
	ProxyPass              *String
	ProxyUser              *String
	RejectNonStd           *Bool
	RelayNonStd            *Bool
	RPCCert                *String
	RPCConnect             *String
	RPCKey                 *String
	RPCListeners           *Strings
	RPCMaxClients          *Int
	RPCMaxConcurrentReqs   *Int
	RPCMaxWebsockets       *Int
	RPCQuirks              *Bool
	RunAsService           *Bool
	Save                   *Bool
	ServerPass             *String
	ServerTLS              *Bool
	ServerUser             *String
	SigCacheMaxSize        *Int
	Solo                   *Bool
	TLS                    *Bool
	TLSSkipVerify          *Bool
	TorIsolation           *Bool
	TrickleInterval        *Duration
	TxIndex                *Bool
	UPNP                   *Bool
	UserAgentComments      *Strings
	Username               *String
	UUID                   *Int
	Wallet                 *Bool
	WalletFile             *String
	WalletOff              *Bool
	WalletPass             *String
	WalletRPCListeners     *Strings
	WalletRPCMaxClients    *Int
	WalletRPCMaxWebsockets *Int
	WalletServer           *String
	Whitelists             *Strings
}
