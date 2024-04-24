package solanacore

import (
	"context"
	"github.com/gagliardetto/solana-go/rpc"
)

func (c *Core) GetLatestBlockHash(ctx context.Context, commitment rpc.CommitmentType) (*rpc.GetLatestBlockhashResult, error) {
	return c.latestBlockHashSchedule.getResult(ctx, commitment)
}

func (c *Core) GetMinimumBalanceForRentExemption(ctx context.Context, dataSize uint64) (uint64, error) {
	return c.minimumBalanceForRentExemptionHandler.getResult(ctx, dataSize)
}
