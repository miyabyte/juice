# easy websocket 

```
ws := &juice.Juice{Conf: juice.Config{
		Addr:                   "localhost:8000",
		HandlerFuncPattern:     "/ws",
		ReadBufferSize:         1024,
		WriteBufferSize:        1024,
		HeartbeatCheckInterval: himConf.HeartbeatCheckInterval,
		HeartbeatIdleTime:      himConf.HeartbeatIdleTime,
	}}
	
ws.SetEvent(&websocket.Event{})
log.Println(ws.Exec())
ws.Cmd("listen ws on : 8000")
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
