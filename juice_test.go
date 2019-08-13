package juice_test

import (
	websocket2 "github.com/gorilla/websocket"
	"github.com/pengshuaifei/juice"
	"log"
	"testing"
)

func getJuice() *juice.Juice {
	return &juice.Juice{}
}

func TestJuice_Exec(t *testing.T) {
	ws := &juice.Juice{Conf: juice.Config{
		Addr:                   string("localhost:8000"),
		HandlerFuncPattern:     "/ws",
		ReadBufferSize:         juice.ReadBufferSize,
		WriteBufferSize:        juice.WriteBufferSize,
		HeartbeatCheckInterval: juice.HeartbeatCheckInterval,
		HeartbeatIdleTime:      juice.HeartbeatIdleTime,

		EnableAnalyzeUid: true,
	}}
	ws.SetEvent(&juice.DefaultEvent{})
	log.Fatalln(ws.Exec())
}

func TestCliManager_AddClient(t *testing.T) {
	cids := make([]uint32, 0)

	cliM := juice.GetCliManager(&juice.DefaultEvent{})

	for i := 0; i < 100; i++ {
		client, _ := juice.NewClient(&websocket2.Conn{})
		cliM.AddClient(client)
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
