package juice

import (
	"context"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Client struct {
	sync.Mutex
	conn     *websocket.Conn
	LastTime time.Time
	UUID     uint32
	Ctx      context.Context
	Cancel   context.CancelFunc

	Uid interface{}
}
