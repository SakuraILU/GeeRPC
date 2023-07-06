package codec

import "reflect"

type Request struct {
	Head   Head
	Argv   reflect.Value
	Replyv reflect.Value
}
