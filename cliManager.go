package juice

import (
	"context"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"time"
)

var cm *cliManager

type cliManager struct {
	clients map[uint32]*Client
	Event
}

func GetCliManager(e Event) *cliManager {
	if cm == nil {
		cm = &cliManager{
			clients: make(map[uint32]*Client),
			Event:   e,
		}
	}
	return cm
}

func NewClient(conn *websocket.Conn) (client *Client, err error) {
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

func (cm *cliManager) GetClient(uuid uint32) (cli *Client, flag bool) {
	cli, flag = cm.clients[uuid]
	return
}

// up
func (cm *cliManager) AddClient(cli *Client) *cliManager {
	cm.clients[cli.UUID] = cli
	return cm
}

func (cm *cliManager) CloseClient(c *Client) {
	_ = c.conn.Close()
	_ = cm.RemoveClient(c)
}

// down
func (cm *cliManager) RemoveClient(cli *Client) *cliManager {
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
