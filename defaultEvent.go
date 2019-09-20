package juice

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type DefaultEvent struct {
}

func (e *DefaultEvent) Open(cli *Client, r *http.Request) error {
	//fmt.Println(r.Header)
	fmt.Println("open")
	return nil
}

type ginMux struct {
	*gin.Engine
}

func (g *ginMux) HandleFunc(path string, f func(w http.ResponseWriter, r *http.Request)) {
	g.GET(path, func(c *gin.Context) {
		f(c.Writer, c.Request)
	})
}

func (e *DefaultEvent) NewMux() Mux {
	return &ginMux{gin.Default()}
}

func (e *DefaultEvent) AnalyzeUid(r *http.Request) (uid int, err error) {
	//fmt.Println(r.Header)
	fmt.Println("az uid")
	return 2760, nil
}

func (e *DefaultEvent) Close(cli *Client) {
	fmt.Println("ev close")
	return

}

func (e *DefaultEvent) Message(cli *Client, p []byte) {
	fmt.Println("ev msg")
	return
}

func (e *DefaultEvent) BinaryMessage(cli *Client, p []byte) {
	fmt.Println("ev msg")
	return
}

func (e *DefaultEvent) ErrorHandler(err JError) {
	return
}

func WsInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("info/index"))

	fmt.Println("===")
	cs := GetCliManager().GetClients()
	ucs := GetUserCliManager().GetUserClients()
	ous := GetUserCliManager().GetOnlineUsers()
	for _, c := range cs {
		fmt.Printf("fd : %d , uid %d \n", c.UUID, c.Uid)
	}
	fmt.Println("===")
	for uid, uc := range ucs {
		var cFdsInfo string
		for _, cs := range uc.Clients {
			cFdsInfo = cFdsInfo + cs.Info()
		}
		fmt.Printf("uid@%d :\n %s \n", uid, cFdsInfo)
	}
	fmt.Println("===")
	for uid, ou := range ous {
		fmt.Printf("uid : %d , online_t : %s \n", uid, ou.Format(time.StampMilli))
	}
	return
}
