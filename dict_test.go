// Copyright (c) 2025 srfrog - https://srfrog.dev
// Use of this source code is governed by the license in the LICENSE file.

package dict

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDict(t *testing.T) {
	d := &Dict{values: make(map[uint64]interface{})}
	wrapFn := func(fn func(*testing.T, *Dict)) func(*testing.T) {
		return func(t *testing.T) { fn(t, d) }
	}

	tests := []struct {
		name string
		fn   func(t *testing.T, d *Dict)
	}{
		{"new", testNew},
		{"new chan", testNewChan},
		{"set", testSet},
		{"get", testGet},
		{"set update", testSetUpdate},
		{"insert", testSetInsert},
		{"set embed", testSetEmbed},
		{"set chan", testSetChan},
		{"key", testKey},
		{"isempty", testIsEmpty},
		{"del", testDel},
		{"keys", testKeys},
		{"values", testValues},
		{"clear", testClear},
		{"string", testPrint},
		{"update", testUpdate},
		{"workers", testWorkers},
	}
	for _, tc := range tests {
		if !t.Run(tc.name, wrapFn(tc.fn)) {
			break
		}
	}
}

func testNew(t *testing.T, d *Dict) {
	tests := []struct {
		in  interface{}
		len int
	}{
		{in: nil, len: 0},
		{in: 1, len: 1},
		{in: []interface{}{}, len: 0},
		{in: []struct{}{}, len: 0},
		{in: []int{}, len: 0},
		{in: []int{1, 2, 3, 4, 5}, len: 5},
		{in: []uint{1, 2, 3, 4, 5}, len: 5},
		{in: []float64{1, 2, 3, 4, 5}, len: 5},
		{in: []string{}, len: 0},
		{in: []string{"1", "2", "3"}, len: 3},
		{in: map[float32]struct{}{}, len: 0},
		{in: map[int]int{1: 11, 2: 22, 3: 33}, len: 3},
		{in: map[string]int{"1": 1, "2": 2, "3": 3}, len: 3},
		{in: map[int]string{1: "one", 2: "two", 3: "three"}, len: 3},
	}
	for _, tc := range tests {
		d := New(tc.in)
		require.Equal(t, tc.len, d.Len())
	}
}

func testNewChan(t *testing.T, d *Dict) {
	ch1 := make(chan int)
	go func() {
		defer close(ch1)
		for i := 0; i < 100; i++ {
			ch1 <- i
		}
	}()

	t.Run("blocking channel", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, 100, New(ch1).Len())
	})

	ch2 := make(chan int, 100)
	go func() {
		defer close(ch2)
		for i := 0; i < 100; i++ {
			ch2 <- i
		}
	}()

	t.Run("buffered channel", func(t *testing.T) {
		t.Parallel()
		require.Equal(t, 100, New(ch2).Len())
	})
}

type testDevice int

func (dev testDevice) String() string { return fmt.Sprintf("%#x", int(dev)) }

func testSet(t *testing.T, d *Dict) {
	data := []struct {
		key, value interface{}
	}{
		{key: "5DFD011F-6123-4C4D-8BBF-9C26B4D1AD0F", value: testDevice(0x1)},
		{key: "0E22688F-7E76-4F41-9351-243DD0824428", value: testDevice(0x2)},
		{key: "B3FCB096-C0AF-42BB-9AFA-BBAA9CDA1CBC", value: "9872931"},
		{key: "B3FCB096-C0AF-42BB-9AFA-BBAA9CDA1CBC", value: testDevice(0x3)},
		{key: "F1A6671B-0E1D-4D80-9182-19667418F9C2", value: testDevice(0x4)},
		{key: "5F0D34B3-24A5-4F76-BBC5-BF29E97C15AD", value: testDevice(0x5)},
		{key: testDevice(0x1), value: "device 1"},
		{key: testDevice(0x2), value: "device 2"},
		{key: testDevice(0x3), value: "device 3"},
		{key: testDevice(0x4), value: "device 4"},
		{key: testDevice(0x5), value: "some device 5"},
		{key: testDevice(0x5), value: "device 5"},
	}
	for i := range data {
		out := d.Set(data[i].key, data[i].value)
		require.EqualValues(t, d, out)
		require.False(t, out.size == 0, "expected size > 10 but got %d", out.size)
	}
	require.True(t, d.Len() == 10, "expected dict length to be 10 but got %d", d.Len())
}

func testGet(t *testing.T, d *Dict) {
	tests := []struct {
		in, out, alt interface{}
	}{
		{in: nil, out: nil},
		{in: "5DFD011F-6123-4C4D-8BBF-9C26B4D1AD0F", out: testDevice(0x1)},
		{in: "0E22688F-7E76-4F41-9351-243DD0824428", out: testDevice(0x2)},
		{in: "B3FCB096-C0AF-42BB-9AFA-BBAA9CDA1CBC", out: testDevice(0x3)},
		{in: "F1A6671B-0E1D-4D80-9182-19667418F9C2", out: testDevice(0x4)},
		{in: "5F0D34B3-24A5-4F76-BBC5-BF29E97C15AD", out: testDevice(0x5)},
		{in: testDevice(0x1), out: "device 1"},
		{in: testDevice(0x2), out: "device 2"},
		{in: testDevice(0x3), out: "device 3"},
		{in: testDevice(0x4), out: "device 4"},
		{in: testDevice(0x5), out: "device 5"},
		{alt: nil, out: nil},
		{alt: uint64(0), out: uint64(0)},
		{alt: "default value", out: "default value"},
	}
	for _, tc := range tests {
		out := d.Get(tc.in, tc.alt)
		require.EqualValues(t, tc.out, out)
	}
}

func testSetUpdate(t *testing.T, d *Dict) {
	tests := []struct {
		in  Item
		out string
	}{
		{in: Item{Key: testDevice(0x1), Value: "new device 1"}, out: "new device 1"},
		{in: Item{Key: testDevice(0x2), Value: "new device 2"}, out: "new device 2"},
		{in: Item{Key: testDevice(0x3), Value: "new device 3"}, out: "new device 3"},
		{in: Item{Key: testDevice(0x4), Value: "new device 4"}, out: "new device 4"},
		{in: Item{Key: testDevice(0x5), Value: "new device 5"}, out: "new device 5"},
	}
	size := d.Len()
	for _, tc := range tests {
		dd := d.Set(tc.in.Key, tc.in.Value)
		require.EqualValues(t, d, dd)
		require.Equal(t, size, dd.Len())
		out := d.Get(tc.in.Key)
		require.Equal(t, tc.out, out)
	}
}

type testEmployee uint

func (dev testEmployee) String() string { return fmt.Sprintf("%#x", uint(dev)) }

func testSetInsert(t *testing.T, d *Dict) {
	data := []struct {
		key, value interface{}
	}{
		{key: testEmployee(0xa), value: New().Set("name", "Stephanie Alexander")},
		{key: testEmployee(0xb), value: New().Set("name", "Audrey Bishop")},
		{key: testEmployee(0xc), value: New().Set("name", "Jack Fields")},
		{key: testEmployee(0xd), value: New().Set("name", "Antonio Lambert")},
		{key: testEmployee(0xe), value: New().Set("name", "Tim Hicks")},
	}
	len := d.Len()
	for i := range data {
		out := d.Set(data[i].key, data[i].value)
		require.EqualValues(t, d, out)
		require.False(t, out.Len() == len, "expected size > %d but got %d", len, out.Len())
	}
	require.True(t, d.Len() == 15, "expected dict length to be 15 but got %d", d.Len())
}

func testSetEmbed(t *testing.T, d *Dict) {
	tests := []struct {
		key, value interface{}
	}{
		{key: testEmployee(0xa), value: float64(50000.0)},
		{key: testEmployee(0xb), value: float64(246000.0)},
		{key: testEmployee(0xc), value: float64(115023.0)},
		{key: testEmployee(0xd), value: float64(0)},
		{key: testEmployee(0xe), value: float64(23768.53)},
	}
	len := d.Len()
	for _, tc := range tests {
		rec := d.Get(tc.key, float64(0))
		require.NotNil(t, rec)
		require.NotNil(t, rec.(*Dict).Set("salary", tc.value))
		require.True(t, d.Len() == len, "expected size > %d but got %d", len, d.Len())
	}
	require.True(t, d.Len() == 15, "expected dict length to be 15 but got %d", d.Len())
}

func testSetChan(t *testing.T, d *Dict) {}

func testKey(t *testing.T, d *Dict) {
	tests := []struct {
		in  interface{}
		out bool
	}{
		{in: nil, out: false},
		{in: "", out: false},
		{in: struct{}{}, out: false},
		{in: "5DFD011F-6123-4C4D-8BBF-9C26B4D1AD0F", out: true},
		{in: "0E22688F-7E76-4F41-9351-243DD0824428", out: true},
		{in: "B3FCB096-C0AF-42BB-9AFA-BBAA9CDA1CBC", out: true},
		{in: "F1A6671B-0E1D-4D80-9182-19667418F9C2", out: true},
		{in: "5F0D34B3-24A5-4F76-BBC5-BF29E97C15AD", out: true},
		{in: "4C091301-FA37-4317-94D9-434C9D675B91", out: false},
		{in: "EF2BBB14-719A-4CFF-B1A0-62C79750FCD2", out: false},
		{in: "23F305C3-9D2C-49E0-A108-A8B12FC6534E", out: false},
		{in: "4FD676C8-2FC4-4BC9-A2B7-C78BF2CF21DC", out: false},
		{in: "7249B312-1CD5-4142-81EB-DB438AD45C54", out: false},
		{in: "0x3", out: true},
		{in: testDevice(0x3), out: true},
		{in: "0xc", out: true},
		{in: testEmployee(0xc), out: true},
		{in: float32(math.MaxFloat32), out: false},
		{in: float64(math.MaxFloat64), out: false},
		{in: int8(math.MaxInt8), out: false},
		{in: int16(math.MaxInt16), out: false},
		{in: int32(math.MaxInt32), out: false},
		{in: int64(math.MaxInt64), out: false},
		{in: uint8(math.MaxUint8), out: false},
		{in: uint16(math.MaxUint16), out: false},
		{in: uint32(math.MaxUint32), out: false},
		{in: uint64(math.MaxUint64), out: false},
	}
	for _, tc := range tests {
		out := d.Key(tc.in)
		require.Equal(t, tc.out, out)
	}
}

func testDel(t *testing.T, d *Dict) {
	tests := []struct {
		in  interface{}
		out bool
	}{
		{in: "5DFD011F-6123-4C4D-8BBF-9C26B4D1AD0F", out: true},
		{in: "0E22688F-7E76-4F41-9351-243DD0824428", out: true},
		{in: "B3FCB096-C0AF-42BB-9AFA-BBAA9CDA1CBC", out: true},
		{in: "F1A6671B-0E1D-4D80-9182-19667418F9C2", out: true},
		{in: "5F0D34B3-24A5-4F76-BBC5-BF29E97C15AD", out: true},
		{in: "4C091301-FA37-4317-94D9-434C9D675B91", out: false},
		{in: "EF2BBB14-719A-4CFF-B1A0-62C79750FCD2", out: false},
		{in: "23F305C3-9D2C-49E0-A108-A8B12FC6534E", out: false},
		{in: "4FD676C8-2FC4-4BC9-A2B7-C78BF2CF21DC", out: false},
		{in: "7249B312-1CD5-4142-81EB-DB438AD45C54", out: false},
		{in: testDevice(0x1), out: true},
		{in: testDevice(0x2), out: true},
		{in: testDevice(0x3), out: true},
		{in: testDevice(0x4), out: true},
		{in: testDevice(0x5), out: true},
	}
	for _, tc := range tests {
		out := d.Del(tc.in)
		require.Equal(t, tc.out, out, "FAIL: %v", tc.in)
	}
}

func testIsEmpty(t *testing.T, d *Dict) {
	tests := []struct {
		in  *Dict
		out bool
	}{
		{in: &Dict{}, out: true},
		{in: New(), out: true},
		{in: New([]int{1, 2, 3}), out: false},
		{in: d, out: false},
	}
	for _, tc := range tests {
		out := tc.in.IsEmpty()
		require.Equal(t, tc.out, out)
	}
}

func testClear(t *testing.T, d *Dict) {
	tests := []struct {
		in  *Dict
		out bool
	}{
		{in: &Dict{}, out: false},
		{in: New(), out: false},
		{in: New(1, 2, 3), out: true},
		{in: New(nil), out: false},
		{in: New([]int{1, 2, 3}), out: true},
		{in: d, out: true},
	}
	for _, tc := range tests {
		out := tc.in.Clear()
		require.Equal(t, tc.out, out)
	}
}

func testKeys(t *testing.T, d *Dict) {
	tests := []struct {
		in  *Dict
		out []string
	}{
		{in: New(), out: nil},
		{in: New(nil), out: nil},
		{in: New([]int{}), out: nil},
		{in: New([]int{1, 2, 3}), out: []string{"0", "1", "2"}},
		{in: New(1, 2, 3), out: []string{"0", "1", "2"}},
		{in: New(map[int]struct{}{}), out: nil},
		{in: New(map[string]int{"one": 1, "two": 2, "three": 3}),
			out: []string{"one", "three", "two"}},
	}
	for _, tc := range tests {
		out := tc.in.Keys()
		if out != nil {
			sort.Strings(out)
		}
		require.EqualValues(t, tc.out, out)
	}
}

func testValues(t *testing.T, d *Dict) {
	tests := []struct {
		in  *Dict
		out []interface{}
	}{
		{in: New(), out: nil},
		{in: New(nil), out: nil},
		{in: New([]int{1, 2, 3}), out: []interface{}{int(1), int(2), int(3)}},
		{in: New(1, 2, 3), out: []interface{}{int(1), int(2), int(3)}},
		{in: New(1.1, 2.2, 3.3), out: []interface{}{float64(1.1), float64(2.2), float64(3.3)}},
		{in: New(map[int]string{1: "one", 2: "two", 3: "three"}),
			out: []interface{}{"one", "two", "three"}},
	}
	for _, tc := range tests {
		out := tc.in.Values()
		require.EqualValues(t, tc.out, out)
	}
}

func testPrint(t *testing.T, d *Dict) {
	tests := []struct {
		in  *Dict
		out string
	}{
		{in: New(), out: "{}"},
		{in: New(nil), out: "{}"},
		{in: New([]int{1, 2, 3}), out: "{0: 1, 1: 2, 2: 3}"},
		{in: New(1, 2, 3), out: "{0: 1, 1: 2, 2: 3}"},
		{in: New(1.1, 2.2, 3.3), out: "{0: 1.1, 1: 2.2, 2: 3.3}"},
	}
	for _, tc := range tests {
		out := fmt.Sprintf("%v", tc.in)
		require.Equal(t, tc.out, out)
	}
}

func testUpdate(t *testing.T, d *Dict) {
	tests := []struct {
		name string
		fn   func(*testing.T)
	}{
		{"Value slice", func(t *testing.T) {
			values := []struct {
				key, value interface{}
			}{
				{key: "testDevice11", value: "device 11"},
				{key: "testDevice22", value: "device 22"},
				{key: "testDevice33", value: "device 33"},
				{key: "testDevice44", value: "device 44"},
				{key: "testDevice55", value: "device 55"},
			}
			size, version := d.Len(), d.Version()
			ok := d.Update(values)
			require.True(t, ok)
			require.True(t, size != d.Len(), "expected size == %d but got %d", size, d.Len())
			require.True(t, version != d.Version(), "expected version == %d but got %d", version, d.Version())
		}},
		{"Item slice", func(t *testing.T) {
			values := []Item{
				{Key: testDevice(0x1), Value: "device 1"},
				{Key: testDevice(0x2), Value: "device 2"},
				{Key: testDevice(0x3), Value: "device 3"},
				{Key: testDevice(0x4), Value: "device 4"},
				{Key: testDevice(0x5), Value: "device 5"},
			}
			size, version := d.Len(), d.Version()
			ok := d.Update(values)
			require.True(t, ok)
			require.True(t, size != d.Len(), "expected size == %d but got %d", size, d.Len())
			require.True(t, version != d.Version(), "expected version == %d but got %d", version, d.Version())
		}},
		{"Value map", func(t *testing.T) {
			values := map[testDevice]interface{}{
				testDevice(0x1): "update device 1",
				testDevice(0x2): "update device 2",
				testDevice(0x3): "update device 3",
				testDevice(0x4): "update device 4",
				testDevice(0x5): "update device 5",
			}
			size, version := d.Len(), d.Version()
			ok := d.Update(values)
			require.True(t, ok)
			require.True(t, size == d.Len(), "expected size == %d but got %d", size, d.Len())
			require.True(t, version != d.Version(), "expected version == %d but got %d", version, d.Version())
		}},
		{"Value channel", func(t *testing.T) {
			values := []Item{
				{Key: testDevice(0x1), Value: "device 1"},
				{Key: testDevice(0x2), Value: "device 2"},
				{Key: testDevice(0x3), Value: "device 3"},
				{Key: testDevice(0x4), Value: "device 4"},
				{Key: testDevice(0x5), Value: "device 5"},
			}

			ci := make(chan Item)
			go func() {
				defer close(ci)
				for i := range values {
					ci <- values[i]
				}
			}()

			size, version := d.Len(), d.Version()
			ok := d.Update(ci)

			require.True(t, ok)
			require.True(t, size == d.Len(), "expected size == %d but got %d", size, d.Len())
			require.True(t, version != d.Version(), "expected version == %d but got %d", version, d.Version())
		}},
	}
	for _, tc := range tests {
		if !t.Run(tc.name, tc.fn) {
			break
		}
	}
}

func testWorkers(t *testing.T, d *Dict) {
	var wg sync.WaitGroup

	d.Clear()

	numWorkers := 7
	numItems := 100

	ch := make(chan Item)

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for item := range ch {
				d.Set(item.Key, item.Value)
			}
		}()
	}

	for i := 0; i < numItems; i++ {
		key := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		ch <- Item{key, value}
	}
	close(ch)
	wg.Wait()

	require.True(t, d.Len() == numItems, "expected size = %d, but got %d instead", numItems, d.Len())
	require.True(t, d.Version() > 0, "expected non-zero version, got %d", d.Version())
	v := d.Get("key99")
	require.NotNil(t, v)
	require.Equal(t, v, "value99")
}

func TestPopDoesNotDeadlock(t *testing.T) {
	d := New()
	d.Set("deadlock", "sentinel")

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = d.Pop("deadlock")
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Pop should complete but appears to be blocked")
	}
}

func TestPopItemShrinksLength(t *testing.T) {
	d := New(1, 2)
	require.Equal(t, 2, d.Len())

	item := d.PopItem()
	require.NotNil(t, item)
	require.Equal(t, 1, d.Len(), "PopItem should reduce the reported length")
}

func TestPopItemAfterLastElement(t *testing.T) {
	d := New()
	d.Set("only", "value")

	require.NotPanics(t, func() {
		out := d.PopItem()
		require.NotNil(t, out)
	})

	require.Zero(t, d.Len(), "Len should be zero after removing the last element")

	require.NotPanics(t, func() {
		require.Nil(t, d.PopItem(), "PopItem on an empty dict should return nil")
	})
}

func TestDelNoPanicWhenKeysSliceShrinks(t *testing.T) {
	d := New()
	d.Set("key", "value")

	d.mu.Lock()
	id := d.keys[0]
	d.keys = d.keys[:0]
	d.mu.Unlock()

	require.NotNil(t, id)

	require.NotPanics(t, func() {
		require.False(t, d.Del("key"))
	})
}

func TestItemsSnapshotSurvivesClear(t *testing.T) {
	d := New(1, 2, 3)
	ch := d.Items()
	d.Clear()

	var got []Item
	for item := range ch {
		got = append(got, item)
	}

	require.Len(t, got, 3, "Items should emit the snapshot taken before Clear")
}

func TestNilDictSafety(t *testing.T) {
	require.NotPanics(t, func() {
		var d *Dict
		require.True(t, d.IsEmpty())       // still treated as empty
		require.Nil(t, d.Get("missing"))   // read operations stay harmless
		require.False(t, d.Del("missing")) // write-style helpers no-op
	})
}

func TestPopItemConcurrentDoesNotPanic(t *testing.T) {
	d := New()
	d.Set("entry", 1)

	var wg sync.WaitGroup
	panicCh := make(chan interface{}, 2)
	start := make(chan struct{})

	worker := func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				panicCh <- r
			}
		}()
		<-start
		_ = d.PopItem()
	}

	wg.Add(2)
	go worker()
	go worker()

	close(start)
	wg.Wait()
	close(panicCh)

	for p := range panicCh {
		require.Failf(t, "unexpected panic", "%v", p)
	}
}

func TestPopItemDoesNotPanicWithStaleSize(t *testing.T) {
	d := New()
	atomic.StoreInt64(&d.size, 1)

	require.NotPanics(t, func() {
		require.Nil(t, d.PopItem())
	})
}
