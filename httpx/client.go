package httpx

import (
	"net/http"
	"time"
)

func ClientSeconds(seconds int) *http.Client {
	return &http.Client{
		Timeout: time.Duration(seconds) * time.Second,
	}
}
