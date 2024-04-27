package solanacore

import (
	"context"
	"github.com/dde-group/blockchain-core/solanacore/raydium"
	"github.com/dde-group/blockchain-core/utils/logutils"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"go.uber.org/zap"
	"testing"
	"time"
)

func Test_BlockSub(t *testing.T) {
	logutils.InitLogger(logutils.DefaultZapConfig)
	rpcEndpoint := "https://wiser-wild-vineyard.solana-mainnet.quiknode.pro/efe8fb2f8241640535b662658afb570bce5cf227/"
	wsEndpoint := "wss://wiser-wild-vineyard.solana-mainnet.quiknode.pro/efe8fb2f8241640535b662658afb570bce5cf227/"

	core, err := NewCore(rpcEndpoint, wsEndpoint)
	if nil != err {
		logutils.Panic("core init failed", zap.Error(err))
	}

	tableCache := make(map[string]solana.PublicKeySlice)

	now := time.Now()
	_, err = core.WsAgent.BlockSubscribeMentions(
		solana.MustPublicKeyFromBase58(raydium.LiquidityPoolProgramV4),
		rpc.CommitmentConfirmed,
		func(result *ws.BlockResult) {
			slot := result.Context.Slot
			blockTime := result.Value.Block.BlockTime
			blockHeight := uint64(0)
			if nil != result.Value.Block.BlockHeight {
				blockHeight = *result.Value.Block.BlockHeight
			}

			success := 0
			//var pcTokenAmount, coinTokenAmount uint64
			//detailList := make([]*SwapPoolTransactionDetail, 0, len(result.Value.Block.Transactions))
			for txIdx, tx := range result.Value.Block.Transactions {
				table := make(map[solana.PublicKey]solana.PublicKeySlice)

				if nil != tx.Meta.Err {
					continue
				}
				rawTx := solana.Transaction{}

				err = bin.NewBinDecoder(tx.Transaction.GetBinary()).Decode(&rawTx)
				if nil != err {
					continue
				}

				slice := rawTx.Message.GetAddressTableLookups().GetTableIDs()

				var state *addresslookuptable.AddressLookupTableState
				for _, account := range slice {
					address := account.String()
					address = address
					//logutils.Info("addresslookuptable", zap.String("address", address))
					//tableCache[address] = struct{}{}

					keys, ok := tableCache[address]
					if ok {
						table[account] = keys
					} else {
						state, err = addresslookuptable.GetAddressLookupTable(context.TODO(),
							core.client,
							account,
						)
						if nil != err {
							logutils.Error("addresslookuptable get error", zap.Error(err))
							continue
						} else {
							logutils.Info("addresslookuptable get success", zap.String("address", address))
						}
						table[account] = state.Addresses
						tableCache[address] = state.Addresses
					}

				}
				err = rawTx.Message.SetAddressTables(table)
				err = rawTx.Message.ResolveLookups()
				if nil != err {
					logutils.Error("resolveLookups error", zap.Error(err))
				}
				swapResult := ParseSwapRayV4ResultFromRawTx(slot, txIdx, &rawTx)

				swapResult = swapResult
				success++
			}
			since := time.Since(now)
			var tm time.Time
			if nil != blockTime {
				tm = time.Unix(int64(*blockTime), 0)
			}
			logutils.Info("block binary", zap.Int("table", len(tableCache)), zap.Uint64("slot", slot),
				zap.Uint64("height", blockHeight),
				zap.Time("blockTime", tm), zap.Int("total", len(result.Value.Block.Transactions)),
				zap.Int("success", success), zap.Duration("duration", since),
			)
			now = time.Now()
		})
	if nil != err {
		logutils.Panic("ws sub err", zap.Error(err))
	} else {
		logutils.Info("ws sub success")
	}
	select {}
}

func Test_ParsedBlockSub(t *testing.T) {
	logutils.InitLogger(logutils.DefaultZapConfig)
	rpcEndpoint := rpc.MainNetBeta_RPC
	wsEndpoint := "wss://wiser-wild-vineyard.solana-mainnet.quiknode.pro/efe8fb2f8241640535b662658afb570bce5cf227/"

	core, err := NewCore(rpcEndpoint, wsEndpoint)
	if nil != err {
		logutils.Panic("core init failed", zap.Error(err))
	}

	now := time.Now()
	_, err = core.WsAgent.ParsedBlockSubscribeMentions(
		solana.MustPublicKeyFromBase58(raydium.LiquidityPoolProgramV4),
		rpc.CommitmentConfirmed,
		func(result *ws.ParsedBlockResult) {
			slot := result.Context.Slot
			blockTime := result.Value.Block.BlockTime
			blockHeight := uint64(0)
			if nil != result.Value.Block.BlockHeight {
				blockHeight = *result.Value.Block.BlockHeight
			}
			//var pcTokenAmount, coinTokenAmount uint64
			//detailList := make([]*SwapPoolTransactionDetail, 0, len(result.Value.Block.Transactions))
			success := 0
			dealNow := time.Now()
			for txIdx, tx := range result.Value.Block.Transactions {
				//txHash = result.Value.Block.Signatures[idx]
				//txHash = txHash
				if nil != tx.Meta.Err {
					continue
				}
				swapResult := ParseSwapRaydiumV4ResultFromTx(slot, txIdx, &tx)

				swapResult = swapResult
				success++
			}

			since := time.Since(now)
			var tm time.Time
			if nil != blockTime {
				tm = time.Unix(int64(*blockTime), 0)
			}
			logutils.Info("block jsonRpc", zap.Uint64("slot", slot),
				zap.Uint64("height", blockHeight),
				zap.Time("blockTime", tm), zap.Int("total", len(result.Value.Block.Transactions)),
				zap.Int("success", success), zap.Duration("duration", since),
				zap.Duration("deal", time.Since(dealNow)))
			now = time.Now()

		})
	if nil != err {
		logutils.Panic("ws sub err", zap.Error(err))
	} else {
		logutils.Info("ws sub success", zap.Duration("duration", time.Since(now)))
	}
	select {}
}
