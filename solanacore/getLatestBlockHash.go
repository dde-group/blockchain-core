package solanacore

import (
	"context"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/juju/ratelimit"
	"time"
)

type getLatestBlockHashHandler struct {
	bucket *ratelimit.Bucket
	result *rpc.GetLatestBlockhashResult
	client *rpc.Client
}

func newGetLatestBlockHashHandler(client *rpc.Client) *getLatestBlockHashHandler {
	ret := &getLatestBlockHashHandler{
		//500ms 只调用一次 时间内用缓存
		bucket: ratelimit.NewBucketWithQuantum(500*time.Millisecond, 1, 1),
		client: client,
	}

	return ret
}

func (h *getLatestBlockHashHandler) getResult(ctx context.Context, commitment rpc.CommitmentType) (*rpc.GetLatestBlockhashResult, error) {
	if h.bucket.TakeAvailable(1) > 0 {
		if nil != h.result {
			return h.result, nil
		}
	}

	result, err := h.client.GetLatestBlockhash(ctx, commitment)
	if nil != err {
		return nil, err
	}
	h.result = result
	return result, nil
}
