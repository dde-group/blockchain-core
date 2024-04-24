package blockchain

type TokenSupply struct {
	TokenAddress   string //地址
	Supply         uint64 //整型数量
	Decimals       uint   //精度
	SupplyDecimals string //小数数量
}
