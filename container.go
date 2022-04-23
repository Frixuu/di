package di

type Container interface {
	Put(typename string, service Service)
	PutNamed(typename string, named string, service Service)
	Get(typename string) (Service, bool)
	GetNamed(typename string, named string) (Service, bool)
}

type containerImpl struct {
	m typedSyncMap[string, *typedSyncMap[string, Service]]
}

func NewContainer() Container {
	return &containerImpl{
		m: typedSyncMap[string, *typedSyncMap[string, Service]]{},
	}
}

func (c *containerImpl) mapOfNamed(named string) *typedSyncMap[string, Service] {
	return c.m.ComputeIfAbsent(named, func() *typedSyncMap[string, Service] {
		return &typedSyncMap[string, Service]{}
	})
}

func (c *containerImpl) Put(typename string, service Service) {
	c.PutNamed(typename, "", service)
}

func (c *containerImpl) PutNamed(typename string, named string, service Service) {
	namedMap := c.mapOfNamed(named)
	namedMap.Put(typename, service)
}

func (c *containerImpl) Get(typename string) (s Service, ok bool) {
	return c.GetNamed(typename, "")
}

func (c *containerImpl) GetNamed(typename string, named string) (s Service, ok bool) {
	namedMap := c.mapOfNamed(named)
	s, ok = namedMap.Get(typename)
	return
}
