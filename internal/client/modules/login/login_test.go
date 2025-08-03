package login

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Save(ctx context.Context, filename string, data io.Reader) error {
	args := m.Called(ctx, filename, data)
	return args.Error(0)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Post(endpoint string, body any) ([]byte, error) {
	args := m.Called(endpoint, body)
	return args.Get(0).([]byte), args.Error(1)
}

func TestRun(t *testing.T) {
	mockRepo := new(MockFileRepository)
	mockHTTP := new(MockHTTPClient)
	ctx := context.Background()

	tests := []struct {
		name     string
		username string
		password string
		token    string
		postErr  error
		saveErr  error
		wantErr  bool
	}{
		{
			name:     "success",
			username: "user",
			password: "pass",
			token:    "token123",
			postErr:  nil,
			saveErr:  nil,
			wantErr:  false,
		},
		{
			name:     "post error",
			username: "user",
			password: "pass",
			postErr:  assert.AnError,
			wantErr:  true,
		},
		{
			name:     "save error",
			username: "user",
			password: "pass",
			token:    "token123",
			postErr:  nil,
			saveErr:  assert.AnError,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.token != "" {
				mockHTTP.On("Post", "/login", mock.Anything).Return([]byte(`{"token":"`+tt.token+`"}`), tt.postErr).Once()
			} else {
				mockHTTP.On("Post", "/login", mock.Anything).Return([]byte{}, tt.postErr).Once()
			}

			if tt.postErr == nil && tt.token != "" {
				mockRepo.On("Save", ctx, tt.username, mock.Anything).Return(tt.saveErr).Once()
			}

			err := Run(tt.username, tt.password, mockRepo, mockHTTP)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockHTTP.AssertExpectations(t)
			if tt.postErr == nil && tt.token != "" {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
