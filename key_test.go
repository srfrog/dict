// Copyright (c) 2025 srfrog - https://srfrog.dev
// Use of this source code is governed by the license in the LICENSE file.

package dict_test

import (
	"testing"

	"github.com/srfrog/dict"
	"github.com/stretchr/testify/require"
)

func TestMakeKey(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{"nil key value", func(t *testing.T) {
			key := dict.MakeKey(nil)
			require.Nil(t, key)
		}},
		{"invalid key type", func(t *testing.T) {
			var value = struct{}{}
			key := dict.MakeKey(value)
			require.Nil(t, key)
		}},
		{"valid key type", func(t *testing.T) {
			var value interface{} = 1
			key := dict.MakeKey(value)
			require.NotNil(t, key)
			require.Equal(t, key.Name, "1")
		}},
		{"equal key Ids", func(t *testing.T) {
			value1 := 123
			key1 := dict.MakeKey(value1)
			require.NotNil(t, key1)

			value2 := "123"
			key2 := dict.MakeKey(value2)
			require.NotNil(t, key2)

			require.Equal(t, key1.ID, key2.ID)
			require.Equal(t, key1.Name, key2.Name)
		}},
	}
	for _, tc := range tests {
		t.Run(tc.name, tc.fn)
	}
}
