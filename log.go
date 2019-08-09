package juice

import (
	"fmt"
	"log"
)

type Log struct {
}

func (Log) FileLog(v interface{}) {
	fmt.Println(v)
}

func (Log) Cmd(v interface{}) {
	log.Println(v)
}
