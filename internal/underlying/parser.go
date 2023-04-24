package underlying

import (
	"reflect"
	"time"
)

type numbers interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type numbersPtr interface {
	~*float32 | ~*float64 |
		~*int | ~*int8 | ~*int16 | ~*int32 | ~*int64 |
		~*uint | ~*uint8 | ~*uint16 | ~*uint32 | ~*uint64
}

// Basic sets underlying basic types and its pointer types.
type Basic interface {
	~string | ~*string | ~bool | ~*bool | numbers | numbersPtr | time.Time
}

// CastSlice converts `[]T` from `U` (defined type of slice).
// Because if `U` represents as defined type, convert process will be so complicated.
func CastSlice[T Basic, U ~[]T](val *U) (s *[]T) {
	rv := reflect.ValueOf(val)
	return rv.Convert(reflect.TypeOf(s)).Interface().(*[]T)
}

// CastMap converts `map[T]struct{}` from `U` (defined type of map).
// Because if `U` represents as defined type, convert process will be so complicated.
func CastMap[T Basic, U ~map[T]struct{}](val *U) (m *map[T]struct{}) {
	rv := reflect.ValueOf(val)
	return rv.Convert(reflect.TypeOf(m)).Interface().(*map[T]struct{})
}
