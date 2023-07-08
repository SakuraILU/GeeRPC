package service

import (
	"log"
	"reflect"
)

type Method struct {
	Name   string
	Fun    reflect.Value
	Argt   reflect.Type
	Replyt reflect.Type
}

func NewMethod(method reflect.Method) *Method {
	return &Method{
		Name:   method.Name,
		Fun:    method.Func,
		Argt:   method.Type.In(1),
		Replyt: method.Type.In(2),
	}
}

func (m *Method) NewArgv() (argv reflect.Value) {
	if m.Argt.Kind() == reflect.Ptr {
		// m.argt    -->  .Elem()  -->   relect.New
		// 	*int type     int type	     *int value
		argv = reflect.New(m.Argt.Elem())
	} else {
		// m.argt  -->   reflect.New --> .Elem()
		// int type     *int value       int value
		argv = reflect.New(m.Argt).Elem()
	}
	return
}

func (m *Method) NewReply() (replyv reflect.Value) {
	if m.Replyt.Kind() != reflect.Ptr {
		log.Fatal("reply type must be ptr")
	}
	// m.argt    -->  .Elem()  -->   relect.New
	// 	*int type     int type	     *int value
	replyv = reflect.New(m.Replyt.Elem())
	switch m.Replyt.Elem().Kind() {
	case reflect.Slice:
		// arr := make([]int, 0, 0)
		replyv.Elem().Set(reflect.MakeSlice(m.Replyt.Elem(), 0, 0))
	case reflect.Map:
		// arr := make(map[int]string)
		replyv.Elem().Set(reflect.MakeMap(m.Replyt.Elem()))
		// case reflect.Chan:
		// arr := make(chan bool, 0)
		// replyv.Set(reflect.MakeChan(m.replyt.Elem(), 0))
	}
	return
}
