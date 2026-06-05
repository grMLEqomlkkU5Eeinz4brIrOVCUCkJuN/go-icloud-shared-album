package utils

import (
	"bytes"
	"context"
	"fmt"
	"io" // Changed from "io/ioutil"
	"log" // Added log import
	"net/http"
)

func GetRedirectedBaseUrl(baseUrl string, token string) (string, error) {
	url := baseUrl + "webstream"
	dataString := []byte(`{"streamCtag":null}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataString))
	if err != nil {
		return "", err
	}

	req.Header.Set("Origin", "https://www.icloud.com")
	req.Header.Set("Accept-Language", "en-US,en;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "https://www.icloud.com/sharedalbum/")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Don't follow redirects automatically
		},
	}
	resp, err := client.Do(req.WithContext(context.Background()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 330 {
		host := resp.Header.Get("X-Apple-MMe-Host")
		if host == "" {
			return "", fmt.Errorf("missing X-Apple-MMe-Host in 330 response")
		}

		newBaseUrl := fmt.Sprintf(
			"https://%s/%s/sharedstreams/",
			host,
			token,
		)
		return newBaseUrl, nil
	}

	// If not 330, return the original baseUrl.
	// In some cases, the initial baseUrl might be correct, or it's a different error.
	if resp.StatusCode != http.StatusOK {
		// Attempt to read body for more error context if it's not a 330 redirect and not OK
		bodyBytes, _ := io.ReadAll(resp.Body) // Changed ioutil.ReadAll to io.ReadAll
		log.Printf("GetRedirectedBaseUrl received non-OK status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}
	
	return baseUrl, nil
}
