package filestore

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/aube/keeper/internal/client/entities"
	"github.com/aube/keeper/internal/client/utils/apperrors"
	"github.com/aube/keeper/internal/client/utils/logger"
	"github.com/aube/keeper/internal/client/utils/progress"
	"github.com/rs/zerolog"
)

const (
	chunkSize    = 4096 // Размер блока для чтения/шифрования
	gcmNonceSize = 12   // Стандартный размер nonce для GCM
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

func (r *FileSystemRepository) GetPath(filename string) string {
	return filepath.Join(r.storagePath, filename)
}

func (r *FileSystemRepository) Save(ctx context.Context, filename string, data io.Reader) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := r.GetPath(filename)
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

func (r *FileSystemRepository) Delete(ctx context.Context, filename string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := r.GetPath(filename)
	if err := os.Remove(filePath); err != nil {
		r.log.Debug().Err(err).Msg("Delete")
		if os.IsNotExist(err) {
			return apperrors.ErrFileNotFound
		}
		return err
	}
	return nil
}

func (r *FileSystemRepository) GetFile(ctx context.Context, filename string) (io.ReadCloser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filePath := filepath.Join(r.storagePath, filename)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, apperrors.ErrFileNotFound
	}

	return os.Open(filePath)
}

func (r *FileSystemRepository) GetFileContent(ctx context.Context, filename string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	file, err := r.GetFile(ctx, filename)
	if err != nil {
		return "", err
	}

	bodyBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
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

func (r *FileSystemRepository) EncryptFile(inputPath, outputName, password string) error {

	fi, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	bar := progress.NewBar(fi.Size(), "Шифрую файл...")

	r.mu.RLock()
	defer r.mu.RUnlock()

	outputPath := r.GetPath(outputName)

	// Открываем исходный файл
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть входной файл: %w", err)
	}
	defer inputFile.Close()

	// Создаем выходной файл
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("не удалось создать выходной файл: %w", err)
	}
	defer outputFile.Close()

	// Генерируем ключ из пароля
	key := deriveKey(password)

	// Инициализируем шифр
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("ошибка создания блока шифра: %w", err)
	}

	// Генерируем уникальный nonce
	nonce := make([]byte, gcmNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return fmt.Errorf("ошибка генерации nonce: %w", err)
	}

	// Записываем nonce в начало выходного файла
	if _, err := outputFile.Write(nonce); err != nil {
		return fmt.Errorf("ошибка записи nonce: %w", err)
	}

	// Создаем GCM режим
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("ошибка создания GCM: %w", err)
	}

	// Буфер для чтения данных
	buf := make([]byte, chunkSize)

	for {
		// Читаем порцию данных из файла
		n, err := inputFile.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("ошибка чтения файла: %w", err)
		}

		if n == 0 {
			break
		}

		// Шифруем данные
		ciphertext := gcm.Seal(nil, nonce, buf[:n], nil)

		// Записываем зашифрованные данные
		if _, err := outputFile.Write(ciphertext); err != nil {
			return fmt.Errorf("ошибка записи зашифрованных данных: %w", err)
		}

		// Обновляем прогресс-бар
		if err := bar.Add(n); err != nil {
			return fmt.Errorf("ошибка обновления прогресс-бара: %w", err)
		}
	}

	return nil
}

func (r *FileSystemRepository) DecryptFile(inputName, outputPath, password string) error {
	inputPath := r.GetPath(inputName)

	fi, err := os.Stat(inputPath)
	if err != nil {
		return err
	}
	bar := progress.NewBar(fi.Size(), "Дешифрую файл...")

	// Открываем зашифрованный файл
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("не удалось открыть входной файл: %w", err)
	}
	defer inputFile.Close()

	// Создаем файл для расшифрованных данных
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("не удалось создать выходной файл: %w", err)
	}
	defer outputFile.Close()

	// Получаем ключ из пароля
	key := deriveKey(password)

	// Инициализируем AES блок
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("ошибка создания блока шифра: %w", err)
	}

	// Читаем nonce из начала файла
	nonce := make([]byte, gcmNonceSize)
	if _, err := io.ReadFull(inputFile, nonce); err != nil {
		return fmt.Errorf("ошибка чтения nonce: %w", err)
	}

	// Создаем GCM режим
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("ошибка создания GCM: %w", err)
	}

	// Буфер для чтения данных
	buf := make([]byte, chunkSize+gcm.Overhead()) // Учитываем overhead аутентификации

	for {
		// Читаем порцию зашифрованных данных
		n, err := inputFile.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("ошибка чтения файла: %w", err)
		}

		if n == 0 {
			break
		}

		// Расшифровываем данные
		plaintext, err := gcm.Open(nil, nonce, buf[:n], nil)
		if err != nil {
			return fmt.Errorf("ошибка дешифрования: %w", err)
		}

		// Записываем расшифрованные данные
		if _, err := outputFile.Write(plaintext); err != nil {
			return fmt.Errorf("ошибка записи данных: %w", err)
		}

		// Обновляем прогресс-бар
		if err := bar.Add(n); err != nil {
			return fmt.Errorf("ошибка обновления прогресс-бара: %w", err)
		}
	}

	return nil
}

// var _ appFile.FileRepository = (*FileSystemRepository)(nil)

func deriveKey(password string) []byte {
	// В реальном приложении используйте PBKDF2, scrypt или аналогичные функции
	// для безопасного преобразования пароля в ключ
	key := make([]byte, 32) // AES-256 требует 32-байтный ключ
	copy(key, password)

	// Если пароль короче 32 байт, оставшиеся байты останутся нулевыми
	// Если длиннее - обрежем
	if len(key) > 32 {
		key = key[:32]
	}

	return key
}
