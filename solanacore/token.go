package solanacore

import (
	"context"
	"fmt"
	"github.com/dde-group/blockchain-core/blockchain"
	"github.com/gagliardetto/solana-go"
	"github.com/shopspring/decimal"
	"math/big"
)

func (c *Core) GetTokenSupply(tokenAddr string) (*blockchain.TokenSupply, error) {
	tokenMint, err := solana.PublicKeyFromBase58(tokenAddr)
	if nil != err {
		return nil, fmt.Errorf("invalid tokenAddr, err: %s", err.Error())
	}
	mint, err := c.GetTokenMintInfo(context.TODO(), tokenMint)
	if nil != err {
		return nil, fmt.Errorf("GetTokenMintInfo err: %s", err.Error())
	}

	amount := big.NewInt(0).SetUint64(mint.Supply)
	amountDecimal := decimal.NewFromBigInt(amount, 1).Mul(decimal.NewFromBigInt(big.NewInt(1), int32(-mint.Decimals)))

	ret := &blockchain.TokenSupply{
		TokenAddress:   tokenAddr,
		Supply:         mint.Supply,
		Decimals:       uint(mint.Decimals),
		SupplyDecimals: amountDecimal.String(),
	}

	return ret, nil
}

func (c *Core) GetTokenDecimals(tokenAddr string) (uint, error) {
	ret, ok := c.decimalsCache.Load(tokenAddr)
	if !ok {
		tokenSupply, err := c.GetTokenSupply(tokenAddr)
		if nil != err {
			return 0, fmt.Errorf("GetTokenSupply err: %s", err.Error())
		}
		c.decimalsCache.Store(tokenAddr, tokenSupply.Decimals)
		return tokenSupply.Decimals, nil
	}

	return ret.(uint), nil
}

func (c *Core) GetAmountWithDecimals(amount decimal.Decimal, tokenAddr string) (uint64, error) {
	decimals := uint(SOLDecimals)
	var err error
	if tokenAddr != NativeSOL {
		decimals, err = c.GetTokenDecimals(tokenAddr)
		if nil != err {
			return 0, fmt.Errorf("get token decimals err: %s", err.Error())
		}
	}
	realAmount := amount.Mul(decimal.NewFromBigInt(solana.DecimalsInBigInt(uint32(decimals)), 0)).BigInt().Uint64()

	return realAmount, nil
}

func (c *Core) GetDecimalAmount(amount uint64, tokenAddr string) (decimal.Decimal, error) {
	decimals := uint(SOLDecimals)
	var err error

	if tokenAddr != NativeSOL {
		decimals, err = c.GetTokenDecimals(tokenAddr)
		if nil != err {
			return decimal.Zero, fmt.Errorf("get token decimals err: %s", err.Error())
		}
	}

	amountDecimal := big.NewInt(0).SetUint64(amount)
	ret := decimal.NewFromBigInt(amountDecimal, 0).Mul(decimal.NewFromBigInt(big.NewInt(1), int32(-decimals)))
	return ret, nil
}
