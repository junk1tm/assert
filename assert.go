// Package assert provides common assertions to use with the standard [testing] package.
package assert

import (
	"errors"
	"reflect"
)

// TB is a tiny subset of [testing.TB] used by [assert].
type TB interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

// Param controls the behaviour of an assertion in case it fails.
// Either [E] or [F] must be specified as a type parameter when calling the assertion.
type Param interface {
	method(t TB) func(format string, args ...any)
}

// E is a [Param] that marks the test as having failed but continues its execution (similar to [testing.T.Errorf]).
type E struct{}

func (E) method(t TB) func(format string, args ...any) { return t.Errorf }

// F is a [Param] that marks the test as having failed and stops its execution (similar to [testing.T.Fatalf]).
type F struct{}

func (F) method(t TB) func(format string, args ...any) { return t.Fatalf }

// Equal asserts that got and want are equal.
func Equal[T Param, V any](t TB, got, want V, formatAndArgs ...any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		fail[T](t, formatAndArgs, "\ngot\t%v\nwant\t%v", got, want)
	}
}

// NoErr asserts that err is nil.
func NoErr[T Param](t TB, err error, formatAndArgs ...any) {
	t.Helper()
	if err != nil {
		fail[T](t, formatAndArgs, "\ngot\t%v\nwant\tno error", err)
	}
}

// IsErr asserts that [errors.Is](err, target) is true.
func IsErr[T Param](t TB, err, target error, formatAndArgs ...any) {
	t.Helper()
	if !errors.Is(err, target) {
		fail[T](t, formatAndArgs, "\ngot\t%v\nwant\t%v", err, target)
	}
}

// AsErr asserts that [errors.As](err, target) is true.
func AsErr[T Param](t TB, err error, target any, formatAndArgs ...any) {
	t.Helper()
	if !errors.As(err, target) {
		fail[T](t, formatAndArgs, "\ngot\t%T\nwant\t%T", err, target)
	}
}

// fail marks the test as having failed and continues/stops its execution based on T's type.
func fail[T Param](t TB, customFormatAndArgs []any, format string, args ...any) {
	t.Helper()
	if len(customFormatAndArgs) > 0 {
		format = customFormatAndArgs[0].(string)
		args = customFormatAndArgs[1:]
	}
	(*new(T)).method(t)(format, args...)
}
