package sync

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

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

func (m *MockFileRepository) GetPath(filename string) string {
	args := m.Called(filename)
	return args.String(0)
}

func (m *MockFileRepository) Exists(filename string) bool {
	args := m.Called(filename)
	return args.Bool(0)
}

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) Save(ctx context.Context, filename string, data io.Reader) error {
	args := m.Called(ctx, filename, data)
	return args.Error(0)
}

func (m *MockTokenRepository) GetFileContent(ctx context.Context, uuid string) (string, error) {
	args := m.Called(ctx, uuid)
	return args.String(0), args.Error(1)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(endpoint string, queryParams map[string]string) ([]byte, error) {
	args := m.Called(endpoint, queryParams)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockHTTPClient) DownloadFile(fileURL, outputPath string) error {
	args := m.Called(fileURL, outputPath)
	return args.Error(0)
}

func TestRun(t *testing.T) {
	now := time.Now()
	lastSyncTime := now.AddDate(-1, 0, 0)
	username := "testuser"

	tests := []struct {
		name          string
		prevSyncTime  string
		getSyncErr    error
		deletedFiles  []Row
		newFiles      []Row
		deletedExists []bool
		newExists     []bool
		deleteErr     error
		downloadErr   error
		getDeletedErr error
		getNewErr     error
		saveSyncErr   error
		wantErr       bool
	}{
		{
			name:          "success with previous sync",
			prevSyncTime:  lastSyncTime.Format(time.RFC3339),
			deletedFiles:  []Row{{Name: "file1.txt"}, {Name: "file2.txt"}},
			newFiles:      []Row{{Name: "file3.txt"}, {Name: "file4.txt"}},
			deletedExists: []bool{true, false},
			newExists:     []bool{false, true},
			wantErr:       false,
		},
		{
			name:          "get deleted files error",
			prevSyncTime:  lastSyncTime.Format(time.RFC3339),
			getDeletedErr: errors.New("server error"),
			wantErr:       true,
		},
		{
			name:          "delete file error",
			prevSyncTime:  lastSyncTime.Format(time.RFC3339),
			deletedFiles:  []Row{{Name: "file1.txt"}},
			deletedExists: []bool{true},
			deleteErr:     errors.New("delete failed"),
			wantErr:       true,
		},
		{
			name:          "get new files error",
			prevSyncTime:  lastSyncTime.Format(time.RFC3339),
			deletedFiles:  []Row{{Name: "file1.txt"}},
			deletedExists: []bool{false},
			getNewErr:     errors.New("server error"),
			wantErr:       true,
		},
		{
			name:         "download file error",
			prevSyncTime: lastSyncTime.Format(time.RFC3339),
			newFiles:     []Row{{Name: "file1.txt"}},
			newExists:    []bool{false},
			downloadErr:  errors.New("download failed"),
			wantErr:      true,
		},
		{
			name:         "save sync time error",
			prevSyncTime: lastSyncTime.Format(time.RFC3339),
			deletedFiles: []Row{},
			newFiles:     []Row{},
			saveSyncErr:  errors.New("save failed"),
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockFileRepo := new(MockFileRepository)
			mockTokenRepo := new(MockTokenRepository)
			mockHTTP := new(MockHTTPClient)

			// Setup expectations for getting previous sync time
			mockTokenRepo.On("GetFileContent", ctx, username).Return(tt.prevSyncTime, tt.getSyncErr).Once()

			// Calculate expected uploaded_at parameter
			uploadedAt := lastSyncTime.Format("2006-01-02")
			if tt.getSyncErr != nil {
				uploadedAt = time.Time{}.Format("2006-01-02") // default value when no previous sync
			}

			// Only proceed if no error getting sync time (or error is "not found")
			if tt.getSyncErr == nil || errors.Is(tt.getSyncErr, errors.New("not found")) {
				// Setup expectations for deleted files
				deletedParams := map[string]string{
					"deleted":     "true",
					"uploaded_at": uploadedAt,
				}
				deletedFilesJSON, _ := json.Marshal(Response{Rows: tt.deletedFiles})
				mockHTTP.On("Get", "/uploads", deletedParams).Return(deletedFilesJSON, tt.getDeletedErr).Once()

				// Process deleted files if no error
				if tt.getDeletedErr == nil {
					for i, row := range tt.deletedFiles {
						mockFileRepo.On("Exists", row.Name).Return(tt.deletedExists[i]).Once()
						if tt.deletedExists[i] {
							mockFileRepo.On("Delete", ctx, row.Name).Return(tt.deleteErr).Once()
							if tt.deleteErr != nil {
								break
							}
						}
					}

					// Only proceed if no errors so far
					if tt.deleteErr == nil {
						// Setup expectations for new files
						newParams := map[string]string{
							"deleted":     "false",
							"uploaded_at": uploadedAt,
						}
						newFilesJSON, _ := json.Marshal(Response{Rows: tt.newFiles})
						mockHTTP.On("Get", "/uploads", newParams).Return(newFilesJSON, tt.getNewErr).Once()

						if tt.getNewErr == nil {
							for i, row := range tt.newFiles {
								mockFileRepo.On("Exists", row.Name).Return(tt.newExists[i]).Once()
								if !tt.newExists[i] {
									mockFileRepo.On("GetPath", row.Name).Return("/path/to/" + row.Name).Once()
									mockHTTP.On("DownloadFile", "/file?name="+row.Name, "/path/to/"+row.Name).Return(tt.downloadErr).Once()
									if tt.downloadErr != nil {
										break
									}
								}
							}

							// Only save sync time if everything succeeded
							if tt.downloadErr == nil {
								mockTokenRepo.On("Save", ctx, username, mock.MatchedBy(func(r io.Reader) bool {
									buf := new(strings.Builder)
									io.Copy(buf, r)
									_, err := time.Parse("2006-01-02", buf.String())
									return err == nil
								})).Return(tt.saveSyncErr).Once()
							}
						}
					}
				}
			}

			err := Run(username, mockFileRepo, mockTokenRepo, mockHTTP)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockFileRepo.AssertExpectations(t)
			mockTokenRepo.AssertExpectations(t)
			mockHTTP.AssertExpectations(t)
		})
	}
}

func TestExtractRows(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []Row
		wantErr bool
	}{
		{
			name: "valid response",
			input: []byte(`{
				"pagination": {"size": 10, "page": 1, "total": 2},
				"rows": [
					{"uuid": "1", "name": "file1.txt"},
					{"uuid": "2", "name": "file2.txt"}
				]
			}`),
			want: []Row{
				{UUID: "1", Name: "file1.txt"},
				{UUID: "2", Name: "file2.txt"},
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			input:   []byte(`invalid`),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractRows(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
