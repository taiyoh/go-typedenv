package typedenv

import (
	"encoding"
	"fmt"
	"os"
	"strings"

	"github.com/taiyoh/go-typedenv/internal/underlying"
)

// Scanner sets some rules to lookup and convert env var.
type Scanner struct {
	keyName     string
	required    bool
	defaultVal  *string
	unmarshaler encoding.TextUnmarshaler
}

// Scan scans environment variables and converts to supplied types.
func Scan(scanners ...Scanner) error {
	for _, scanner := range scanners {
		val, ok := os.LookupEnv(scanner.keyName)
		if !ok {
			if scanner.required {
				return fmt.Errorf("%s is required", scanner.keyName)
			}
			if scanner.defaultVal != nil {
				val = *scanner.defaultVal
			}
		}

		// skip if value is not set
		if val == "" {
			continue
		}

		if err := scanner.unmarshaler.UnmarshalText([]byte(val)); err != nil {
			return err
		}
	}
	return nil
}

// Required returns Scanner with required flag.
// This function requires Converter.
func Required(key string, u encoding.TextUnmarshaler) Scanner {
	return Scanner{
		keyName:     key,
		required:    true,
		unmarshaler: u,
	}
}

// Default returns Scanner with default value.
// This function requires Converter.
func Default(key string, u encoding.TextUnmarshaler, defaultVal string) Scanner {
	return Scanner{
		keyName:     key,
		defaultVal:  &defaultVal,
		unmarshaler: u,
	}
}

// Lookup returns Scanner without required and default.
// This function requires Converter.
func Lookup(key string, u encoding.TextUnmarshaler) Scanner {
	return Scanner{
		keyName:     key,
		unmarshaler: u,
	}
}

type mapUnmarshaler[T underlying.Basic] struct {
	target    *map[T]struct{}
	separator string
}

var _ encoding.TextUnmarshaler = (*mapUnmarshaler[int])(nil)

func (c *mapUnmarshaler[T]) UnmarshalText(source []byte) error {
	sources := strings.Split(string(source), c.separator)
	targets := make([]T, len(sources))
	for i, s := range sources {
		if err := underlying.Resolve(s, &targets[i]); err != nil {
			return err
		}
	}
	*c.target = make(map[T]struct{}, len(targets))
	for _, tgt := range targets {
		(*c.target)[tgt] = struct{}{}
	}
	return nil
}

// Map returns encoding.TextUnmarshaler implementation for map.
func Map[T underlying.Basic, U ~map[T]struct{}](val *U, seps ...string) encoding.TextUnmarshaler {
	sep := "," // default separator
	if len(seps) > 0 {
		sep = seps[0]
	}
	return &mapUnmarshaler[T]{target: underlying.CastMap(val), separator: sep}
}

type sliceUnmarshaler[T underlying.Basic] struct {
	target    *[]T
	separator string
}

var _ encoding.TextUnmarshaler = (*sliceUnmarshaler[int])(nil)

func (c *sliceUnmarshaler[T]) UnmarshalText(source []byte) error {
	sources := strings.Split(string(source), c.separator)
	targets := make([]T, len(sources))
	for i, s := range sources {
		if err := underlying.Resolve(s, &targets[i]); err != nil {
			return err
		}
	}
	*c.target = targets
	return nil
}

// Slice returns encoding.TextUnmarshaler implementation for slice.
func Slice[T underlying.Basic, U ~[]T](val *U, seps ...string) encoding.TextUnmarshaler {
	sep := "," // default separator
	if len(seps) > 0 {
		sep = seps[0]
	}
	return &sliceUnmarshaler[T]{target: underlying.CastSlice(val), separator: sep}
}

type directUnmarshaler[T underlying.Basic] struct {
	target *T
}

var _ encoding.TextUnmarshaler = (*directUnmarshaler[int])(nil)

func (c *directUnmarshaler[T]) UnmarshalText(source []byte) error {
	var val T
	if err := underlying.Resolve(string(source), &val); err != nil {
		return err
	}
	*c.target = val
	return nil
}

// Direct returns encoding.TextUnmarshaler implementation for basic types.
func Direct[T underlying.Basic](val *T) encoding.TextUnmarshaler {
	return &directUnmarshaler[T]{target: val}
}

// RequiredDirect returns Scanner with required flag.
func RequiredDirect[T underlying.Basic](key string, val *T) Scanner {
	return Required(key, Direct(val))
}

// DefaultDirect returns Scanner with default value.
func DefaultDirect[T underlying.Basic](key string, val *T, defaultVal string) Scanner {
	return Default(key, Direct(val), defaultVal)
}

// LookupDirect returns Scanner without required and default.
func LookupDirect[T underlying.Basic](key string, val *T) Scanner {
	return Lookup(key, Direct(val))
}
