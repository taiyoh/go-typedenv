package typedenv

import (
	"reflect"
	"testing"
)

func TestScan(t *testing.T) {
	t.Setenv("FOO", "123")
	t.Setenv("BAR", "true")
	t.Setenv("BAZ", "abc,def")
	t.Setenv("BAZ2", "ghi:jklm")

	type bazTyp []string
	type fooTyp map[string]struct{}
	type hogeInt int

	var (
		foo   string
		fooM  fooTyp
		bar   bool
		baz   bazTyp
		hoge  []hogeInt
		hoge2 []hogeInt
		hoge3 bazTyp
		fuga  map[hogeInt]struct{}
		fuga2 map[hogeInt]struct{}
		fuga3 fooTyp
		piyo  hogeInt
	)

	if err := Scan(
		RequiredDirect("FOO", &foo),
		LookupDirect("BAR", &bar),
		DefaultDirect("PIYO", &piyo, "33"),
		Required("BAZ", Slice(&baz)),
		Default("HOGE", Slice(&hoge), "123,456"),
		Lookup("HOGE2", Slice(&hoge2)),
		Lookup("BAZ2", Slice(&hoge3, ":")),
		Required("BAZ", Map(&fooM)),
		Default("HOGE", Map(&fuga), "123,456"),
		Lookup("HOGE2", Map(&fuga2)),
		Lookup("BAZ2", Map(&fuga3, ":")),
	); err != nil {
		t.Error(err)
	}
	if foo != "123" {
		t.Errorf("foo is captured: %s", foo)
	}
	if !bar {
		t.Errorf("bar is captured: %v", bar)
	}
	if piyo != 33 {
		t.Errorf("piyo is captured: %v", piyo)
	}
	if !reflect.DeepEqual(baz, bazTyp{"abc", "def"}) {
		t.Errorf("baz is captured: %+v", baz)
	}
	if !reflect.DeepEqual(hoge, []hogeInt{123, 456}) {
		t.Errorf("hoge is captured: %+v", hoge)
	}
	if len(hoge2) != 0 {
		t.Errorf("found hoge2: %#v", hoge2)
	}
	if !reflect.DeepEqual(hoge3, bazTyp{"ghi", "jklm"}) {
		t.Errorf("hoge3 is captured: %+v", hoge3)
	}
	if !reflect.DeepEqual(fooM, fooTyp{"abc": {}, "def": {}}) {
		t.Errorf("fooM is captured: %+v", fooM)
	}
	if !reflect.DeepEqual(fuga, map[hogeInt]struct{}{123: {}, 456: {}}) {
		t.Errorf("fuga is captured: %+v", fuga)
	}
	if len(fuga2) != 0 {
		t.Errorf("found fuga2: %#v", fuga2)
	}
	if !reflect.DeepEqual(fuga3, fooTyp{"ghi": {}, "jklm": {}}) {
		t.Errorf("fuga3 is captured: %+v", fuga3)
	}
}
