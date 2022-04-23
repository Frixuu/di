package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagValidGetsParsed(t *testing.T) {
	type myStruct struct {
		foo int `x:"a:1,b:2" y:"c:3,d:4,e:5" z:"f:6"`
	}

	ty := getType[myStruct]()
	field, _ := ty.FieldByName("foo")

	x := getTagMetaMap(field, "x")
	assert.Equal(t, 2, len(x))
	assert.Equal(t, "1", x["a"])
	assert.Equal(t, "2", x["b"])

	y := getTagMetaMap(field, "y")
	assert.Equal(t, 3, len(y))

	z := getTagMetaMap(field, "z")
	assert.Equal(t, 1, len(z))

	m := getTagMetaMap(field, "m")
	assert.Equal(t, 0, len(m))
}
