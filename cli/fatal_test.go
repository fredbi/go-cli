package cli

import (
	"sync"
	"testing"
)

type fatalMock struct {
	mx       sync.Mutex
	calledBy []testing.TB
}

func newFatalMock() *fatalMock {
	return &fatalMock{}
}

func (m *fatalMock) Fatalf(format string, args ...any) {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.calledBy = append(m.calledBy, nil) // TODO: captures caller go routine
}
