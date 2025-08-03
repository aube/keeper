package card

import (
	"os"
	"path/filepath"
	"strings"
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
		name     string
		number   string
		date     string
		cvv      string
		password string
		mockErr  error
		wantErr  bool
	}{
		{
			name:     "success",
			number:   "1234 5678 9012 3456",
			date:     "12/25",
			cvv:      "123",
			password: "secure",
			mockErr:  nil,
			wantErr:  false,
		},
		{
			name:     "encryption error",
			number:   "1234 5678 9012 3456",
			date:     "12/25",
			cvv:      "123",
			password: "secure",
			mockErr:  assert.AnError,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("EncryptFile", mock.Anything, mock.Anything, tt.password).Return(tt.mockErr).Once()

			_, err := Run(tt.number, tt.date, tt.cvv, tt.password, mockRepo, nil)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Clean up temp file
				filename := "card_" + strings.ReplaceAll(tt.number, " ", "") + ".json"
				os.Remove(filepath.Join(os.TempDir(), filename))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExtractUUID(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:    "valid response",
			input:   []byte(`{"uuid": "12345"}`),
			want:    "12345",
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   []byte(`invalid`),
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractUUID(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
