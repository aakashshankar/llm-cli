package mocks

import "github.com/stretchr/testify/mock"

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Prompt(prompt string, stream bool, tokens int, model string, system string, clear bool) (string, error) {
	args := m.Called(prompt, stream, tokens, model, system, clear)
	return args.String(0), args.Error(1)
}
