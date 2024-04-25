package network

import "time"

const (
	GrpcDialTimeout = 3 * time.Second
)

type GrpcCfg struct {
	Name string
	Host string
	Port uint
}
