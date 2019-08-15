# easy websocket 

```
ws := juice.NewJuice(
		juice.Config{
			Addr:                   string("localhost:8000"),
			HandlerFuncPattern:     "/ws",
			ReadBufferSize:         juice.ReadBufferSize,
			WriteBufferSize:        juice.WriteBufferSize,
			HeartbeatCheckInterval: juice.HeartbeatCheckInterval,
			HeartbeatIdleTime:      juice.HeartbeatIdleTime,

			EnableAnalyzeUid: true,
		},
		&juice.DefaultEvent{},
	)

	ws.Mux.HandleFunc("/ws/info", juice.WsInfo)

	log.Fatalln(ws.Exec())
```

```
type Event struct {
}

func (e *Event) Open(cli *juice.Client,r *http.Request) error {
	fmt.Println(r.Header)
	fmt.Println("open")
	return nil
}

func (e *Event) Close(cli *juice.Client) {
	fmt.Println("ev close")
	return

}

func (e *Event) Message(cli *juice.Client, p []byte) {
	fmt.Println("ev msg")
	return
}

func (e *Event) BinaryMessage(cli *juice.Client, p []byte) {
	fmt.Println("ev msg")
	return
}

func (e *Event) ErrorHandler(err juice.JError) {

}
```
### options:
当开启选项时需满足对应接口
- juice.config.EnableAnalyzeUid  :  true 
