package blockchain

const (
	ChainSolana = "solana"
	ChainEth    = "eth"
	ChainBsc    = "bsc"
)

const (
	ChainIdSolana = 101
)

const (
	MainTokenUnknown = "Main Token"
	MainTokenSolana  = "SOL"
	MainTokenEth     = "ETH"
	MainTokenBsc     = "BNB"
)

type BlockChainSetting struct {
	Name             string
	ChainId          int
	MainToken        string
	MainTokenAddress string
}

var (
	CacheBlockChainSetting = map[string]BlockChainSetting{
		ChainSolana: {
			Name:      ChainSolana,
			ChainId:   ChainIdSolana,
			MainToken: MainTokenSolana,
		},
		ChainEth: {
			Name: ChainEth,
		},
		ChainBsc: {
			Name: ChainBsc,
		},
	}
)
