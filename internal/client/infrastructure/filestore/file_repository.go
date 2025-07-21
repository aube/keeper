package filestore

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/aube/keeper/internal/client/entities"
	"github.com/aube/keeper/internal/client/utils/apperrors"
	"github.com/aube/keeper/internal/client/utils/logger"
	"github.com/rs/zerolog"
)

type FileSystemRepository struct {
	storagePath string
	mu          sync.RWMutex
	log         zerolog.Logger
}

func NewFileSystemRepository(storagePath string) (*FileSystemRepository, error) {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, err
	}
	return &FileSystemRepository{
		storagePath: storagePath,
		log:         logger.Get().With().Str("fs", "file_repository").Logger(),
	}, nil
}

func (r *FileSystemRepository) Save(ctx context.Context, file *entities.File, data io.Reader) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := filepath.Join(r.storagePath, file.Name)
	dst, err := os.Create(filePath)
	if err != nil {
		r.log.Debug().Err(err).Msg("Save")
		r.log.Debug().Msg(filePath)
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, data); err != nil {
		r.log.Debug().Err(err).Msg("Save")
		return err
	}

	return nil
}

func (r *FileSystemRepository) FindAll(ctx context.Context) (*entities.Files, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	files, err := os.ReadDir(r.storagePath)
	if err != nil {
		r.log.Debug().Err(err).Msg("FindAll")
		return nil, err
	}

	var result entities.Files
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		result = append(result, *entities.NewFile(
			file.Name(),
			filepath.Join(r.storagePath, file.Name()),
			fileInfo.Size(),
		))
	}

	return &result, nil
}

func (r *FileSystemRepository) Delete(ctx context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := filepath.Join(r.storagePath, uuid)
	if err := os.Remove(filePath); err != nil {
		r.log.Debug().Err(err).Msg("Delete")
		if os.IsNotExist(err) {
			return apperrors.ErrFileNotFound
		}
		return err
	}
	return nil
}

func (r *FileSystemRepository) GetFileContent(ctx context.Context, uuid string) (io.ReadCloser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filePath := filepath.Join(r.storagePath, uuid)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, apperrors.ErrFileNotFound
	}

	return os.Open(filePath)
}

// var _ appFile.FileRepository = (*FileSystemRepository)(nil)
