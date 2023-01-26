package choice

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Get(ctx context.Context, input GetInput) (*GetOutput, error) {
	ret := m.Called(ctx, input)

	if ret.Error(1) != nil {
		return nil, ret.Error(1)
	}

	return ret.Get(0).(*GetOutput), nil
}

func (m *Mock) Put(ctx context.Context, input *PutInput) error {
	ret := m.Called(ctx, input)

	return ret.Error(0)
}
