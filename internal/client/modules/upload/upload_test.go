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

func (m *MockFileRepository) GetPath(filename string) string {
	args := m.Called(filename)
	return args.String(0)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) UploadFile(ctx context.Context, endpoint string, filePath string, formFields map[string]string) ([]byte, error) {
	args := m.Called(ctx, endpoint, filePath, formFields)
	return args.Get(0).([]byte), args.Error(1)
}

func TestRun(t *testing.T) {
	mockRepo := new(MockFileRepository)
	mockHTTP := new(MockHTTPClient)
	ctx := context.Background()

	tests := []struct {
		name      string
		filename  string
		category  string
		path      string
		uploadErr error
		wantErr   bool
	}{
		{
			name:     "success",
			filename: "test.txt",
			category: "docs",
			path:     "/path/to/test.txt",
			wantErr:  false,
		},
		{
			name:      "upload error",
			filename:  "test.txt",
			category:  "docs",
			path:      "/path/to/test.txt",
			uploadErr: assert.AnError,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedFields := map[string]string{
				"description": "ololo alala",
				"category":    tt.category,
			}

			mockRepo.On("GetPath", tt.filename).Return(tt.path).Once()
			mockHTTP.On("UploadFile", ctx, "/upload", tt.path, expectedFields).Return([]byte(`{"uuid":"123"}`), tt.uploadErr).Once()

			err := Run(mockRepo, tt.filename, tt.category, mockHTTP)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockHTTP.AssertExpectations(t)
		})
	}
}
