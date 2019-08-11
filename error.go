package juice

import "fmt"

type JError struct {
	code int
	msg  string
}

func NewJError(code int, msg string) JError {
	return JError{
		code: code,
		msg:  msg,
	}
}

func (je *JError) Error() string {
	return fmt.Sprintf("[%d] %s", je.code, je.msg)
}
