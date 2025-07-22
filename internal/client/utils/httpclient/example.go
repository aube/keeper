package httpclient

import (
	"fmt"
	"log"
)

func main() {
	// Пример использования клиента

	// Создаем клиент
	client := NewHTTPClient("https://jsonplaceholder.typicode.com")

	// Устанавливаем заголовки
	client.SetHeader("User-Agent", "MyGoHTTPClient/1.0")
	client.SetHeader("Accept", "application/json")

	// Пример GET-запроса
	fmt.Println("Making GET request...")
	response, err := client.Get("/todos/1", nil)
	if err != nil {
		log.Fatalf("GET error: %v", err)
	}
	fmt.Printf("GET Response: %s\n", response)

	// Пример POST-запроса
	fmt.Println("\nMaking POST request...")
	postData := map[string]interface{}{
		"title":  "foo",
		"body":   "bar",
		"userId": 1,
	}
	response, err = client.Post("/posts", postData)
	if err != nil {
		log.Fatalf("POST error: %v", err)
	}
	fmt.Printf("POST Response: %s\n", response)

	// Пример загрузки файла
	fmt.Println("\nDownloading file...")
	fileURL := "https://example.com/sample.jpg"
	err = client.DownloadFile(fileURL, "sample.jpg")
	if err != nil {
		log.Printf("Download error: %v (this might fail if the URL doesn't point to a real file)", err)
	} else {
		fmt.Println("File downloaded successfully")
	}
}
