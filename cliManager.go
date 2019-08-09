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
}

func GetCliManager() *cliManager {
	if cm == nil {
		cm = &cliManager{
			clients: make(map[uint32]*Client),
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
			Cancel: cancel,
		}, nil
	}
}

func (cm *cliManager) GetClients () map[uint32]*Client {
	return cm.clients
}

func (cm *cliManager) GetClient (uuid uint32) (cli *Client,flag bool) {
	cli,flag = cm.clients[uuid]
	return
}

func (cm *cliManager) AddClient(j *Juice, cli *Client) *cliManager {
	cli.Lock()
	defer cli.Unlock()

	cm.clients[cli.UUID] = cli
	return cm
}

func (cm *cliManager) RemoveClient(j *Juice, cli *Client) *cliManager {
	cli.Lock()
	defer cli.Unlock()

	delete(cm.clients, cli.UUID)
	j.Event.Close(cli)
	cli.Cancel()
	return cm
}

func (cm *cliManager) getMessage(j *Juice, cli *Client) {
	conn := cli.conn
	for {
		// msgType 1 text 2 binary
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if messageType == websocket.TextMessage {
			j.Event.Message(cli, p)
		}

		if messageType == websocket.BinaryMessage {
			// 二进制
		}

		// wm = nextWriter\write\close
		//if err := conn.WriteMessage(messageType, p); err != nil {
		//	j.Cmd(err)
		//	return
	}
}
