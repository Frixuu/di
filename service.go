package di

import (
	"reflect"
)

type Service interface {
	Build(c Container) reflect.Value
}

type SingletonService struct {
	ImplType reflect.Type
	IsBuilt  bool
	Instance reflect.Value
}

func (s *SingletonService) Build(c Container) (instance reflect.Value) {
	if s.IsBuilt {
		instance = s.Instance
		return
	}
	instance = reflect.New(s.ImplType)
	s.Instance = instance
	s.IsBuilt = true
	populateService(&instance, c, s.ImplType)
	return
}

type TransientService struct {
	ImplType reflect.Type
}

func (s *TransientService) Build(c Container) (instance reflect.Value) {
	instance = reflect.New(s.ImplType)
	populateService(&instance, c, s.ImplType)
	return
}

func populateService(s *reflect.Value, c Container, t reflect.Type) {
	elem := s.Elem()
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		if field.Kind() != reflect.Interface {
			continue
		}

		if !field.CanSet() || !field.IsNil() {
			continue
		}

		var svc Service
		meta := getTagMetaMap(t.Field(i), "di")
		name, ok := meta["named"]
		if ok {
			svc, ok = c.GetNamed(field.Type(), name)
		} else {
			svc, ok = c.Get(field.Type())
		}

		if ok {
			field.Set(svc.Build(c))
		}
	}
}
