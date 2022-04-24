package di

import (
	"reflect"
)

func Register[TSvc, TImpl any](c Container) error {
	return RegisterNamed[TSvc, TImpl](c, "")
}

func MustRegister[TSvc, TImpl any](c Container) {
	must(Register[TSvc, TImpl](c))
}

func RegisterNamed[TSvc, TImpl any](c Container, name string) (err error) {
	tImpl := getType[TImpl]()
	tIf := getType[TSvc]()
	if err = areTypesValidForDi(tIf, tImpl, false); err != nil {
		return
	}
	c.PutNamed(tIf, name, &SingletonService{
		ImplType: tImpl,
	})
	return
}

func MustRegisterNamed[TSvc, TImpl any](c Container, name string) {
	must(RegisterNamed[TSvc, TImpl](c, name))
}

func RegisterTransient[TSvc, TImpl any](c Container) error {
	return RegisterTransientNamed[TSvc, TImpl](c, "")
}

func MustRegisterTransient[TSvc, TImpl any](c Container) {
	must(RegisterTransient[TSvc, TImpl](c))
}

func RegisterTransientNamed[TSvc, TImpl any](c Container, name string) (err error) {
	tImpl := getType[TImpl]()
	tIf := getType[TSvc]()
	if err = areTypesValidForDi(tIf, tImpl, false); err != nil {
		return
	}
	c.PutNamed(tIf, name, &TransientService{
		ImplType: tImpl,
	})
	return
}

func MustRegisterTransientNamed[TSvc, TImpl any](c Container, name string) {
	must(RegisterTransientNamed[TSvc, TImpl](c, name))
}

func RegisterInstance[TSvc, TImpl any](c Container, i TImpl) error {
	return RegisterInstanceNamed[TSvc](c, "", i)
}

func MustRegisterInstance[TSvc, TImpl any](c Container, i TImpl) {
	must(RegisterInstance[TSvc](c, i))
}

func RegisterInstanceNamed[TSvc, TImpl any](c Container, name string, i TImpl) (err error) {
	tImpl := getType[TImpl]()
	tIf := getType[TSvc]()
	if err = areTypesValidForDi(tIf, tImpl, true); err != nil {
		return
	}
	c.PutNamed(tIf, name, &SingletonService{
		ImplType: tImpl,
		IsBuilt:  true,
		Instance: reflect.ValueOf(i),
	})
	return
}

func MustRegisterInstanceNamed[TSvc, TImpl any](c Container, name string, i TImpl) {
	must(RegisterInstanceNamed[TSvc](c, name, i))
}

func Get[T any](c Container) (s T, err error) {
	s, err = GetNamed[T](c, "")
	return
}

func MustGet[T any](c Container) (s T) {
	s, err := Get[T](c)
	must(err)
	return
}

func GetNamed[T any](c Container, name string) (s T, err error) {
	tIf := getType[T]()
	err = checkServiceType(tIf)
	if err != nil {
		return
	}
	sb, ok := c.GetNamed(tIf, name)
	if !ok {
		err = ErrNotRegistered
		return
	}
	v := sb.Build(c).Interface()
	s, ok = v.(T)
	if !ok {
		err = ErrInvalidImplementation
	}
	return
}

func MustGetNamed[T any](c Container, name string) (s T) {
	s, err := GetNamed[T](c, name)
	must(err)
	return
}

// checkServiceType checks if a type can be registered as a service type.
func checkServiceType(tIf reflect.Type) error {
	if isInterface(tIf) {
		return nil
	}
	if isPointerToStruct(tIf) {
		return nil
	}
	return ErrInvalidServiceType
}

// areTypesValidForDi checks if a pair of an abstract and concrete type are valid for injection.
func areTypesValidForDi(tIf reflect.Type, tImpl reflect.Type, allowSamePtr bool) error {
	if err := checkServiceType(tIf); err != nil {
		return err
	}
	isIf := isInterface(tIf)
	isImpl := isIf && (tImpl.Implements(tIf) || reflect.PointerTo(tImpl).Implements(tIf))
	if isImpl {
		return nil
	}
	isPtr := tIf == reflect.PointerTo(tImpl) || (allowSamePtr && tIf == tImpl)
	if isPtr {
		return nil
	}
	return ErrDoesNotImplInterface
}
