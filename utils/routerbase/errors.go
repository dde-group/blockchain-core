package routerbase

/**
 * @Author: lee
 * @Description:
 * @File: errors
 * @Date: 2022/4/1 6:50 下午
 */
const Success = 200
const (
	ErrorRateLimit = iota + 1001
	ErrorPermissionDenied
)

const (
	ErrorInvalidParams = iota + 2001
	ErrorInvalidRequest
	ErrorNoData
)

const (
	ErrorInvalidRaftNode = iota + 3001
	ErrorRaftJoinFailed
	ErrorRaftNotALeader
)

var ErrorMessage = map[int]string{
	Success:        "success",
	ErrorRateLimit: "request has reached the rate limit, please retry later",

	ErrorPermissionDenied: "permission denied",
	ErrorInvalidParams:    "invalid parameters",
	ErrorInvalidRequest:   "invalid request",
	ErrorNoData:           "no data found",

	ErrorInvalidRaftNode: "raft node has not launch",
	ErrorRaftJoinFailed:  "raft join cluster failed",
}
