package solanacore

import (
	"context"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
)

func (c *Core) GetAccountToken(ctx context.Context, account solana.PublicKey) (solana.PublicKey, error) {
	out, err := c.client.GetAccountInfo(ctx, account)
	if nil != err {
		return solana.PublicKey{}, fmt.Errorf("GetAccountInfo err: %s", err.Error())
	}

	tokenAcc := token.Account{}
	data := out.Value.Data.GetBinary()
	dec := bin.NewBinDecoder(data)
	err = dec.Decode(&tokenAcc)

	return out.Value.Owner, nil

}

func (c *Core) GetTokenMintInfo(ctx context.Context, tokenMint solana.PublicKey) (*token.Mint, error) {
	result, err := c.client.GetAccountInfo(ctx,
		tokenMint)
	if nil != err {
		return nil, fmt.Errorf("GetAccountInfo err: %s", err.Error())
	}

	mint := token.Mint{}

	if err = bin.NewBinDecoder(result.Value.Data.GetBinary()).Decode(&mint); nil != err {
		//if err = mint.Decode(result.Value.Data.GetBinary()); err != nil {
		return nil, fmt.Errorf("unable to retrieve mint information: %s", err.Error())
	}

	if !result.Value.Owner.Equals(solana.TokenProgramID) {
		return nil, fmt.Errorf("not a token mint address")
	}

	return &mint, nil
}
