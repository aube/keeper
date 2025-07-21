package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
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

	// Читаем тело ответа
	return io.ReadAll(resp.Body)
}

// DownloadFile скачивает файл по URL
func (c *HTTPClient) DownloadFile(fileURL, outputPath string) error {
	// Создаем запрос
	req, err := http.NewRequest("GET", fileURL, nil)
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
