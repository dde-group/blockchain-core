package solanacore

import (
	"context"
	"github.com/gagliardetto/solana-go"
	"testing"
)

func Test_Token(t *testing.T) {
	endpoint := "https://ancient-spring-dream.solana-mainnet.quiknode.pro/a3c521971fd30abcc981696f5b03872d4a5fe90e/"

	core, _ := NewCore(endpoint, "")

	addr := "7caQgZMcyDkuwTCFTSZKgofgcuWvhWLV63w311eNZN1c"
	account, _ := solana.PublicKeyFromBase58(addr)
	_, _ = core.GetAccountToken(context.TODO(), account)
}
