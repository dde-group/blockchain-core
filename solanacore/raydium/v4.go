package raydium

import (
	"bytes"
	"encoding/binary"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type SwapInstructionV4 struct {
	bin.BaseVariant
	InAmount                uint64
	MinimumOutAmount        uint64
	solana.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (inst *SwapInstructionV4) ProgramID() solana.PublicKey {
	return solana.MustPublicKeyFromBase58(LiquidityPoolProgramV4)
}

func (inst *SwapInstructionV4) Accounts() (out []*solana.AccountMeta) {
	return inst.Impl.(solana.AccountsGettable).GetAccounts()
}

func (inst *SwapInstructionV4) Data() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := bin.NewBorshEncoder(buf).Encode(inst); err != nil {
		return nil, fmt.Errorf("unable to encode instruction: %w", err)
	}
	return buf.Bytes(), nil
}

func (inst *SwapInstructionV4) MarshalWithEncoder(encoder *bin.Encoder) (err error) {
	// Swap instruction is number 9
	err = encoder.WriteUint8(9)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.InAmount, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteUint64(inst.MinimumOutAmount, binary.LittleEndian)
	if err != nil {
		return err
	}
	return nil
}

func NewSwapInstructionV4(
	inAmount uint64,
	minimumOutAmount uint64,
	tokenProgram solana.PublicKey,
	pool *PoolInfo,
	userSourceTokenAccount solana.PublicKey,
	userDestTokenAccount solana.PublicKey,
	userOwner solana.PublicKey) *SwapInstructionV4 {
	inst := SwapInstructionV4{
		InAmount:         inAmount,
		MinimumOutAmount: minimumOutAmount,
		AccountMetaSlice: make(solana.AccountMetaSlice, 18),
	}
	inst.BaseVariant = bin.BaseVariant{
		Impl: inst,
	}

	ammId := solana.MustPublicKeyFromBase58(pool.AmmId)
	ammAuthority := solana.MustPublicKeyFromBase58(pool.AmmAuthority)
	ammOpenOrders := solana.MustPublicKeyFromBase58(pool.AmmOpenOrders)
	ammTargetOrders := solana.MustPublicKeyFromBase58(pool.AmmTargetOrders)
	poolCoinTokenAccount := solana.MustPublicKeyFromBase58(pool.AmmCoinTokenAccount)
	poolPcTokenAccount := solana.MustPublicKeyFromBase58(pool.AmmPcTokenAccount)
	serumProgramId := solana.MustPublicKeyFromBase58(pool.SerumProgramId)
	serumMarket := solana.MustPublicKeyFromBase58(pool.SerumMarket)
	serumBids := solana.MustPublicKeyFromBase58(pool.SerumBids)
	serumAsks := solana.MustPublicKeyFromBase58(pool.SerumAsks)
	serumEventQueue := solana.MustPublicKeyFromBase58(pool.SerumEventQueue)
	serumCoinVaultAccount := solana.MustPublicKeyFromBase58(pool.SerumCoinVaultAccount)
	serumPcVaultAccount := solana.MustPublicKeyFromBase58(pool.SerumPcVaultAccount)
	serumVaultSigner := solana.MustPublicKeyFromBase58(pool.SerumVaultSigner)

	inst.AccountMetaSlice[0] = solana.Meta(tokenProgram)
	inst.AccountMetaSlice[1] = solana.Meta(ammId).WRITE()
	inst.AccountMetaSlice[2] = solana.Meta(ammAuthority)
	inst.AccountMetaSlice[3] = solana.Meta(ammOpenOrders).WRITE()
	inst.AccountMetaSlice[4] = solana.Meta(ammTargetOrders).WRITE()
	inst.AccountMetaSlice[5] = solana.Meta(poolCoinTokenAccount).WRITE()
	inst.AccountMetaSlice[6] = solana.Meta(poolPcTokenAccount).WRITE()
	inst.AccountMetaSlice[7] = solana.Meta(serumProgramId)
	inst.AccountMetaSlice[8] = solana.Meta(serumMarket).WRITE()
	inst.AccountMetaSlice[9] = solana.Meta(serumBids).WRITE()
	inst.AccountMetaSlice[10] = solana.Meta(serumAsks).WRITE()
	inst.AccountMetaSlice[11] = solana.Meta(serumEventQueue).WRITE()
	inst.AccountMetaSlice[12] = solana.Meta(serumCoinVaultAccount).WRITE()
	inst.AccountMetaSlice[13] = solana.Meta(serumPcVaultAccount).WRITE()
	inst.AccountMetaSlice[14] = solana.Meta(serumVaultSigner)
	inst.AccountMetaSlice[15] = solana.Meta(userSourceTokenAccount).WRITE()
	inst.AccountMetaSlice[16] = solana.Meta(userDestTokenAccount).WRITE()
	inst.AccountMetaSlice[17] = solana.Meta(userOwner).SIGNER()

	return &inst
}
