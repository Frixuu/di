package di

import "reflect"

type Container interface {
	Put(t reflect.Type, service Service)
	PutNamed(t reflect.Type, named string, service Service)
	Get(t reflect.Type) (Service, bool)
	GetNamed(t reflect.Type, named string) (Service, bool)
}

type containerImpl struct {
	m typedSyncMap[string, *typedSyncMap[reflect.Type, Service]]
}

func NewContainer() Container {
	return &containerImpl{
		m: typedSyncMap[string, *typedSyncMap[reflect.Type, Service]]{},
	}
}

func (c *containerImpl) mapOfNamed(named string) *typedSyncMap[reflect.Type, Service] {
	return c.m.ComputeIfAbsent(named, func() *typedSyncMap[reflect.Type, Service] {
		return &typedSyncMap[reflect.Type, Service]{}
	})
}

func (c *containerImpl) Put(t reflect.Type, service Service) {
	c.PutNamed(t, "", service)
}

func (c *containerImpl) PutNamed(t reflect.Type, named string, service Service) {
	namedMap := c.mapOfNamed(named)
	namedMap.Put(t, service)
}

func (c *containerImpl) Get(t reflect.Type) (s Service, ok bool) {
	return c.GetNamed(t, "")
}

func (c *containerImpl) GetNamed(t reflect.Type, named string) (s Service, ok bool) {
	namedMap := c.mapOfNamed(named)
	s, ok = namedMap.Get(t)
	return
}
