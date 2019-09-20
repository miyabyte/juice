package juice_test

import (
	websocket2 "github.com/gorilla/websocket"
	"juice"
	"log"
	"testing"
)

func TestJuice_Exec(t *testing.T) {
	ws, _ := juice.NewJuice(
		&juice.Config{
			Addr:                   string("localhost:8000"),
			HandlerFuncPattern:     "/ws",
			ReadBufferSize:         juice.ReadBufferSize,
			WriteBufferSize:        juice.WriteBufferSize,
			HeartbeatCheckInterval: juice.HeartbeatCheckInterval,
			HeartbeatIdleTime:      juice.HeartbeatIdleTime,
			//开启解析用户
			EnableAnalyzeUid: true,
			//开启自定义mux   示例：gin
			EnableChangeMux: true,
		},
		&juice.DefaultEvent{},
	)

	ws.Mux.HandleFunc("/ws/info", juice.WsInfo)

	log.Fatalln(ws.Exec())
}

func TestCliManager_AddClient(t *testing.T) {
	go TestJuice_Exec(t)

	cids := make([]uint32, 0)

	cliM := juice.NewCliManager(&juice.DefaultEvent{})

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
