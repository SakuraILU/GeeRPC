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
	// One thing worth noting: we don't need to care about Slice or Map,
	// because after we create refargv, we will convert it to argvp interface{} (if not a ptr type, cvt2ptr type)
	// and call ReadBody(arvp), which will call gob.Decode(arvp), it will create a new slice or map for us if necessary
	// but in NewReply(), no one will do this for us, so we need to do it by ourselves in NewReply()
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
	}
	return
}
