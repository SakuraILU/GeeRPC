package service

import (
	"testing"
)

// define a simple add struct service
type AddService struct {
}

type ArgType struct {
	A, B int
}

func (s *AddService) Add(argv ArgType, reply *int) error {
	*reply = argv.A + argv.B
	return nil
}

func Test1(t *testing.T) {
	// log.SetFlags(log.Lshortfile | log.LstdFlags)
	// addsvic := &AddService{}
	// svic := NewService(addsvic)

	// argv := reflect.ValueOf(ArgType{A: 1, B: 2})
	// replyv := reflect.ValueOf(new(int))

	// err := svic.Call("Add", argv, replyv)
	// if err != nil {
	// 	t.Error(err)
	// }
	// if replyv.Elem().Int() != 3 {
	// 	t.Error("Add method error")
	// }
}
