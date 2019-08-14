package juice

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

var cm *cliManager
var onceCM sync.Once

type Observer interface {
	Online(cli *Client)
	Offline(cli *Client)
}

type cliManager struct {
	clients map[uint32]*Client
	Event

	observers []Observer
}

func GetCliManager(e Event) *cliManager {
	onceCM.Do(func() {
		cm = &cliManager{
			clients: make(map[uint32]*Client),
			Event:   e,
		}
	})
	return cm
}

func NewClient(conn *websocket.Conn) (cli *Client, err error) {
	if _UUID, err := uuid.NewUUID(); err != nil {
		return nil, err
	} else {
		ctx, cancel := context.WithCancel(context.Background())
		return &Client{
			conn:     conn,
			UUID:     _UUID.ID(),
			LastTime: time.Now(),
			Ctx:      ctx,
			Cancel:   cancel,
		}, nil
	}
}

func (cm *cliManager) GetClients() map[uint32]*Client {
	return cm.clients
}

func (cm *cliManager) GetClient(uuid uint32) (cli *Client, ok bool) {
	cli, ok = cm.clients[uuid]
	return
}

// up  observer mode
func (cm *cliManager) AddClient(cli *Client) *cliManager {
	if GetConfig().EnableAnalyzeUid {
		ucm := GetUserCliManager()
		//lock
		ucm.Lock()
		defer ucm.Unlock()

		ucm.AddClient(cli)
	}
	cm.clients[cli.UUID] = cli
	return cm
}

func (cm *cliManager) CloseClient(c *Client) {
	_ = c.conn.Close()
	_ = cm.RemoveClient(c)
}

// down observer mode
func (cm *cliManager) RemoveClient(cli *Client) *cliManager {
	if GetConfig().EnableAnalyzeUid {
		ucm := GetUserCliManager()
		//lock
		ucm.Lock()
		defer ucm.Unlock()

		ucm.RemoveClient(cli)
	}
	delete(cm.clients, cli.UUID)
	cm.Close(cli)
	cli.Cancel()
	return cm
}

func (cm *cliManager) getMessage(cli *Client) {
	conn := cli.conn
	for {
		// msgType 1 text 2 binary
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			cm.ErrorHandler(NewJError(ErrWsGetMsg, err.Error()))
			return
		}

		if messageType == websocket.TextMessage {
			cm.Event.Message(cli, p)
		}

		if messageType == websocket.BinaryMessage {
			cm.Event.BinaryMessage(cli, p)
		}

		// wm = nextWriter\write\close
		//if err := conn.WriteMessage(messageType, p); err != nil {
		//	j.Cmd(err)
		//	return
	}
}
