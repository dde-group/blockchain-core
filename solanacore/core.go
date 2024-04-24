package solanacore

import (
	"github.com/gagliardetto/solana-go/rpc"
	"sync"
)

type Core struct {
	client *rpc.Client

	decimalsCache sync.Map
	//获取最新区块
	latestBlockHashSchedule *getLatestBlockHashHandler
	//获取租金
	minimumBalanceForRentExemptionHandler *getMinimumBalanceForRentExemptionHandler
}

func NewCore(client *rpc.Client) *Core {
	ret := &Core{
		client:                                client,
		latestBlockHashSchedule:               newGetLatestBlockHashHandler(client),
		minimumBalanceForRentExemptionHandler: newGetMinimumBalanceForRentExemptionSchedule(client),
	}

	return ret
}
