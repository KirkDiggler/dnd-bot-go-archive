package types

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type TimeClock interface {
	Now() time.Time
}

type Clock struct{}

func (c *Clock) Now() time.Time {
	return time.Now()
}

type MockClock struct {
	mock.Mock
}

func (m *MockClock) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
