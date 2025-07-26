package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/aube/keeper/internal/client/utils/progress"
	"github.com/schollz/progressbar/v3"
)

// HTTPClient представляет наш клиент
type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
}

// NewHTTPClient создает новый экземпляр HTTPClient
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		headers: make(map[string]string),
	}
}

// SetHeader устанавливает заголовок для всех последующих запросов
func (c *HTTPClient) SetHeader(key, value string) {
	c.headers[key] = value
}

// Get выполняет GET-запрос
func (c *HTTPClient) Get(endpoint string, queryParams map[string]string) ([]byte, error) {
	// Создаем URL
	u, err := url.Parse(c.baseURL + endpoint)
	if err != nil {
		return nil, err
	}

	// Добавляем параметры запроса
	if queryParams != nil {
		q := u.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		u.RawQuery = q.Encode()
	}

	// Создаем запрос
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Добавляем заголовки
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	return io.ReadAll(resp.Body)
}

// Post выполняет POST-запрос с JSON телом
func (c *HTTPClient) Post(endpoint string, body interface{}) ([]byte, error) {
	// Преобразуем тело в JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Добавляем заголовки
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	return io.ReadAll(resp.Body)
}

// DownloadFile скачивает файл по URL
func (c *HTTPClient) DownloadFile(endpoint, outputPath string) error {
	// Создаем запрос
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return err
	}

	// Добавляем заголовки
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Создаем файл
	out, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Копируем данные в файл
	_, err = io.Copy(out, resp.Body)
	return err
}

// UploadFileWithProgress отправляет файл с отображением прогресса
func (c *HTTPClient) UploadFile(ctx context.Context, endpoint string, filePath string, formFields map[string]string) ([]byte, error) {
	// Получаем информацию о файле для прогресс-бара
	fi, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения информации о файле: %w", err)
	}

	// Создаем pipe для отслеживания прогресса
	pr, pw := io.Pipe()
	defer pr.Close()

	// Создаем прогресс-бар
	bar := progress.NewBar(fi.Size(), "Отправка файла...")

	// Горутина для записи файла в pipe с отслеживанием прогресса
	go func() {
		defer pw.Close()

		// Открываем файл
		file, err := os.Open(filePath)
		if err != nil {
			pw.CloseWithError(fmt.Errorf("ошибка открытия файла: %w", err))
			return
		}
		defer file.Close()

		// Копируем файл в pipe с отслеживанием прогресса
		reader := progressbar.NewReader(file, bar)
		_, err = io.Copy(pw, &reader)
		if err != nil {
			pw.CloseWithError(fmt.Errorf("ошибка копирования файла: %w", err))
			return
		}
	}()

	// Создаем multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файл
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания части файла: %w", err)
	}

	// Копируем данные из pipe в multipart
	if _, err = io.Copy(part, pr); err != nil {
		return nil, fmt.Errorf("ошибка копирования данных: %w", err)
	}

	// Добавляем текстовые поля
	for key, value := range formFields {
		_ = writer.WriteField(key, value)
	}

	// Закрываем writer
	if err = writer.Close(); err != nil {
		return nil, fmt.Errorf("ошибка закрытия writer: %w", err)
	}

	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	// Устанавливаем заголовки
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус код
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("неожиданный статус код: %d", resp.StatusCode)
	}

	// Читаем ответ
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	return responseBody, nil
}
