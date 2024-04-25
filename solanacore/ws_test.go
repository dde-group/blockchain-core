package solanacore

import (
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"gitlab.xbit.trade/blockchain/blockchain-core/solanacore/raydium"
	"gitlab.xbit.trade/blockchain/blockchain-core/utils/logutils"
	"go.uber.org/zap"
	"testing"
)

func Test_BlockSub(t *testing.T) {
	logutils.InitLogger(logutils.DefaultZapConfig)
	rpcEndpoint := rpc.MainNetBeta_RPC
	wsEndpoint := "wss://wiser-wild-vineyard.solana-mainnet.quiknode.pro/efe8fb2f8241640535b662658afb570bce5cf227/"

	core, err := NewCore(rpcEndpoint, wsEndpoint)
	if nil != err {
		logutils.Panic("core init failed", zap.Error(err))
	}

	_, err = core.WsAgent.BlockSubscribeMentions(
		solana.MustPublicKeyFromBase58(raydium.LiquidityPoolProgramV4),
		rpc.CommitmentFinalized,
		func(result *ws.BlockResult) {
			//slot := result.Context.Slot

			txHash := ""
			count := 0
			//var pcTokenAmount, coinTokenAmount uint64
			//detailList := make([]*SwapPoolTransactionDetail, 0, len(result.Value.Block.Transactions))
			for _, tx := range result.Value.Block.Transactions {
				//txHash = result.Value.Block.Signatures[idx]
				//txHash = txHash
				if nil != tx.Meta.Err {
					continue
				}

				rawTx := solana.Transaction{}
				err = bin.NewBinDecoder(tx.Transaction.GetBinary()).Decode(&rawTx)
				if nil != err {
					continue
				}

				accountKeys, err := rawTx.Message.GetAllKeys()
				if nil != err {
					continue
				}
				accountKeys = accountKeys
				hash := rawTx.Signatures[0]
				txHash = hash.String()
				txHash = txHash
				count++
			}
			if count < 0 {
				return
			}
		})
	if nil != err {
		logutils.Panic("ws sub err", zap.Error(err))
	}
	select {}
}
