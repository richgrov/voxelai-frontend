package app

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
)

type mockConstructorService struct{}

func (*mockConstructorService) build(id string, prompt string) (string, error) {
	if rand.Intn(10) == 0 {
		return "", errors.New("fake error")
	}

	return "/schem/stone.schematic", nil
}

type httpConstructorService struct {
	url string
}

func (service *httpConstructorService) build(id string, prompt string) (string, error) {
	req, err := http.NewRequest("POST", service.url, nil)
	if err != nil {
		return "", err
	}

	query := req.URL.Query()
	query.Add("id", id)
	query.Add("prompt", prompt)
	req.URL.RawQuery = query.Encode()
	req.URL.Path = "/generate"

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", errors.New(string(result))
	}

	return string(result), nil
}
