package solanacore

import (
	"context"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/juju/ratelimit"
	"sync"
	"time"
)

type getMinimumBalanceForRentExemptionHandler struct {
	bucket      map[uint64]*ratelimit.Bucket
	resultCache map[uint64]uint64
	client      *rpc.Client
	mtx         sync.Mutex
}

func newGetMinimumBalanceForRentExemptionSchedule(client *rpc.Client) *getMinimumBalanceForRentExemptionHandler {
	ret := &getMinimumBalanceForRentExemptionHandler{
		//500ms 只调用一次 时间内用缓存
		bucket:      make(map[uint64]*ratelimit.Bucket),
		resultCache: make(map[uint64]uint64),
		client:      client,
	}

	return ret
}

func (h *getMinimumBalanceForRentExemptionHandler) getResult(ctx context.Context, dataSize uint64) (uint64, error) {
	h.mtx.Lock()
	defer h.mtx.Unlock()
	bucket, ok := h.bucket[dataSize]
	if !ok {
		bucket = ratelimit.NewBucketWithQuantum(1000*time.Millisecond, 1, 1)
		h.bucket[dataSize] = bucket
	}
	if bucket.TakeAvailable(1) > 0 {
		if result, ok := h.resultCache[dataSize]; ok {
			return result, nil
		}
	}

	result, err := h.client.GetMinimumBalanceForRentExemption(ctx, dataSize, rpc.CommitmentConfirmed)
	if nil != err {
		return 0, err
	}

	h.resultCache[dataSize] = result
	return result, nil
}
