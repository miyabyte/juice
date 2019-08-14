package juice

import (
	"fmt"
	"net/http"
)

type DefaultEvent struct {
}

func (e *DefaultEvent) Open(cli *Client, r *http.Request) error {
	//fmt.Println(r.Header)
	fmt.Println("open")
	return nil
}

func (e *DefaultEvent) AnalyzeUid(r *http.Request) (uid int) {
	//fmt.Println(r.Header)
	fmt.Println("az uid")
	return 2760
}

func (e *DefaultEvent) Close(cli *Client) {
	fmt.Println("ev close")
	return

}

func (e *DefaultEvent) Message(cli *Client, p []byte) {
	fmt.Println("ev msg")
	us := getOnlineUsers()
	for k, v := range us {
		fmt.Println(k, v)
	}
	return
}

func (e *DefaultEvent) BinaryMessage(cli *Client, p []byte) {
	fmt.Println("ev msg")
	return
}

func (e *DefaultEvent) ErrorHandler(err JError) {
	return
}
