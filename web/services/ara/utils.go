package ara

import (
	"io"
	"net/http"
)

//go:generate mockery --name=GetJson --output mocksutils

type GetJson func(query string) ([]byte, error)

var getJson GetJson = func(query string) ([]byte, error) {
	var err error
	resp, err := http.Get(query)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
