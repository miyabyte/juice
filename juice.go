package juice

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var j *Juice
var onceJ sync.Once

type Mux interface {
	http.Handler
	HandleFunc(string, func(http.ResponseWriter, *http.Request))
}

type Juice struct {
	event Event
	Conf  *Config
	Log
	cm *cliManager

	upGrader websocket.Upgrader
	Mux      Mux

	srv *http.Server
}

func NewJuice(conf *Config, e Event) *Juice {
	onceJ.Do(func() {
		j = &Juice{
			Conf:  conf,
			event: e,
			Mux:   http.NewServeMux(),
		}
	})
	return j
}

type Event interface {
	Open(cli *Client, r *http.Request) error
	//close when heartbeat or manually(CM) will trigger once
	Close(cli *Client)
	Message(cli *Client, p []byte)
	BinaryMessage(cli *Client, p []byte)
	ErrorHandler(err JError)
}

func (j *Juice) Exec() (err error) {
	setConfig(j.Conf)
	GetUserCliManager()

	if err = j.wsSet(); err != nil {
		return
	}

	j.Mux.HandleFunc(j.Conf.HandlerFuncPattern, j.initialize)

	if err = j.heartbeat(); err != nil {
		return
	}

	j.srv = &http.Server{Addr: j.Conf.Addr, Handler: j.Mux}
	if err = j.srv.ListenAndServe(); err != nil {
		return
	}
	return
}

func (j *Juice) Close() error {
	return j.srv.Close()
}

func (j *Juice) initialize(w http.ResponseWriter, r *http.Request) {
	var (
		conn   *websocket.Conn
		err    error
		client *Client
	)

	if conn, err = j.upGrader.Upgrade(w, r, nil); err != nil {
		j.Cmd(err)
		return
	}

	client, err = NewClient(conn)
	if err != nil {
		return
	}

	//lifecycle
	// analyzeUid handler
	if j.Conf.EnableAnalyzeUid {
		if client.Uid, err = j.event.(EnableAnalyzeUid).AnalyzeUid(r); err != nil {
			j.Cmd(err)
			conn.Close()
			return
		}
	}

	// 	onopen handler
	if err = j.event.Open(client, r); err != nil {
		j.Cmd(conn.Close())
		return
	}

	// onclose handler
	conn.SetCloseHandler(func(closeCode int, closeText string) error {
		j.cm.RemoveClient(client)
		return nil
	})

	j.cm.AddClient(client)

	//  onmessage handler
	go j.cm.getMessage(client)
}

func (j *Juice) heartbeat() (err error) {
	hb := &heartbeat{j.Conf.HeartbeatCheckInterval, j.Conf.HeartbeatIdleTime}
	if err = hb.run(j.cm); err != nil {
		return
	}
	return
}

func (j *Juice) wsSet() (err error) {
	//upGrader
	var upGrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	if j.Conf.ReadBufferSize != 0 && j.Conf.WriteBufferSize != 0 {
		upGrader.ReadBufferSize = j.Conf.ReadBufferSize
		upGrader.WriteBufferSize = j.Conf.WriteBufferSize
	}

	if j.Conf.CheckOrigin != nil {
		upGrader.CheckOrigin = j.Conf.CheckOrigin
	}

	j.upGrader = upGrader

	//event
	if j.event == nil {
		j.event = &DefaultEvent{}
	}

	if j.Conf.EnableAnalyzeUid {
		if _, ok := j.event.(EnableAnalyzeUid); !ok {
			return &JError{
				code: ErrEnableAnalyzeUid,
				msg:  "your EVENT must implement AnalyzeUid interface",
			}
		}
	}

	j.cm = NewCliManager(j.event)

	return
}
