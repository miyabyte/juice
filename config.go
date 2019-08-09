package juice

import "time"

type Config struct {
	Addr string
	HandlerFuncPattern string

	ReadBufferSize uint16
	WriteBufferSize uint16

	HeartbeatCheckInterval time.Duration
	HeartbeatIdleTime time.Duration
}