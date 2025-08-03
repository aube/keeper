package decrypt

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aube/keeper/internal/client/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) DecryptFile(inputName, outputPath, password string) error {
	args := m.Called(inputName, outputPath, password)
	return args.Error(0)
}

func (m *MockFileRepository) Save(ctx context.Context, filename string, data io.Reader) error {
	return nil
}

func (m *MockFileRepository) FindAll(ctx context.Context) (*entities.Files, error) {
	return nil, nil
}

func (m *MockFileRepository) Delete(ctx context.Context, uuid string) error {
	return nil
}

func (m *MockFileRepository) GetFileContent(ctx context.Context, uuid string) (string, error) {
	return "", nil
}

func TestRun(t *testing.T) {
	mockRepo := new(MockFileRepository)

	tests := []struct {
		name      string
		password  string
		inputName string
		output    string
		mockErr   error
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "success",
			password:  "pass",
			inputName: "input",
			output:    "output",
			mockErr:   nil,
			wantErr:   false,
		},
		{
			name:      "empty password",
			password:  "",
			inputName: "input",
			output:    "output",
			wantErr:   true,
			errMsg:    "empty password",
		},
		{
			name:      "decryption error",
			password:  "pass",
			inputName: "input",
			output:    "output",
			mockErr:   errors.New("decryption failed"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.password != "" {
				mockRepo.On("DecryptFile", tt.inputName, tt.output, tt.password).Return(tt.mockErr).Once()
			}

			err := Run(tt.password, tt.inputName, tt.output, mockRepo)
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
