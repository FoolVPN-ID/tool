package common

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

func MakeHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func RecoverFromPanic() any {
	if err := recover(); err != nil {
		log.Printf("Recovered from panic: %v\n", err)
		return err
	}

	return ""
}
