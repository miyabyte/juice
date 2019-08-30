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

func Cmd(v interface{}) {
	log.Println(v)
}

func Panic(v interface{}) {
	log.Panic(v)
}
