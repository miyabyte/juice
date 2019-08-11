package juice

import (
	"net/http"
	"sync"
	"time"
)

var c *Config
var onceConf sync.Once

type Config struct {
	Addr               string
	HandlerFuncPattern string

	ReadBufferSize  int
	WriteBufferSize int

	HeartbeatCheckInterval time.Duration
	HeartbeatIdleTime      time.Duration

	CheckOrigin func(r *http.Request) bool

	EnableAnalyzeUid bool
}

type EnableAnalyzeUid interface {
	AnalyzeUid(r *http.Request) int
}

func setConfig(config *Config) {
	onceConf.Do(func() {
		c = config
	})
}

func GetConfig() *Config {
	return c
}
