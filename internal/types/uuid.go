package types

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UUIDGenerator interface {
	New() string
}

type GoogleUUID struct{}

func (g *GoogleUUID) New() string {
	return uuid.New().String()
}

type MockUUID struct {
	mock.Mock
}

func (m *MockUUID) New() string {
	args := m.Called()
	return args.String(0)
}
