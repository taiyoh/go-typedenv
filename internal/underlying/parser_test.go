package underlying

import (
	"reflect"
	"testing"
)

func TestCast(t *testing.T) {
	type foo1 int32
	type foo2 []foo1
	type foo3 []int32

	type bar1 uint16
	type bar2 map[bar1]struct{}
	type bar3 map[uint16]struct{}

	t.Run("CastSlice 1", func(t *testing.T) {
		f := foo2{1, 2, 3}
		ff := CastSlice(&f)
		if !reflect.DeepEqual(*ff, []foo1{1, 2, 3}) {
			t.Errorf("want: []foo1{1, 2, 3}, got: %+v", ff)
		}
	})

	t.Run("CastSlice 2", func(t *testing.T) {
		f := foo3{1, 2, 3}
		ff := CastSlice(&f)
		if !reflect.DeepEqual(*ff, []int32{1, 2, 3}) {
			t.Errorf("want: []int32{1, 2, 3}, got: %+v", ff)
		}
	})

	t.Run("CastMap 1", func(t *testing.T) {
		b := bar2{1: {}, 2: {}, 3: {}}
		bb := CastMap(&b)
		if !reflect.DeepEqual(*bb, map[bar1]struct{}{1: {}, 2: {}, 3: {}}) {
			t.Errorf("want: map[bar1]struct{}{1: {}, 2: {}, 3: {}}, got: %+v", bb)
		}
	})

	t.Run("CastMap 2", func(t *testing.T) {
		b := bar3{1: {}, 2: {}, 3: {}}
		bb := CastMap(&b)
		if !reflect.DeepEqual(*bb, map[uint16]struct{}{1: {}, 2: {}, 3: {}}) {
			t.Errorf("want: map[uint16]struct{}{1: {}, 2: {}, 3: {}}, got: %+v", bb)
		}
	})
}
