package goweb

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCastUint(t *testing.T) {
	var vOfUint uint
	typo := reflect.ValueOf(vOfUint).Type()
	v, err := StringCast("11", &typo)
	if err != nil {
		panic(err)
	}
	fmt.Println(reflect.ValueOf(v).Type().Kind() == reflect.Uint)
}
