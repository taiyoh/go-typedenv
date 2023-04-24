package underlying

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

func intConverter(str string) (interface{}, error) {
	return strconv.ParseInt(str, 10, 64)
}

func uintConverter(str string) (interface{}, error) {
	return strconv.ParseUint(str, 10, 64)
}

func nopConverter(str string) (interface{}, error) {
	return str, nil
}

func floatConverter(str string) (interface{}, error) {
	return strconv.ParseFloat(str, 64)
}

func boolConverter(str string) (interface{}, error) {
	return strconv.ParseBool(str)
}

func timeConverter(str string) (interface{}, error) {
	return time.Parse(time.RFC3339Nano, str)
}

func durationConverter(str string) (interface{}, error) {
	return time.ParseDuration(str)
}

var ptrTable = map[reflect.Kind]func(str string) (interface{}, error){
	reflect.Int:     intConverter,
	reflect.Int8:    intConverter,
	reflect.Int16:   intConverter,
	reflect.Int32:   intConverter,
	reflect.Int64:   intConverter,
	reflect.Uint:    uintConverter,
	reflect.Uint8:   uintConverter,
	reflect.Uint16:  uintConverter,
	reflect.Uint32:  uintConverter,
	reflect.Uint64:  uintConverter,
	reflect.String:  nopConverter,
	reflect.Float32: floatConverter,
	reflect.Float64: floatConverter,
	reflect.Bool:    boolConverter,
}

// Resolve resolves underlying type of value which it receives.
// supported: (via https://cs.opensource.google/go/x/exp/+/master:constraints/constraints.go)
// - int, int8, int16, int32, int64 (equivarent of constraints.Signed)
// - uint, uint8, uint16, uint32, uint64 (equivarent of constraints.Unsigned)
// - float32, float64 (equivarent of constraints.Float)
// - string
// - bool
// - time.Time
// - time.Duration
// if this function receives value except aboves, returns nil.
func Resolve(env string, val any) error {
	var fn func(string) (interface{}, error)
	switch val.(type) {
	case *time.Time:
		fn = timeConverter
	case *time.Duration:
		fn = durationConverter
	}
	rv := reflect.ValueOf(val)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		kind = rv.Elem().Kind()
		rv = rv.Elem()
		rv.Set(reflect.New(rv.Type()).Elem())
	}
	// if val is pointer-pointer type like **int,
	// dereference again
	if kind == reflect.Ptr {
		rv.Set(reflect.New(rv.Type().Elem()))
		kind = rv.Elem().Kind()
		rv = rv.Elem()
	}
	if fn == nil {
		var ok bool
		if fn, ok = ptrTable[kind]; !ok {
			return ErrNotFound
		}
	}
	vv, err := fn(env)
	if err != nil {
		return err
	}
	rv.Set(reflect.ValueOf(vv).Convert(rv.Type()))
	return nil
}

// ErrNotFound shows supplied value could not be found for supported underlying type.
var ErrNotFound = errors.New("not found")
