package common

// TaprootActivation
// todo might be inapplicable due to transactions that have taproot prevouts from before the activation
//
//	is relevant for the height-to-hash lookup in the db
const TaprootActivation uint32 = 709632
const ConfigFileName string = "blindbit.toml"
const DefaultBaseDirectory = "~/.blindbit-oracle"

var TweaksOnly bool
var TweakIndexFullNoDust bool
var TweakIndexFullIncludingDust bool
var TweaksCutThroughWithDust bool

var (
	RpcEndpoint = "http://127.0.0.1:8332" // default local node
	CookiePath  = ""
	RpcUser     = ""
	RpcPass     = ""

	BaseDirectory = ""
	DBPath        = ""
	LogsPath      = ""

	Host = "127.0.0.1:8000"
)

type chain int

const (
	Unknown chain = iota
	Mainnet
	Signet
	Regtest
	Testnet3
)

// control vars
var (
	SyncStartHeight uint32 = 833_000 // random block where BIP352 was not active yet. todo change to actual number
	// MinHeightToProcess No block below this number will be processed
	// todo is this actually needed
	//MinHeightToProcess uint32 = 833_000

	Chain = Unknown

	// SyncHeadersMaxPerCall how many headers will maximally be requested in one batched RPC call
	SyncHeadersMaxPerCall uint32 = 10_000
	// MaxParallelRequests sets how many RPC calls will be made in parallel to the Node
	MaxParallelRequests uint16 = 2
	// MaxParallelTweakComputations number of parallel processes which will be spawned in order to compute the tweaks for a given block
	MaxParallelTweakComputations = 2
	// PruneFrequency every x blocks the data will be checked and pruned
	// possible routines: -remove utxos for 100% spent transaction
	PruneFrequency = 72
)

// one has to call SetDirectories otherwise common.DBPath will be empty
var (
	DBPathHeaders        string
	DBPathHeadersInv     string
	DBPathFilters        string
	DBPathTweaks         string
	DBPathTweakIndex     string
	DBPathTweakIndexDust string
	DBPathUTXOs          string
)

// NumsH = 0x50929b74c1a04954b78b4b6035e97a5e078a5a0f28ec96d547bfee9ace803ac0
var NumsH = []byte{80, 146, 155, 116, 193, 160, 73, 84, 183, 139, 75, 96, 53, 233, 122, 94, 7, 138, 90, 15, 40, 236, 150, 213, 71, 191, 238, 154, 206, 128, 58, 192}

func SetDirectories() {
	BaseDirectory = ResolvePath(BaseDirectory)

	DBPath = BaseDirectory + "/data"
	LogsPath = BaseDirectory + "/logs"

	DBPathHeaders = DBPath + "/headers"
	DBPathHeadersInv = DBPath + "/headers-inv"
	DBPathFilters = DBPath + "/filters"
	DBPathTweaks = DBPath + "/tweaks"
	DBPathTweakIndex = DBPath + "/tweak-index"
	DBPathTweakIndexDust = DBPath + "/tweak-index-dust"
	DBPathUTXOs = DBPath + "/utxos"
}

func HeaderMustSyncHeight() uint32 {
	switch Chain {
	case Mainnet:
		return 500_000
	case Signet:
		return 1
	case Regtest:
		return 1
	case Testnet3:
		return 1
	case Unknown:
		panic("chain not defined")
	default:
		return 1
	}
}
