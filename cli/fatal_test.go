package cli

import (
	"bytes"
	stdlog "log"
	"runtime/debug"
	"sync"
	"testing"
)

var (
	onceFatalMock sync.Once
	dieMock       *fatalMock
)

// fatalMock tracks calls to Die (aka log.Fatal) from tests run in parallel.
type fatalMock struct {
	mx            sync.Mutex
	once          sync.Once
	fatalCalledBy map[string]bool
	reset         func()
	fatals        map[string]func(string, ...interface{})
}

func newFatalMock() *fatalMock {
	return &fatalMock{
		fatalCalledBy: make(map[string]bool),
		fatals:        make(map[string]func(string, ...interface{})),
		reset: func() {
			SetDie(stdlog.Fatalf)
		},
	}
}

// Reset returns the original defaut Fatal behavior for this package.
func (f *fatalMock) Reset() {
	f.reset()
}

// Called indicates whether Fatal was called by this test.
func (f *fatalMock) Called(t testing.TB) bool {
	f.mx.Lock()
	defer f.mx.Unlock()

	return f.fatalCalledBy[t.Name()]
}

// Register a "log.Fatal" mock for the current test, in the current go routine.
func (f *fatalMock) Register(t testing.TB) {
	f.mx.Lock()
	defer f.mx.Unlock()

	localFatal := func(msg string, args ...interface{}) {
		// a Fatal method that is private to this goroutine
		f.fatalCalledBy[t.Name()] = true

		t.Logf(msg, args...)
	}

	current := whoami()
	f.fatals[current] = localFatal

	f.once.Do(func() {
		SetDie(func(msg string, args ...interface{}) {
			me := whoami()
			f.mx.Lock()

			for who, fn := range f.fatals {
				if who != me {
					continue
				}

				fn(msg, args...)
			}

			f.mx.Unlock()

		})
	})
}

// whoami returns the current goroutine
func whoami() string {
	return string(bytes.Fields(debug.Stack())[1])
}
