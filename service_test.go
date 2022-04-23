package di

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	c := NewContainer()
	Register[a, struct{}](c)

	s := &SingletonService{
		ImplType: reflect.TypeOf(testImpl{}),
		IsBuilt:  false,
	}

	tInstance := s.Build(c)
	assert.NotNil(t, tInstance)
	assert.IsType(t, &testImpl{}, tInstance.Interface())

	impl := tInstance.Interface().(*testImpl)
	assert.NotNil(t, impl.A)

	assert.Equal(t, tInstance, s.Instance)
	assert.True(t, s.IsBuilt)

	tInstance2 := s.Build(c)
	assert.Equal(t, tInstance, tInstance2)
}

func TestTransient(t *testing.T) {
	c := NewContainer()
	Register[a, struct{}](c)

	s := &TransientService{
		ImplType: reflect.TypeOf(testImpl{}),
	}

	tInstance := s.Build(c)
	assert.NotNil(t, tInstance)
	assert.IsType(t, &testImpl{}, tInstance.Interface())

	impl := tInstance.Interface().(*testImpl)
	assert.NotNil(t, impl.A)

	tInstance2 := s.Build(c)
	assert.NotEqual(t, tInstance, tInstance2)
}
