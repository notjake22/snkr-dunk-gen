package module

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func (t *SnkrDunkTask) addAddress() error {
	csrfToken, err := t.getAddressSession()
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	form := url.Values{}
	form.Add("firstName", "Mike")
	form.Add("lastName", "Smith")
	form.Add("phoneNumber", "8455789804")
	form.Add("country", "US")
	form.Add("streetAddress", "1150 N Damen Ave")
	form.Add("aptSuite", "")
	form.Add("city", "Chicago")
	form.Add("region", "IL")
	form.Add("postCode", "60622")
	form.Add("csrf_token", csrfToken)

	req, err := http.NewRequest("POST", "https://snkrdunk.com/en/account/address?slide=right", strings.NewReader(form.Encode()))
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	cookies := fmt.Sprintf("csrf=%s", t.csrfHeaderToken)
	req.Header = http.Header{
		"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"accept-encoding":           []string{"gzip, deflate, br"},
		"accept-language":           []string{"en-US,en;q=0.9"},
		"cache-control":             []string{"max-age=0"},
		"content-type":              []string{"application/x-www-form-urlencoded"},
		"cookie":                    []string{cookies},
		"dnt":                       []string{"1"},
		"origin":                    []string{"https://snkrdunk.com"},
		"referer":                   []string{"https://snkrdunk.com/en/account/address?slide=right"},
		"sec-ch-ua":                 []string{"\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"103\", \"Chromium\";v=\"103\""},
		"sec-ch-ua-mobile":          []string{"?0"},
		"sec-ch-ua-platform":        []string{"\"Windows\""},
		"sec-fetch-dest":            []string{"document"},
		"sec-fetch-mode":            []string{"navigate"},
		"sec-fetch-site":            []string{"same-origin"},
		"sec-fetch-user":            []string{"?1"},
		"upgrade-insecure-requests": []string{"1"},
		"user-agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"},
	}

	res, err := t.client.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	if res.StatusCode >= 400 {
		log.Println("Request error, code: ", res.StatusCode)
		msg := fmt.Sprintf("Request error, code: %v", res.StatusCode)
		sentry.CaptureMessage(msg)
		return errors.New("request error")
	}

	return nil
}

func (t *SnkrDunkTask) getAddressSession() (string, error) {
	req, err := http.NewRequest("GET", "https://snkrdunk.com/en/account/address?slide=right", nil)
	if err != nil {
		log.Println("Error starting request")
		sentry.CaptureException(err)
		return "", err
	}

	cookies := fmt.Sprintf("csrf=%s", t.csrfHeaderToken)
	req.Header = http.Header{
		"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
		"accept-encoding":           []string{"gzip, deflate, br"},
		"accept-language":           []string{"en-US,en;q=0.9"},
		"cookie":                    []string{cookies},
		"dnt":                       []string{"1"},
		"referer":                   []string{"https://snkrdunk.com/en/account"},
		"sec-ch-ua":                 []string{"\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"103\", \"Chromium\";v=\"103\""},
		"sec-ch-ua-mobile":          []string{"?0"},
		"sec-ch-ua-platform":        []string{"\"Windows\""},
		"sec-fetch-dest":            []string{"document"},
		"sec-fetch-mode":            []string{"navigate"},
		"sec-fetch-site":            []string{"same-origin"},
		"sec-fetch-user":            []string{"?1"},
		"upgrade-insecure-requests": []string{"1"},
		"user-agent":                []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"},
	}

	res, err := t.client.Do(req)
	if err != nil {
		log.Println("Error making session request to site")
		sentry.CaptureException(err)
		return "", err
	}

	if res.StatusCode != 200 {
		log.Println("Request error to site starting session")
		er := fmt.Sprintf("Request error: $%v", res.StatusCode)
		sentry.CaptureMessage(er)
		return "", errors.New(er)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading body")
		sentry.CaptureException(err)
		return "", err
	}

	findCsrf := strings.Split(string(body), "csrf_token")
	if len(findCsrf) > 0 {
		csrf := strings.Split(findCsrf[1], "\"")
		return csrf[2], nil
	} else {
		log.Println("No csrf tokens found")
		sentry.CaptureMessage("No csrf token found on submitting address session")
		return "", errors.New("no csrf token")
	}
}
