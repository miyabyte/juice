package juice_test

import (
	websocket2 "github.com/gorilla/websocket"
	"him/chat/juice"
	"him/chat/socket/websocket"
	"testing"
)

func getJuice() *juice.Juice {
	return &juice.Juice{Event: &websocket.Event{}}
}

func TestJuice_Exec(t *testing.T) {
	ws := &juice.Juice{Event: &websocket.Event{}, Conf: juice.Config{
		Addr:                   string("localhost:8000"),
		HandlerFuncPattern:     "/ws",
		ReadBufferSize:         1024,
		WriteBufferSize:        1024,
		HeartbeatCheckInterval: 10,
		HeartbeatIdleTime:      10,
	}}
	ws.Exec()
}

func TestCliManager_AddClient(t *testing.T) {
	cids := make([]uint32, 0)

	cliM := juice.GetCliManager()

	for i := 0; i < 100; i++ {
		client, _ := juice.NewClient(&websocket2.Conn{})
		cliM.AddClient(getJuice(), client)
		cids = append(cids, client.UUID)
	}

	clients := cliM.GetClients()
	for k, v := range cids {
		if v, f := clients[v]; f == true {
			t.Logf("num: %d , succ, v : %v \n", k, v)
		} else {
			t.Error(k, v, f)
			return
		}
	}

	t.Log(clients)
}
