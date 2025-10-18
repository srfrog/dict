// Copyright (c) 2019 srfrog - https://srfrog.dev
// Use of this source code is governed by the license in the LICENSE file.

package dict

import (
	"encoding/json"
	"reflect"
	"strings"
)

// MarshalJSON implements the json.MarshalJSON interface.
// The JSON representation of dict is just a JSON object.
func (d *Dict) MarshalJSON() ([]byte, error) {
	if d.IsEmpty() {
		return []byte("null"), nil
	}

	var (
		err error
		sb  strings.Builder
		cnt int
	)

	sb.WriteByte('{')
	for item := range d.Items() {
		var p []byte

		sb.WriteByte('"')
		sb.WriteString(item.Key.(string))
		sb.WriteByte('"')
		sb.WriteByte(':')

		p, err = json.Marshal(item.Value)
		if err != nil {
			return nil, err
		}
		sb.Write(p)
		cnt++
		if cnt < d.Len() {
			sb.WriteByte(',')
		}
	}
	sb.WriteByte('}')

	return []byte(sb.String()), nil
}

// UnmarshalJSON implements the json.UnmarshalJSON interface.
// The JSON representation of dict is just a JSON object.
func (d *Dict) UnmarshalJSON(p []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(p, &m); err != nil {
		return err
	}

	// Unforunately json.Unmarshal will produce dynamic interface types for JSON arrays
	// and objects - https://golang.org/pkg/encoding/json/#Unmarshal
	// So here we try to convert []interface{} (JSON array) values into a slice if all the
	// value types are the same. e.g., []string, []float64, etc...
	// Also convert map[string]interface{} (JSON object) values into embedded dict objects.
	for k, v := range m {
		switch x := v.(type) {
		// JSON array -> slice
		case []interface{}:
			kind, ok := hasSameKind(x)
			if !ok {
				break
			}
			switch kind {
			case reflect.Bool:
				var bs []bool
				for i := range x {
					bv, _ := x[i].(bool)
					bs = append(bs, bv)
				}
				m[k] = bs
			case reflect.Float64:
				var fs []float64
				for i := range x {
					fv, _ := x[i].(float64)
					fs = append(fs, fv)
				}
				m[k] = fs
			case reflect.String:
				var ss []string
				for i := range x {
					sv, _ := x[i].(string)
					ss = append(ss, sv)
				}
				m[k] = ss
			}

		// JSON object -> dict
		case map[string]interface{}:
			m[k] = New(x)
		}
	}
	d.Update(m)

	return nil
}

func hasSameKind(a []interface{}) (reflect.Kind, bool) {
	var k, kseen reflect.Kind
	for i := range a {
		switch a[i].(type) {
		case nil:
			// If at least one value isn't nil (JSON null) convert it to the zero value of
			// the type.
		case bool:
			k = reflect.Bool
		case float64:
			k = reflect.Float64
		case string:
			k = reflect.String
		default:
			// TODO: Array of arrays and array of objects.
			return reflect.Invalid, false
		}
		if kseen == 0 {
			kseen = k
			continue
		}
		if k != kseen {
			return reflect.Invalid, false
		}
	}
	return kseen, kseen != reflect.Invalid
}
