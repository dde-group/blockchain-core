package solanacore

import (
	"fmt"
	"github.com/gagliardetto/solana-go/rpc"
	"sync"
)

type Core struct {
	client  *rpc.Client
	WsAgent *WebsocketAgent

	decimalsCache sync.Map
	//获取最新区块
	latestBlockHashSchedule *getLatestBlockHashHandler
	//获取租金
	minimumBalanceForRentExemptionHandler *getMinimumBalanceForRentExemptionHandler
}

func NewCore(rpcEndpoint, wsEndpoint string) (*Core, error) {
	client := rpc.New(rpcEndpoint)

	wsAgent := NewWebsocketAgent(wsEndpoint)
	err := wsAgent.Connect()
	if nil != err {
		return nil, fmt.Errorf("ws connect err: %s", err.Error())
	}

	ret := &Core{
		client:                                client,
		WsAgent:                               wsAgent,
		latestBlockHashSchedule:               newGetLatestBlockHashHandler(client),
		minimumBalanceForRentExemptionHandler: newGetMinimumBalanceForRentExemptionSchedule(client),
	}

	return ret, nil
}
