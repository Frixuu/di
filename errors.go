package di

import "errors"

var (
	ErrInvalidServiceType    = errors.New("the passed type parameter does not represent an interface nor a pointer to a struct")
	ErrDoesNotImplInterface  = errors.New("the passed service instance does not implement the specified interface")
	ErrNotRegistered         = errors.New("no service is registered for the interface requested")
	ErrInvalidImplementation = errors.New("the retrieved service does not implement the requested interface")
)
