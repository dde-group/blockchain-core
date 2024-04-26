package solanacore

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/goccy/go-json"
	"github.com/shopspring/decimal"
)

type ATAAccountInfo struct {
	Account solana.PublicKey //账户地址
	Mint    solana.PublicKey //token 地址
	Owner   solana.PublicKey //owner
}

type TokenInfo struct {
	Mint     solana.PublicKey //地址
	Decimals uint8            //精度
}

type BlockParsedTransaction struct {
	Transaction     *rpc.ParsedTransaction
	TransactionMeta *rpc.TransactionMeta
}

type parsedInstructionInfo struct {
	InstructionType string          `json:"type"`
	Info            instructionInfo `json:"info"`
}

type instructionInfo struct {
	Amount      json.Number `json:"amount"`
	Lamports    json.Number `json:"lamports"`
	Source      string      `json:"source"`
	Destination string      `json:"destination"`
	Authority   string      `json:"authority"`
}

type ParsedInstructionInfo struct {
	ProgamId        solana.PublicKey `json:"progamId"`
	InstructionType string           `json:"type"`
	Amount          decimal.Decimal  `json:"amount"`
	Lamports        decimal.Decimal  `json:"lamports"`
	Source          solana.PublicKey `json:"source"`
	Destination     solana.PublicKey `json:"destination"`
	Authority       solana.PublicKey `json:"authority"`
}

func ParserInstructionInfo(inst *rpc.ParsedInstruction) (*ParsedInstructionInfo, error) {
	inBytes, _ := inst.Parsed.MarshalJSON()
	info := parsedInstructionInfo{}
	err := json.Unmarshal(inBytes, &info)
	if nil != err {
		return nil, err
	}

	amount, _ := decimal.NewFromString(info.Info.Amount.String())
	lamports, _ := decimal.NewFromString(info.Info.Lamports.String())
	source, err := solana.PublicKeyFromBase58(info.Info.Source)
	if nil != err {
		return nil, fmt.Errorf("source account invalid, %s", info.Info.Source)
	}
	destionation, err := solana.PublicKeyFromBase58(info.Info.Destination)
	if nil != err {
		return nil, fmt.Errorf("source account invalid, %s", info.Info.Destination)
	}
	authority, err := solana.PublicKeyFromBase58(info.Info.Authority)
	if nil != err {
		return nil, fmt.Errorf("source account invalid, %s", info.Info.Authority)
	}

	ret := &ParsedInstructionInfo{
		ProgamId:        inst.ProgramId,
		InstructionType: info.InstructionType,
		Amount:          amount,
		Lamports:        lamports,
		Source:          source,
		Destination:     destionation,
		Authority:       authority,
	}

	return ret, nil
}

type SwapInstructionResult struct {
	Slot           uint64
	Interact       solana.PublicKey //交易所account
	Index          uint64           //tx index，因为订阅的不是全部交易，因此仅仅是相关交易所的交易的index
	InstIndex      uint64           //instruction index
	SubIndex       uint64           //inner instruction index
	Hash           solana.Signature
	AmmId          solana.PublicKey
	AmmCoinAccount solana.PublicKey
	AmmPcAccount   solana.PublicKey
	AmountIn       decimal.Decimal  //uint64 存储 没计算精度
	AmountOut      decimal.Decimal  //uint64 存储
	InFromAccount  solana.PublicKey //in 数据
	InToAccount    solana.PublicKey
	InAuthority    solana.PublicKey
	OutFromAccount solana.PublicKey //out 数据
	OutToAccount   solana.PublicKey
	OutAuthority   solana.PublicKey
	SolPrice       decimal.Decimal
}

func ParseInstruction2SwapResult(instList []*rpc.ParsedInstruction) (*SwapInstructionResult, error) {
	if len(instList) < 2 {
		return nil, fmt.Errorf("inst list length is less than 2")
	}

	baseIn, err := ParserInstructionInfo(instList[0])
	if nil != err {
		return nil, fmt.Errorf("swap in instruction parse err: %s", err.Error())
	}
	//非转账inst不处理
	if !baseIn.ProgamId.Equals(solana.TokenProgramID) || InstTypeTransfer != baseIn.InstructionType {
		return nil, fmt.Errorf("swap in insturction invalid")
	}

	baseOut, err := ParserInstructionInfo(instList[1])
	if nil != err {
		return nil, fmt.Errorf("swap out instruction parse err: %s", err.Error())
	}
	//非转账inst不处理
	if !baseOut.ProgamId.Equals(solana.TokenProgramID) || InstTypeTransfer != baseOut.InstructionType {
		return nil, fmt.Errorf("swap out insturction invalid")
	}

	ret := &SwapInstructionResult{
		AmountIn:       baseIn.Amount,
		InFromAccount:  baseIn.Source,
		InToAccount:    baseIn.Destination,
		InAuthority:    baseIn.Authority,
		AmountOut:      baseOut.Amount,
		OutFromAccount: baseOut.Source,
		OutToAccount:   baseOut.Destination,
		OutAuthority:   baseOut.Authority,
	}

	return ret, nil
}
