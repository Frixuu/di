package di

import (
	"reflect"
)

func getType[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
