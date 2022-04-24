package di

import "reflect"

// getType returns a reflect.Type of a provided generic type.
func getType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

// isInterface checks if a reflected type is an interface.
func isInterface(t reflect.Type) bool {
	return t.Kind() == reflect.Interface
}

// isInterfaceV checks if a reflected value is an interface.
func isInterfaceV(v reflect.Value) bool {
	return v.Kind() == reflect.Interface
}

// isPointerToStruct checks if a reflected type is a pointer to a struct.
func isPointerToStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct
}

// isPointerToStructV checks if a reflected value is a pointer to a struct.
func isPointerToStructV(v reflect.Value) bool {
	return isPointerToStruct(v.Type())
}
