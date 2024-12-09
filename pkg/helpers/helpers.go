package helpers

import (
	"net/http"
	"time"
)

func CreateHttpClient() *http.Client {
	var httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
	return httpClient
}
