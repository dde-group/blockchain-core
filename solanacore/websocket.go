package solanacore

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
)

type WebsocketAgent struct {
	client   *ws.Client
	endpoint string
	isAlive  bool
}

func NewWebsocketAgent(endpoint string) *WebsocketAgent {
	ret := &WebsocketAgent{
		endpoint: endpoint,
	}

	return ret
}

func (agent *WebsocketAgent) Connect() error {
	client, err := ws.Connect(context.TODO(), agent.endpoint)
	if nil != err {
		return fmt.Errorf("ws connect err: %s", err.Error())
	}
	agent.client = client
	agent.isAlive = true
	return nil
}

func (agent *WebsocketAgent) LogsSubscribeMentions(mentions solana.PublicKey, commitment rpc.CommitmentType, handler func(*ws.LogResult)) (*ws.LogSubscription, error) {
	sub, err := agent.client.LogsSubscribeMentions(mentions, commitment)
	if nil != err {
		return nil, fmt.Errorf("subcribe err: %s", err.Error())
	}

	go func() {
		var result *ws.LogResult
		for {
			result, err = sub.Recv()
			if nil != err {
				//TODO reconnect
				//agent.isAlive = false
				break
			}

			if nil != handler {
				handler(result)
			}
		}
	}()
	return sub, nil
}

func (agent *WebsocketAgent) BlockSubscribeMentions(mentions solana.PublicKey, commitment rpc.CommitmentType, handler func(*ws.BlockResult)) (*ws.BlockSubscription, error) {
	maxVersion := uint64(0)
	sub, err := agent.client.BlockSubscribe(
		&ws.BlockSubscribeFilterMentionsAccountOrProgram{
			Pubkey: mentions,
		},
		//ws.BlockSubscribeFilterAll(""),
		&ws.BlockSubscribeOpts{
			Commitment:                     commitment,
			MaxSupportedTransactionVersion: &maxVersion,
			TransactionDetails:             "full",
			Encoding:                       solana.EncodingBase64,
		},
	)
	if nil != err {
		return nil, fmt.Errorf("subcribe err: %s", err.Error())
	}

	go func() {
		var result *ws.BlockResult
		for {
			result, err = sub.Recv()
			if nil != err {
				//TODO reconnect
				break
			}

			if nil != handler {
				handler(result)
			}
		}
	}()

	return sub, nil

}

func (agent *WebsocketAgent) ParsedBlockSubscribeMentions(
	mentions solana.PublicKey,
	commitment rpc.CommitmentType,
	handler func(result *ws.ParsedBlockResult)) (*ws.ParsedBlockSubscription, error) {

	maxVersion := uint64(0)
	sub, err := agent.client.ParsedBlockSubscribe(
		&ws.BlockSubscribeFilterMentionsAccountOrProgram{
			Pubkey: mentions,
		},
		//ws.BlockSubscribeFilterAll(""),
		&ws.BlockSubscribeOpts{
			Commitment:                     commitment,
			MaxSupportedTransactionVersion: &maxVersion,
			TransactionDetails:             "full",
			Encoding:                       solana.EncodingJSONParsed,
		},
	)
	if nil != err {
		return nil, fmt.Errorf("subcribe err: %s", err.Error())
	}

	go func() {
		var result *ws.ParsedBlockResult
		for {
			result, err = sub.Recv()
			if nil != err {
				//TODO reconnect
				break
			}

			if nil != handler {
				handler(result)
			}
		}
	}()

	return sub, nil

}
