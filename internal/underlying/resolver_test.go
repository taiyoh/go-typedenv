package underlying

import (
	"errors"
	"testing"
	"time"
)

func TestResolve(t *testing.T) {
	t.Run("uint", func(t *testing.T) {
		type foo8 uint8
		var f foo8
		if err := Resolve("16", &f); err != nil {
			t.Error(err)
		}
		if f != foo8(16) {
			t.Errorf("want: 16 but got: %d", f)
		}
	})
	t.Run("int", func(t *testing.T) {
		type bar16 int16
		var f bar16
		if err := Resolve("1666", &f); err != nil {
			t.Error(err)
		}
		if f != bar16(1666) {
			t.Errorf("want: 1666 but got: %d", f)
		}
	})
	t.Run("*int32", func(t *testing.T) {
		type baz32 int32
		var f *baz32
		if err := Resolve("1666", &f); err != nil {
			t.Error(err)
		}
		if *f != baz32(1666) {
			t.Errorf("want: 1666 but got: %d", *f)
		}
	})
	t.Run("*int32 part2", func(t *testing.T) {
		type baz32 *int32
		var f baz32
		if err := Resolve("1666", &f); err != nil {
			t.Error(err)
		}
		v := int32(1666)
		if *f != *baz32(&v) {
			t.Errorf("want: 1666 but got: %d", *f)
		}
	})

	t.Run("time", func(t *testing.T) {
		var ts time.Time
		if err := Resolve("2011-03-11T14:45:00+09:00", &ts); err != nil {
			t.Error(err)
		}
		jst, _ := time.LoadLocation("Asia/Tokyo")
		if !ts.Equal(time.Date(2011, 3, 11, 14, 45, 0, 0, jst)) {
			t.Errorf("wrong date: %s", ts)
		}
	})

	t.Run("duration", func(t *testing.T) {
		var d time.Duration
		if err := Resolve("30s", &d); err != nil {
			t.Error(err)
		}
		if d != 30*time.Second {
			t.Errorf("wrong duration: %v", d)
		}
	})

	t.Run("unknown", func(t *testing.T) {
		type hoge struct{}
		var h hoge
		if err := Resolve("hoge", &h); !errors.Is(err, ErrNotFound) {
			t.Error(err)
		}
	})
}
