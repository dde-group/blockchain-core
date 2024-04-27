package solanacore

import (
	"github.com/dde-group/blockchain-core/solanacore/raydium"
	"github.com/dde-group/blockchain-core/utils/logutils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/zap"
)

func ParseSwapRayV4ResultFromRawTx(slot uint64, txIdx int, rawTx *solana.Transaction) []*SwapInstructionResult {
	txHash := rawTx.Signatures[0]
	hash := txHash.String()
	hash = hash

	interact := solana.MustPublicKeyFromBase58(raydium.LiquidityPoolProgramV4)

	raydiumDetailCache := make(map[uint64]*raydium.InstructionPoolAccountsDetail)
	ret := make([]*SwapInstructionResult, 0, 2)
	//解析直接调用的inst
	for index, inst := range rawTx.Message.Instructions {
		programId, err := rawTx.Message.Account(inst.ProgramIDIndex)
		if err != nil {

		}

		if !programId.Equals(interact) {
			continue
		}
		detail := raydium.GetInstructionPoolAccountDetail(inst.Accounts)
		detail.Index = uint64(index)
		raydiumDetailCache[detail.Index] = detail
		//raydiumDetailList = append(raydiumDetailList, detail)
		break
	}

	return nil
}

func ParseSwapRaydiumV4ResultFromTx(slot uint64, txIdx int, tx *rpc.TransactionParsedWithMeta) []*SwapInstructionResult {

	rawTx := tx.Transaction
	txHash := rawTx.Signatures[0]
	hash := txHash.String()

	interact := solana.MustPublicKeyFromBase58(raydium.LiquidityPoolProgramV4)

	raydiumDetailCache := make(map[uint64]*raydium.InstructionPoolAccountsDetail)
	ret := make([]*SwapInstructionResult, 0, 2)
	//解析直接调用的inst
	for index, inst := range rawTx.Message.Instructions {
		programId := inst.ProgramId
		if !programId.Equals(interact) {
			continue
		}
		detail := raydium.GetInstructionPoolAccountDetail(inst.Accounts)
		detail.Index = uint64(index)
		raydiumDetailCache[detail.Index] = detail
		//raydiumDetailList = append(raydiumDetailList, detail)
		break
	}

	//swapResultCache := make(map[string]*solanacore.SwapInstructionResult)

	//解析封装在合约立的inst
	for _, inner := range tx.Meta.InnerInstructions {
		instIndex := inner.Index
		detail, ok := raydiumDetailCache[instIndex]
		//此处的inner Inst 就是单独处理radyium inst
		if ok {
			swapResult, err := ParseInstruction2SwapResult(inner.Instructions)
			if nil != err {
				logutils.Warn("ParseSwapResultFromTx inst ParseInstruction2SwapResult not swap inst", zap.Error(err),
					zap.Uint64("instIndex", instIndex),
					zap.String("txHash", hash),
					zap.Any("inst", inner.Instructions))
				continue
			}
			swapResult.AmmId = detail.AmmId
			swapResult.AmmCoinAccount = detail.AmmCoinAccount
			swapResult.AmmPcAccount = detail.AmmPcAccount

			swapResult.Slot = slot
			swapResult.Interact = interact
			swapResult.Hash = txHash
			swapResult.Index = uint64(txIdx)
			swapResult.InstIndex = instIndex
			swapResult.SubIndex = 0
			ret = append(ret, swapResult)
			continue
		}
		//后面的inner inst 有可能是raydium 作为其他合约的调用inst 他的两个inner inst 跟在raydium 后面
		parsedIdx := 0
		for innerIdx, inst := range inner.Instructions {
			if innerIdx < parsedIdx {
				continue
			}
			programId := inst.ProgramId
			if !programId.Equals(interact) {
				continue
			}
			//处理radium inner inst59*
			detail = raydium.GetInstructionPoolAccountDetail(inst.Accounts)
			//后面连续两个就是swap 的两个transfer
			swapResult, err := ParseInstruction2SwapResult(inner.Instructions[innerIdx+1 : innerIdx+3])
			if nil != err {
				logutils.Warn("ParseSwapResultFromTx inner ParseInstruction2SwapResult not swap inst", zap.Error(err),
					zap.Uint64("instIndex", instIndex),
					zap.Int("innerIdx", innerIdx),
					zap.String("txHash", hash),
					zap.Any("inst", inner.Instructions))
				continue
			}
			swapResult.AmmId = detail.AmmId
			swapResult.AmmCoinAccount = detail.AmmCoinAccount
			swapResult.AmmPcAccount = detail.AmmPcAccount

			swapResult.Slot = slot
			swapResult.Interact = interact
			swapResult.Hash = txHash
			swapResult.Index = uint64(txIdx)
			swapResult.InstIndex = instIndex
			swapResult.SubIndex = uint64(innerIdx)
			ret = append(ret, swapResult)
			//后续循环跳过这两个inst
			innerIdx += 2
		}
	}

	return ret
}
