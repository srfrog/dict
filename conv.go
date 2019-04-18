// Copyright (c) 2019 srfrog - https://srfrog.me
// Use of this source code is governed by the license in the LICENSE file.

package dict

import (
	"reflect"
	"strconv"
)

// Stringer is just like fmt.Stringer without loading that package.
type Stringer interface {
	String() string
}

var stringerType = reflect.TypeOf((*Stringer)(nil)).Elem()

func toFloat64(x interface{}) float64 {
	if v, ok := x.(float32); ok {
		return float64(v)
	}
	return x.(float64)
}

func toUint64(x interface{}) uint64 {
	switch v := x.(type) {
	case uint:
		return uint64(v)
	case uint8:
		return uint64(v)
	case uint16:
		return uint64(v)
	case uint32:
		return uint64(v)
	case uint64:
		return v
	}
	return x.(uint64)
}

func toInt64(x interface{}) int64 {
	switch v := x.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	}
	return x.(int64)
}

func toString(x interface{}) string {
	var s string
	switch v := x.(type) {
	case float32, float64:
		s = strconv.FormatFloat(toFloat64(v), 'f', -1, 64)
	case int, int8, int16, int32, int64:
		s = strconv.FormatInt(toInt64(v), 10)
	case uint, uint8, uint16, uint32, uint64:
		s = strconv.FormatUint(toUint64(v), 10)
	case string:
		s = v
	case Stringer:
		s = v.String()
	}
	return s
}

// Item is a key-value pair.
// Key is the key name value.
// Value is the stored value in dict.
type Item struct {
	Key   interface{}
	Value interface{}
}

func toIterable(i interface{}) <-chan Item {
	ci := make(chan Item)

	go func() {
		defer close(ci)

		if item, ok := i.(Item); ok {
			ci <- item
			return
		}

		t := reflect.TypeOf(i)
		if t == nil {
			return
		}

		switch v := reflect.ValueOf(i); t.Kind() {
		case reflect.Map:
			// The map key must be a hashable key type.
			if !isKeyType(t.Key()) {
				break
			}
			for iter := v.MapRange(); iter.Next(); {
				ci <- Item{Key: iter.Key().Interface(), Value: iter.Value().Interface()}
			}

		case reflect.Chan:
		L:
			for {
				x, ok := v.Recv()
				if !ok {
					break L
				}
				ci <- Item{Value: x.Interface()}
			}

		case reflect.Array, reflect.Slice:
			for j := 0; j < v.Len(); j++ {
				ci <- Item{Key: j, Value: v.Index(j).Interface()}
			}

		default:
			ci <- Item{Value: v.Interface()}
		}
	}()

	return ci
}

func isKeyType(t reflect.Type) bool {
	kind := t.Kind()
	return (kind > reflect.Bool && kind < reflect.Uintptr) ||
		kind == reflect.Float32 || kind == reflect.Float64 ||
		kind == reflect.String ||
		t.Implements(stringerType)
}
