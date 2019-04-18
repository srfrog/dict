// Copyright (c) 2019 srfrog - https://srfrog.me
// Use of this source code is governed by the license in the LICENSE file.

package dict

import "testing"

const N = 1 << 10

func newMap(b *testing.B) map[int]int {
	m := make(map[int]int)
	for i := 0; i < N; i++ {
		m[i] = i
	}

	b.ResetTimer()
	return m
}

func newDict(b *testing.B) *Dict {
	d := New()
	for i := 0; i < N; i++ {
		d.Set(i, i)
	}

	b.ResetTimer()
	return d
}

func BenchmarkGoMapGet(b *testing.B) {
	m := newMap(b)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < N; i++ {
				v := m[i]
				if v != i {
					b.Fail()
				}
			}
		}
	})
}

func BenchmarkDictGet(b *testing.B) {
	d := newDict(b)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < N; i++ {
				v := d.Get(i)
				if v != i {
					b.Fail()
				}
			}
		}
	})
}
