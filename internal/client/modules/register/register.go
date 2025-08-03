package register

type HTTPClient interface {
	Post(endpoint string, body any) ([]byte, error)
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
