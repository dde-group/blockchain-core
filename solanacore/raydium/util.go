package raydium

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math/big"
)

//ComputeAmountOut
/* @Description: 计算换出数量
 * @param amountIn
 * @param pcAmount
 * @param coinAmount
 * @param direction
 * @param pool
 * @return uint64
 */
func ComputeAmountOut(
	amountIn uint64,
	pcAmount uint64,
	coinAmount uint64,
	direction SwapDirection,
	pool *PoolInfo) (uint64, error) {

	//计算手续费
	var fee = big.NewInt(0)
	var amountInDeductFee = big.NewInt(0)
	amountInBN := big.NewInt(0).SetUint64(amountIn)
	fee.Mul(amountInBN, big.NewInt(int64(pool.AmmInfo.Fees.SwapFeeNumerator)))

	feeDecimal := decimal.NewFromBigInt(fee, 0)
	denominator := big.NewInt(int64(pool.AmmInfo.Fees.SwapFeeDenominator))
	fee = feeDecimal.Div(decimal.NewFromBigInt(denominator, 0)).Ceil().BigInt()

	amountInDeductFee.Sub(amountInBN, fee)

	//计算token0 token1 数量
	pcWithoutTakePnl, coinWithoutTakePnl, err := calcTotalWithoutTakePnlNoOrderBook(pcAmount, coinAmount, &pool.AmmInfo)
	if nil != err {
		return 0, err
	}

	amountOut := swapTokenAmountBaseIn(amountInDeductFee.Uint64(), pcWithoutTakePnl, coinWithoutTakePnl, direction)

	return amountOut, nil
}

//swapTokenAmountBaseIn
/* @Description: 计算最终换出数量
 * @param amountInDeductFee 扣出手续费的输入数量
 * @param pcWithoutTakePnl 扣除收益的pcAmount
 * @param coinWithoutTakePnl
 * @param direction
 * @return uint64
 */
func swapTokenAmountBaseIn(amountInDeductFee uint64, pcWithoutTakePnl uint64, coinWithoutTakePnl uint64, direction SwapDirection) uint64 {
	var ret uint64 = 0
	var denominator *big.Int
	amountInBN := big.NewInt(0).SetUint64(amountInDeductFee)
	var amountOut = big.NewInt(0)
	switch direction {
	case Pc2Coin:
		denominator = big.NewInt(0).SetUint64(pcWithoutTakePnl)
		denominator.Add(denominator, amountInBN)
		multiplier := big.NewInt(0).SetUint64(coinWithoutTakePnl)
		multiplier.Mul(multiplier, amountInBN)
		ret = amountOut.Div(multiplier, denominator).Uint64()
	case Coin2Pc:
		denominator = big.NewInt(0).SetUint64(coinWithoutTakePnl)
		denominator.Add(denominator, amountInBN)
		multiplier := big.NewInt(0).SetUint64(pcWithoutTakePnl)
		multiplier.Mul(multiplier, amountInBN)
		ret = amountOut.Div(multiplier, denominator).Uint64()
	}

	return ret

}

func calcTotalWithoutTakePnlNoOrderBook(pcAmount, coinAmount uint64, amm *AmmInfo) (uint64, uint64, error) {
	pcWithoutTakePnl := pcAmount - amm.StateData.NeedTakePnlPc
	if pcWithoutTakePnl < 0 {
		return 0, 0, fmt.Errorf("pcAmount is less than NeedTakePnlPc")
	}
	coinWithoutTakePnl := coinAmount - amm.StateData.NeedTakePnlCoin
	if coinWithoutTakePnl < 0 {
		return 0, 0, fmt.Errorf("coinAmount is less than NeedTakePnlCoin")
	}
	return pcWithoutTakePnl, coinWithoutTakePnl, nil
}
