package gitearelease

import (
	"errors"
	"io"
	"net/http"
)

func fetchData(URL string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		ourErr := errors.New("Net: Failed to create a new request")
		return nil, ourErr
	}

	resp, err := client.Do(req)
	if err != nil {
		ourErr := errors.New("Net: Failed to create a new request")
		return nil, ourErr
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ourErr := errors.New("Net: Server response: " + resp.Status)
		return nil, ourErr
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ourErr := errors.New("IO: Failed to read the response body")
		return nil, ourErr
	}

	return body, nil
}
