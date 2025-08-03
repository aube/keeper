package register

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Post(endpoint string, body any) ([]byte, error) {
	args := m.Called(endpoint, body)
	return args.Get(0).([]byte), args.Error(1)
}

func TestRun(t *testing.T) {
	mockHTTP := new(MockHTTPClient)

	tests := []struct {
		name     string
		username string
		password string
		email    string
		postErr  error
		wantErr  bool
	}{
		{
			name:     "success",
			username: "user",
			password: "pass",
			email:    "user@example.com",
			postErr:  nil,
			wantErr:  false,
		},
		{
			name:     "post error",
			username: "user",
			password: "pass",
			email:    "user@example.com",
			postErr:  errors.New("post failed"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedBody := map[string]interface{}{
				"username": tt.username,
				"password": tt.password,
				"email":    tt.email,
			}

			mockHTTP.On("Post", "/register", expectedBody).Return([]byte{}, tt.postErr).Once()

			err := Run(tt.username, tt.password, tt.email, mockHTTP)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockHTTP.AssertExpectations(t)
		})
	}
}
