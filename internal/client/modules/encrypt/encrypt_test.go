package encrypt

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) EncryptFile(inputPath, outputName, password string) error {
	args := m.Called(inputPath, outputName, password)
	return args.Error(0)
}

func TestRun(t *testing.T) {
	mockRepo := new(MockFileRepository)

	tests := []struct {
		name      string
		password  string
		inputPath string
		output    string
		mockErr   error
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "success",
			password:  "pass",
			inputPath: "input.txt",
			output:    "output.enc",
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "empty password",
			password:  "",
			inputPath: "input.txt",
			output:    "output.enc",
			wantErr:   true,
			errMsg:    "empty password",
		},
		{
			name:      "encryption error",
			password:  "pass",
			inputPath: "input.txt",
			output:    "output.enc",
			mockErr:   errors.New("encryption failed"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.password != "" {
				mockRepo.On("EncryptFile", tt.inputPath, tt.output, tt.password).Return(tt.mockErr).Once()
			}

			err := Run(tt.password, tt.inputPath, tt.output, mockRepo)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
