package register

type HTTPClient interface {
	SetHeader(key, value string)
	Get(endpoint string, queryParams map[string]string) ([]byte, error)
	Post(endpoint string, body any) ([]byte, error)
	DownloadFile(fileURL, outputPath string) error
}

func Run(username string, password string, email string, http HTTPClient) error {

	postData := map[string]interface{}{
		"username": username,
		"password": password,
		"email":    email,
	}
	_, err := http.Post("/register", postData)
	if err != nil {
		return err
	}

	return nil
}
