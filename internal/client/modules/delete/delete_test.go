package upload

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Delete(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

func TestRun(t *testing.T) {
	mockRepo := new(MockFileRepository)
	mockHTTP := new(MockHTTPClient)
	ctx := context.Background()

	tests := []struct {
		name     string
		filename string
		repoErr  error
		httpErr  error
		wantErr  bool
	}{
		{
			name:     "success",
			filename: "test.txt",
			repoErr:  nil,
			httpErr:  nil,
			wantErr:  false,
		},
		{
			name:     "http error",
			filename: "test.txt",
			repoErr:  nil,
			httpErr:  assert.AnError,
			wantErr:  true,
		},
		{
			name:     "repo error",
			filename: "test.txt",
			repoErr:  assert.AnError,
			httpErr:  nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTP.On("Delete", "/delete?name"+tt.filename).Return(tt.httpErr).Once()
			if tt.httpErr == nil {
				mockRepo.On("Delete", ctx, tt.filename).Return(tt.repoErr).Once()
			}

			err := Run(mockRepo, tt.filename, mockHTTP)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockHTTP.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
