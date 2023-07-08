package service

import (
	"errors"
	"fmt"
	"go/ast"
	"log"
	"reflect"
)

type Service struct {
	name    string
	rfvalue reflect.Value
	methods map[string]*Method
}

func NewService(svice interface{}) (s *Service, err error) {
	// service must be exportable
	s = &Service{}
	s.rfvalue = reflect.ValueOf(svice)
	s.name = reflect.Indirect(s.rfvalue).Type().Name()
	// service must be exportable
	if !ast.IsExported(s.name) {
		return nil, fmt.Errorf("service %s must be exportable", s.name)
	}
	s.registerMethods(svice)
	return s, nil
}

func (s *Service) registerMethods(svice interface{}) {
	stp := reflect.TypeOf(svice)
	// fix a bug: s.methods must be initialized...
	// otherwise it will be nil, and s.methods[method.Name] will fail
	s.methods = make(map[string]*Method)
	for i := 0; i < stp.NumMethod(); i++ {
		method := stp.Method(i)
		if !isMethodValid(method) {
			continue
		}
		s.methods[method.Name] = NewMethod(method)
		log.Println("[INFO]: Register the method ", method.Name)
	}
}

func (s *Service) Call(method *Method, argv, reply reflect.Value) error {
	// note: there are actually 3 arguments in method.Fun.Call
	// the first one is self, which is s.rfvalue,
	// the second and third one is argv and reply
	// we need to pass all of them as a slice []reflect.Value to method.Func.Call
	ret := method.Fun.Call([]reflect.Value{s.rfvalue, argv, reply})
	if len(ret) != 1 {
		log.Fatal("return must be a error, something wrong")
	}

	// cvrt ret[0] to error
	// (nil).(error) is not ok...
	// so firstly we need to check it is nil or not
	// if not nil, then check it is error or not
	if ret[0].Interface() != nil {
		err, ok := ret[0].Interface().(error)
		if !ok {
			log.Fatal("return must be a error or nil, something went wrong here")
		}
		return err
	}

	return nil
}

func (s *Service) GetMethod(name string) (m *Method, err error) {
	m, ok := s.methods[name]
	if !ok {
		return nil, fmt.Errorf("method %s is not exist in service %s", name, s.name)
	}
	return
}

func isMethodValid(method reflect.Method) bool {
	mtp := method.Type
	// method is exportable
	if !method.IsExported() {
		return false
	}
	// (self, argv, *replyv) error
	// check number
	if mtp.NumIn() != 3 || mtp.NumOut() != 1 {
		return false
	}
	// print method type format like: func(*service.AddService, service.ArgType, *int) error
	log.Printf("method type format like: %s\n", mtp.String())

	argt, replyt, rett := mtp.In(1), mtp.In(2), mtp.Out(0)
	// argv and *reply is exportable
	if !isExportedOrBuildin(argt) || !isExportedOrBuildin(replyt) {
		return false
	}
	// *reply
	if replyt.Kind() != reflect.Ptr {
		return false
	}
	// return is error type
	if rett == reflect.TypeOf(errors.New("")) {
		return false
	}

	return true
}

func isExportedOrBuildin(tp reflect.Type) bool {
	// exportable: the first letter of tp.Name() is captalized,
	// 		ast.IsExported(string) can check it
	// build-in: build-in type, like int, string, etc, which has no PkgPath
	return ast.IsExported(tp.Name()) || tp.PkgPath() == ""
}

func (s *Service) GetName() string {
	return s.name
}
