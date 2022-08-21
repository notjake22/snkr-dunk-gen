package module

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
)

func (t *SnkrDunkTask) verifyEmailReq(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error making request to verify email")
		sentry.CaptureException(err)
		return err
	}

	req.Header = http.Header{
		"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"accept-encoding":           []string{"gzip, deflate, br"},
		"accept-language":           []string{"en-US,en;q=0.9"},
		"dnt":                       []string{"1"},
		"sec-ch-ua":                 []string{"\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"103\", \"Chromium\";v=\"103\""},
		"sec-ch-ua-mobile":          []string{"?0"},
		"sec-ch-ua-platform":        []string{"\"Windows\""},
		"sec-fetch-dest":            []string{"document"},
		"sec-fetch-mode":            []string{"navigate"},
		"sec-fetch-site":            []string{"same-origin"},
		"upgrade-insecure-requests": []string{"1"},
		"user-agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"},
	}

	res, err := t.client.Do(req)
	if err != nil {
		log.Println("Error sending request to verify email")
		sentry.CaptureException(err)
		return err
	}

	if res.StatusCode != 200 {
		log.Println("Request error, status code: ", res.StatusCode)
		msg := fmt.Sprintf("Request error, status code: %v", res.StatusCode)
		sentry.CaptureMessage(msg)
		return errors.New("request error verifying email")
	}

	defer res.Body.Close()
	for _, cookie := range res.Cookies() {
		if cookie.Name == "csrf" {
			t.csrfHeaderToken = cookie.Value
		}
		if cookie.Name == "ENSID" {
			t.EnsID = cookie.Value
		}
	}

	return nil
}
