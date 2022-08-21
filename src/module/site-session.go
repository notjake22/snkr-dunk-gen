package module

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func (t *SnkrDunkTask) getSession() error {
	req, err := http.NewRequest("GET", "https://snkrdunk.com/en/signup", nil)
	if err != nil {
		log.Println("Error starting request")
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
		log.Println("Error making session request to site")
		sentry.CaptureException(err)
		return err
	}

	if res.StatusCode != 200 {
		log.Println("Request error to site starting session")
		er := fmt.Sprintf("Request error: $%v", res.StatusCode)
		sentry.CaptureException(errors.New(er))
		return errors.New(er)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body")
		sentry.CaptureException(err)
		return err
	}

	findCsrf := strings.Split(string(body), "csrf_token")
	if len(findCsrf) > 0 {
		csrf := strings.Split(findCsrf[1], "\"")
		t.csrfLoginToken = csrf[2]

		defer res.Body.Close()
		for _, cookie := range res.Cookies() {
			if cookie.Name == "csrf" {
				t.csrfHeaderToken = cookie.Value
			}
			//if cookie.Name == "ENSID" {
			//	t.EnsID = cookie.Value
			//}
		}
		return nil
	} else {
		log.Println("No csrf tokens found")
		sentry.CaptureMessage("No csrf token found in initial site session for account creation")
		return errors.New("no csrf token")
	}
}
