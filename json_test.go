// Copyright (c) 2025 srfrog - https://srfrog.dev
// Use of this source code is governed by the license in the LICENSE file.

package dict

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDictMarshalJSON(t *testing.T) {
	tests := []struct {
		in  interface{}
		out string
	}{
		{in: nil, out: `null`},
		{in: int(1), out: `{"0":1}`},
		{in: float64(2.2), out: `{"0":2.2}`},
		{in: "2.2", out: `{"0":"2.2"}`},
		{in: uint(300), out: `{"0":300}`},
		{in: []int{1, 2, 3}, out: `{"0":1,"1":2,"2":3}`},
		{in: [][]int{{1, 2, 3}}, out: `{"0":[1,2,3]}`},
		{in: map[string]int{"one item": 1}, out: `{"one item":1}`},
	}
	for _, tc := range tests {
		d := New(tc.in)
		b, err := json.Marshal(d)
		require.NoError(t, err)
		require.JSONEq(t, tc.out, string(b))
	}
}

func TestDictMarshalJSON_Embed(t *testing.T) {
	d := New(1, 2, 3)
	d.Set(d.Len(), New(4, 5, 6))

	b, err := json.Marshal(d)
	require.NoError(t, err)
	j := `
	{
		"0":1,
		"1":2,
		"2":3,
		"3":{
			"0":4,
			"1":5,
			"2":6
		}
	}`
	require.JSONEq(t, j, string(b))
}

func TestDictMarshalJSONErr(t *testing.T) {
	d := New().Set("x", func() {})
	_, err := json.Marshal(d)
	require.Error(t, err)
}

func TestDictUnmarshalJSON(t *testing.T) {
	j := `{
			"1": true,
			"2": "two",
			"3": 3.30003,
			"4a": ["horse","cow"],
			"4b": [1, 2, 3],
			"4c": [1.1, 2.2, 3.3],
			"4d": [3, "something", 4.4],
			"4e": [null, null, 0.0001, null],
			"4f": [true, false, null],
			"5": {"horse": "neighs", "cow": "moos", "dog": "woofs"},
			"6": null
		}`
	d := New()
	require.NoError(t, json.Unmarshal([]byte(j), d))

	tests := []struct {
		in  string
		out interface{}
	}{
		{in: "1", out: true},
		{in: "2", out: "two"},
		{in: "3", out: float64(3.30003)},
		{in: "4a", out: []string{"horse", "cow"}},
		{in: "4b", out: []float64{1, 2, 3}},
		{in: "4c", out: []float64{1.1, 2.2, 3.3}},
		{in: "4d", out: []interface{}{float64(3), "something", float64(4.4)}},
		{in: "4e", out: []float64{0, 0, 0.0001, 0}},
		{in: "4f", out: []bool{true, false, false}},
		{in: "6", out: nil},
	}
	for _, tc := range tests {
		require.EqualValues(t, tc.out, d.Get(tc.in))
	}

	// Embedded dict
	ed, ok := d.Get("5").(*Dict)
	require.True(t, ok)
	require.True(t, ed.Len() == 3)
	require.EqualValues(t, "neighs", ed.Get("horse"))
	require.EqualValues(t, "moos", ed.Get("cow"))
	require.EqualValues(t, "woofs", ed.Get("dog"))
}

func TestDictUnmarshalJSONErr(t *testing.T) {
	d := New()
	require.Error(t, json.Unmarshal([]byte(nil), d))
}
