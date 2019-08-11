package juice

import (
	"net/http"
	"time"
)

type Config struct {
	Addr               string
	HandlerFuncPattern string

	ReadBufferSize  int
	WriteBufferSize int

	HeartbeatCheckInterval time.Duration
	HeartbeatIdleTime      time.Duration

	CheckOrigin func(r *http.Request) bool
}
