package solanacore

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type ATAAccountInfo struct {
	Account solana.PublicKey //账户地址
	Mint    solana.PublicKey //token 地址
	Owner   solana.PublicKey //owner
}

type TokenInfo struct {
	Mint     solana.PublicKey //地址
	Decimals uint8            //精度
}

type BlockParsedTransaction struct {
	Transaction     *rpc.ParsedTransaction
	TransactionMeta *rpc.TransactionMeta
}
