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
