package download

import (
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

func (m *MockHTTPClient) DownloadFile(fileURL, outputPath string) error {
	args := m.Called(fileURL, outputPath)
	return args.Error(0)
}

func TestRun(t *testing.T) {
	mockRepo := new(MockFileRepository)
	mockHTTP := new(MockHTTPClient)

	tests := []struct {
		name    string
		input   string
		path    string
		httpErr error
		wantErr bool
	}{
		{
			name:    "success",
			input:   "test.txt",
			path:    "/path/to/test.txt",
			httpErr: nil,
			wantErr: false,
		},
		{
			name:    "download error",
			input:   "test.txt",
			path:    "/path/to/test.txt",
			httpErr: assert.AnError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetPath", tt.input).Return(tt.path).Once()
			mockHTTP.On("DownloadFile", "/file?name="+tt.input, tt.path).Return(tt.httpErr).Once()

			err := Run(tt.input, mockRepo, mockHTTP)
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
