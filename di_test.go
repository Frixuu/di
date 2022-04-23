package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type a interface{}

type testImpl struct {
	A a
}

func TestRegister(t *testing.T) {
	c := NewContainer()

	// Infers testImpl as type instead
	// of an interface.
	err := Register[testImpl, testImpl](c)
	assert.ErrorIs(t, err, ErrNoInterface)

	type myInterface interface{}
	key := getType[myInterface]()

	err = Register[myInterface, testImpl](c)
	assert.Nil(t, err)

	v, ok := c.(*containerImpl).Get(key)
	assert.True(t, ok, "no service has been registered")

	svc := v.(*SingletonService)
	assert.Equal(t, svc.ImplType, reflect.TypeOf(testImpl{}))
	assert.False(t, svc.IsBuilt)
	assert.Equal(t, svc.Instance, reflect.Value{})
}

func TestNamed(t *testing.T) {
	c := NewContainer()

	type (
		myInterface  interface{}
		myImplGlobal struct{}
		myImpl1      struct{}
		myImpl2      struct{}
		myImpl3      struct{}
	)

	key := getType[myInterface]()
	ig := &myImplGlobal{}
	i1 := &myImpl1{}
	i2 := &myImpl2{}
	i3 := &myImpl3{}

	c.PutNamed(key, "svc-1", &SingletonService{
		IsBuilt:  true,
		Instance: reflect.ValueOf(i1),
	})

	c.PutNamed(key, "svc-2", &SingletonService{
		IsBuilt:  true,
		Instance: reflect.ValueOf(i2),
	})

	c.Put(key, &SingletonService{
		IsBuilt:  true,
		Instance: reflect.ValueOf(ig),
	})

	c.PutNamed(key, "svc-3", &SingletonService{
		IsBuilt:  true,
		Instance: reflect.ValueOf(i3),
	})

	sg, _ := c.(*containerImpl).Get(key)
	assert.Same(t, ig, sg.Build(c).Interface())

	s1, _ := c.(*containerImpl).GetNamed(key, "svc-1")
	assert.Same(t, i1, s1.Build(c).Interface())

	s2, _ := c.(*containerImpl).GetNamed(key, "svc-2")
	assert.Same(t, i2, s2.Build(c).Interface())

	s3, _ := c.(*containerImpl).GetNamed(key, "svc-3")
	assert.Same(t, i3, s3.Build(c).Interface())

	sg2, _ := c.(*containerImpl).Get(key)
	assert.Same(t, ig, sg2.Build(c).Interface())
}

func TestNamedFromTag(t *testing.T) {
	type (
		svc           interface{}
		svcOne        struct{}
		svcTwo        struct{}
		controller    interface{}
		controllerOne struct {
			Service svc `di:"named:one"`
		}
		controllerTwo struct {
			Service svc `di:"named:two"`
		}
	)

	c := NewContainer()
	MustRegisterNamed[svc, svcOne](c, "one")
	MustRegisterNamed[controller, controllerOne](c, "one")
	MustRegisterNamed[svc, svcTwo](c, "two")
	MustRegisterNamed[controller, controllerTwo](c, "two")

	s1 := MustGetNamed[svc](c, "one")
	assert.NotNil(t, s1)
	c1 := MustGetNamed[controller](c, "one")
	assert.NotNil(t, c1)
	assert.IsType(t, &controllerOne{}, c1)
	cc1, ok := c1.(*controllerOne)
	assert.True(t, ok)
	assert.Same(t, s1, cc1.Service)

	s2 := MustGetNamed[svc](c, "two")
	assert.NotNil(t, s2)
	assert.NotSame(t, s1, s2)
	c2 := MustGetNamed[controller](c, "two")
	assert.NotNil(t, c2)
	assert.IsType(t, &controllerTwo{}, c2)
	cc2, ok := c2.(*controllerTwo)
	assert.True(t, ok)
	assert.Same(t, s2, cc2.Service)

	assert.Panics(t, func() {
		MustGet[svc](c)
	})
}

func TestGet(t *testing.T) {
	c := NewContainer()

	type myInterface interface{}
	type myOtherInterface interface{}

	impl := testImpl{}
	Register[myInterface, testImpl](c)

	s, err := Get[myInterface](c)
	assert.Nil(t, err)
	assert.IsType(t, s, &impl)

	// Ensure that the retrieved value is
	// exactly the same instance on an
	// repeated retrieve.
	s2, err := Get[myInterface](c)
	assert.Nil(t, err)
	assert.Same(t, s2, s)

	_, err = Get[myOtherInterface](c)
	assert.ErrorIs(t, err, ErrNotRegistered)

	_, err = Get[struct{}](c)
	assert.ErrorIs(t, err, ErrNoInterface)
}

func TestTransientGetsDifferent(t *testing.T) {
	c := NewContainer()

	type myInterface interface{}

	impl := testImpl{}
	RegisterTransient[myInterface, testImpl](c)

	s, err := Get[myInterface](c)
	assert.Nil(t, err)
	assert.IsType(t, s, &impl)

	// This time, we want returned instances to be different
	s2, err := Get[myInterface](c)
	assert.Nil(t, err)
	assert.NotSame(t, s2, s)
}

func TestInstance(t *testing.T) {
	c := NewContainer()

	type (
		myInterface interface{}
		myType      struct {
			foo int
		}
	)

	instance := &myType{
		foo: 4,
	}
	MustRegisterInstance[myInterface](c, instance)

	s, err := Get[myInterface](c)
	assert.Nil(t, err)
	assert.IsType(t, instance, s)
	assert.Equal(t, 4, (s.(*myType)).foo)

	instance.foo = 5
	assert.Equal(t, 5, (s.(*myType)).foo)
}

func TestCrossDependency(t *testing.T) {
	type (
		S1 interface{}
		S2 interface{}

		S1Impl struct {
			S S2
		}
		S2Impl struct {
			S S1
		}
	)

	c := NewContainer()

	assert.Nil(t, Register[S1, S1Impl](c))
	assert.Nil(t, Register[S2, S2Impl](c))

	s2, err := Get[S2](c)
	assert.Nil(t, err)

	s2i := s2.(*S2Impl)
	assert.NotNil(t, s2i.S)

	s1i := s2i.S.(*S1Impl)
	assert.NotNil(t, s1i.S)
}

func TestNoInterface(t *testing.T) {
	type (
		S1 interface{}

		S1Impl struct {
			S struct{}
			s S1
		}
	)

	c := NewContainer()

	assert.Nil(t, Register[S1, S1Impl](c))
	_, err := Get[S1](c)
	assert.Nil(t, err)

}
