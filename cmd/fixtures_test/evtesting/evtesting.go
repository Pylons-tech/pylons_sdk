package evtesting

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// T is a modified testing.T
type T struct {
	origin    *testing.T
	useLogPkg bool
}

var listeners = make(map[string]func())

// NewT is function returns modified T from original testing.T
func NewT(origin *testing.T) T {
	newT := T{
		origin:    origin,
		useLogPkg: false,
	}
	if origin == nil {
		orgT := testing.T{}
		newT.origin = &orgT
		newT.useLogPkg = true
	}
	return newT
}

// Fatal is a modified Fatal
func (t *T) Fatal(args ...interface{}) {
	t.DispatchEvent("FAIL")
	if t.useLogPkg {
		log.Fatal(args...)
	} else {
		t.origin.Fatal(args...)
	}
}

// Fatalf is a modified Fatalf
func (t *T) Fatalf(format string, args ...interface{}) {
	t.DispatchEvent("FAIL")
	if t.useLogPkg {
		log.Fatalf(format, args...)
	} else {
		t.origin.Fatalf(format, args...)
	}
}

// MustTrue validate if value is true
func (t *T) MustTrue(value bool) {
	if !value {
		t.DispatchEvent("FAIL")
	}
	if t.useLogPkg {
		if !value {
			log.Fatal("MustTrue validation failed")
		}
	} else {
		require.True(t.origin, value)
	}
}

// MustNil validate if value is nil
func (t *T) MustNil(err error) {
	if err != nil {
		t.Log("comparing \"", err, "\" to nil")
	}
	t.MustTrue(err == nil)
}

// Parallel is modified Parallel
func (t *T) Parallel() {
	t.origin.Parallel()
}

// Log is modified Log
func (t *T) Log(args ...interface{}) {
	if t.useLogPkg {
		log.Println(args...)
	} else {
		t.origin.Log(args...)
	}
}

// Run is modified Run
func (t *T) Run(name string, f func(t *T)) bool {
	return t.origin.Run(name, func(t *testing.T) {
		newT := T{
			origin: t,
		}
		f(&newT)
	})
}

// DispatchEvent process events that are related to the event e.g. failure in one test case make others to fail without continuing
func (t *T) DispatchEvent(event string) {
	if listener, ok := listeners[event]; ok {
		listener()
	}
}
