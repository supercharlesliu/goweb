package goweb

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// StringCast cast string type value to dest type
// Supported dest types: bool, int(int8 ~ int64), uint(uint8 ~ uint64), string, []bool, []int(int8 ~ int64), []uint(uint8 ~ uint64), []string;
func StringCast(src string, typo *reflect.Type) (interface{}, error) {
	var v interface{}
	var err error

	switch (*typo).Kind() {
	case reflect.String:
		v = src
	case reflect.Bool:
		v = src == "true"
	case reflect.Int:
		v, err = strconv.Atoi(src)
	case reflect.Int8:
		v, err = strconv.ParseInt(src, 10, 8)
	case reflect.Int16:
		v, err = strconv.ParseInt(src, 10, 16)
	case reflect.Int32:
		v, err = strconv.ParseInt(src, 10, 32)
	case reflect.Int64:
		v, err = strconv.ParseInt(src, 10, 64)
	case reflect.Uint:
		var vOfU64 uint64
		if vOfU64, err = strconv.ParseUint(src, 10, 0); err == nil {
			v = uint(vOfU64)
		}
	case reflect.Uint8:
		v, err = strconv.ParseUint(src, 10, 8)
	case reflect.Uint16:
		v, err = strconv.ParseUint(src, 10, 16)
	case reflect.Uint32:
		v, err = strconv.ParseUint(src, 10, 32)
	case reflect.Uint64:
		v, err = strconv.ParseUint(src, 10, 64)
	case reflect.Slice:
		list := reflect.MakeSlice(*typo, 0, 0)
		if src == "" {
			return list.Interface(), err
		}

		parts := strings.Split(src, ",")
		elementType := (*typo).Elem()
		if elementType.Kind() == reflect.String {
			return parts, nil
		}

		for _, itemV := range parts {
			vCtv, err := StringCast(itemV, &elementType)
			if err != nil {
				return nil, err
			}
			list = reflect.AppendSlice(list, reflect.ValueOf(vCtv))
		}
		v = list.Interface()
	default:
		return nil, errors.New("Unsupported type cast: " + (*typo).Kind().String())
	}

	if err != nil {
		return nil, err
	}

	return v, nil
}
