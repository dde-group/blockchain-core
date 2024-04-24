package solanacore

import (
	"context"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"testing"
)

func Test_Token(t *testing.T) {
	endpoint := "https://ancient-spring-dream.solana-mainnet.quiknode.pro/a3c521971fd30abcc981696f5b03872d4a5fe90e/"
	client := rpc.New(endpoint)
	core := NewCore(client)

	addr := "7caQgZMcyDkuwTCFTSZKgofgcuWvhWLV63w311eNZN1c"
	account, _ := solana.PublicKeyFromBase58(addr)
	_, _ = core.GetAccountToken(context.TODO(), account)
}
