package juice

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Juice struct {
	Event Event
	Conf Config
	Log
	*cliManager
}


type Event interface {
	Open(cli *Client,r *http.Request) error
	Close(cli *Client)
	Message(cli *Client, p []byte)
}

func (j *Juice) Exec() {
	j.cliManager = GetCliManager()
	http.HandleFunc(j.Conf.HandlerFuncPattern, j.initialize)

	if err := j.heartbeat();err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(j.Conf.Addr, nil))
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (j *Juice) initialize(w http.ResponseWriter, r *http.Request) {
	var (
		conn   *websocket.Conn
		err    error
		client *Client
	)

	if conn, err = upGrader.Upgrade(w, r, nil); err != nil {
		j.Cmd(err)
		return
	}

	client, err = NewClient(conn)
	if err != nil {
		return
	}

	//lifecycle
	// open handler
	if err := j.Event.Open(client, r); err != nil {
		j.Cmd(conn.Close())
		return
	}

	// close handler    [parameter client
	conn.SetCloseHandler(func(closeCode int, closeText string) error {
		_ = j.RemoveClient(j, client)
		return nil
	})

	j.AddClient(j, client)

	go j.getMessage(j, client)

	//heartbeat

}

func (j *Juice) Close(c *Client) {
	_ = c.conn.Close()
	_ = j.RemoveClient(j, c)
}

func (j *Juice) heartbeat() (err error) {
	hb := &heartbeat{j.Conf.HeartbeatCheckInterval, j.Conf.HeartbeatIdleTime}
	if err = hb.run(j);err !=nil {
		return
	}
	return
}
