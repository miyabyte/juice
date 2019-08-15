package juice

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Client struct {
	sync.Mutex
	conn     *websocket.Conn
	LastTime time.Time
	UUID     uint32
	Cancel   context.CancelFunc

	Uid   int
	Ctx   context.Context
	Extra interface{}
}

func (c *Client) Info() (info string) {
	info = fmt.Sprintf(
		"[%d] %d : %s \n ",
		c.UUID,
		c.Uid,
		c.LastTime.Format("2006-01-02 15:04:05.000"),
	)
	return
}
