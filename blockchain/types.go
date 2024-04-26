package blockchain

import "github.com/shopspring/decimal"

type TokenSupply struct {
	TokenAddress   string //地址
	Supply         uint64 //整型数量
	Decimals       uint   //精度
	SupplyDecimals string //小数数量
}

type TokenBalance struct {
	Account        string //即ATA 地址
	Mint           string
	Owner          string
	Amount         uint64          //整型数量
	AmountDecimals decimal.Decimal //小数数量
}
