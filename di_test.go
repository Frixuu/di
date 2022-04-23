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
	const key = "github.com/zekrotja/di.myInterface"

	err = Register[myInterface, testImpl](c)
	assert.Nil(t, err)

	v, ok := c.(*containerImpl).m.Load(key)
	assert.True(t, ok, "no service has been registered")

	svc := v.(*SingletonService)
	assert.Equal(t, svc.ImplType, reflect.TypeOf(testImpl{}))
	assert.False(t, svc.IsBuilt)
	assert.Equal(t, svc.Instance, reflect.Value{})
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
