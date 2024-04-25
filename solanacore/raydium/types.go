package raydium

type PoolInfo struct {
	Owner    string
	OpenTime int64
	Status   string
	//RayProgramID string
	AmmInfo
	SerumMarketInfo
	OpenOrder
}

type AmmStateData struct {
	NeedTakePnlCoin     uint64 `json:"NeedTakePnlCoin"`
	NeedTakePnlPc       uint64 `json:"NeedTakePnlPc"`
	TotalPnlPc          uint64 `json:"TotalPnlPc"`
	TotalPnlCoin        uint64 `json:"TotalPnlCoin"`
	PoolOpenTime        int64  `json:"PoolOpenTime"`
	PunishPcAmount      uint64 `json:"PunishPcAmount"`
	PunishCoinAmount    uint64 `json:"PunishCoinAmount"`
	OrderbookToInitTime int64  `json:"OrderbookToInitTime"`
	SwapCoinInAmount    string `json:"SwapCoinInAmount"`
	SwapPcOutAmount     string `json:"SwapPcOutAmount"`
	SwapTakePcFee       uint64 `json:"SwapTakePcFee"`
	SwapPcInAmount      string `json:"SwapPcInAmount"`
	SwapCoinOutAmount   string `json:"SwapCoinOutAmount"`
	SwapTakeCoinFee     uint64 `json:"SwapTakeCoinFee"`
}

type AmmFees struct {
	MinSeparateNumerator   int `json:"MinSeparateNumerator"`
	MinSeparateDenominator int `json:"MinSeparateDenominator"`
	TradeFeeNumerator      int `json:"TradeFeeNumerator"`
	TradeFeeDenominator    int `json:"TradeFeeDenominator"`
	PnlNumerator           int `json:"PnlNumerator"`
	PnlDenominator         int `json:"PnlDenominator"`
	SwapFeeNumerator       int `json:"SwapFeeNumerator"`
	SwapFeeDenominator     int `json:"SwapFeeDenominator"`
}

type OpenOrder struct {
	NativeBaseTokenFree   uint64
	NativeBaseTokenTotal  uint64
	NativeQuoteTokenFree  uint64
	NativeQuoteTokenTotal uint64
	FreeSlotBits          string
	IsBidBits             string
}

type AmmInfo struct {
	AmmId               string
	AmmAuthority        string
	AmmOpenOrders       string
	AmmTargetOrders     string
	AmmCoinTokenAccount string
	AmmPcTokenAccount   string
	//AmmQuantities   string
	Fees      AmmFees
	StateData AmmStateData
}

type SerumMarketInfo struct {
	SerumProgramId        string
	SerumMarket           string
	SerumBids             string
	SerumAsks             string
	SerumEventQueue       string
	SerumCoinVaultAccount string
	SerumPcVaultAccount   string
	SerumVaultSigner      string
}
