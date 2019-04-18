// Copyright (c) 2019 srfrog - https://srfrog.me
// Use of this source code is governed by the license in the LICENSE file.

package dict

import (
	"hash/fnv"
)

// Key represents a key value. Keys are used to order the items in a dict.
// ID is a 64 bit hash value representation of Name.
// Name is the user-friendly and sortable name.
type Key struct {
	ID   uint64
	Name string
}

func isValidKeyType(t interface{}) bool {
	switch t.(type) {
	case float32, float64:
		return true
	case int8, int16, int32, int64, int:
		return true
	case uint8, uint16, uint32, uint64, uint:
		return true
	case string:
		return true
	case Stringer:
		return true
	}
	return false
}

// MakeKey generates a Key object by hashing the provided value. The value type must be float,
// int, uint, string, or that implements Stringer.
// Returns a new Key object if successful, otherwise returns nil.
func MakeKey(value interface{}) *Key {
	var name string

	if !isValidKeyType(value) {
		return nil
	}

	name = toString(value)
	if name == "" {
		return nil
	}

	h := fnv.New64a()
	_, err := h.Write([]byte(name))
	if err != nil {
		return nil
	}

	return &Key{
		ID:   h.Sum64(),
		Name: name,
	}
}
