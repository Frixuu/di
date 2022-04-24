package di

import (
	"reflect"
	"strings"
)

// getTagMetaMap parses a struct tag to a map of values.
//
// It assumes that the pairs in the tag are separated by a comma
// and the key is separated from the value by a colon.
func getTagMetaMap(field reflect.StructField, key string) map[string]string {
	m := make(map[string]string)
	tag := field.Tag.Get(key)
	tagPairs := strings.Split(tag, ",")
	for _, pair := range tagPairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
