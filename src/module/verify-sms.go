package module

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"main/src/api/sms"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (t *SnkrDunkTask) verifySms() error {
	err := t.login()
	if err != nil {
		log.Println("Error logging in")
		sentry.CaptureMessage("Error logging in for sms")
		return err
	}

	o, err := t.submitSms()
	if err != nil {
		sentry.CaptureMessage("Error submitting sms")
		return err
	}
	log.Println("Submitted number")
	err = t.submitCode(o)
	if err != nil {
		sentry.CaptureMessage("Error submitting code")
		return err
	}

	return nil
}

func (t *SnkrDunkTask) login() error {
	log.Println("Logging in..")

	err := t.getLoginSession()
	if err != nil {
		sentry.CaptureMessage("Error getting login session")
		log.Fatalln("Error getting login session for sms")
	}

	form := url.Values{}
	form.Add("email", t.email)
	form.Add("password", t.password)
	form.Add("csrf_token", t.csrfLoginToken)
	form.Add("tzDatabaseName", "America/New_York")

	req, err := http.NewRequest("POST", "https://snkrdunk.com/en/login?slide=right", strings.NewReader(form.Encode()))
	if err != nil {
		log.Println("Error making request")
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
		"referer":                   []string{"https://snkrdunk.com/en/login?slide=right"},
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
		log.Println("Error sending request")
		sentry.CaptureException(err)
		return err
	}

	if res.StatusCode >= 400 {
		log.Println("Request error, status code:", res.StatusCode)
		msg := fmt.Sprintf("Request error, status code: %v", res.StatusCode)
		sentry.CaptureMessage(msg)
		return errors.New("request error")
	}
	return nil
}

func (t *SnkrDunkTask) getLoginSession() error {
	req, err := http.NewRequest("GET", "https://snkrdunk.com/en/login?slide=right", nil)
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
		sentry.CaptureMessage(er)
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
		}
		return nil
	} else {
		log.Println("No csrf tokens found")
		sentry.CaptureMessage("No csrf tokens found on sms login session")
		return errors.New("no csrf token")
	}
}

func (t *SnkrDunkTask) submitSms() (sms.Options, error) {
	type NumberSubmit struct {
		CountryCode string `json:"countryCode"`
		PhoneNumber string `json:"phoneNumber"`
	}
	var number string
	var o sms.Options
	var err error

	if len(t.Options.SmsAPIKey) > 4 {
		if strings.Contains(t.Options.SmsAPIKey, "1_") {
			number, o, err = sms.GetNewNumber2(t.Options.SmsAPIKey)
			if err != nil {
				log.Println("Error getting number:", err)
				return o, err
			}
		} else {
			number, o, err = sms.GetNewNumber()
			if err != nil {
				log.Println("Error getting number:", err)
				return o, err
			}
		}
	} else {
		number, o, err = sms.GetNumber5Sim()
		if err != nil {
			log.Println("Error getting number:", err)
			return o, err
		}
	}

	log.Println("Submitting number")
	payload, err := json.Marshal(NumberSubmit{
		CountryCode: "US",
		PhoneNumber: number,
	})
	if err != nil {
		log.Println("Error setting payload")
		sentry.CaptureException(err)
		return o, err
	}

	req, err := http.NewRequest("POST", "https://snkrdunk.com/en/v1/account/sms-verification", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error making request")
		sentry.CaptureException(err)
		return o, err
	}

	cookies := fmt.Sprintf("csrf=%s", t.csrfHeaderToken)
	req.Header = http.Header{
		"Accept":                    []string{"application/json, text/plain, */*"},
		"accept-encoding":           []string{"gzip, deflate, br"},
		"accept-language":           []string{"en-US,en;q=0.9"},
		"cache-control":             []string{"max-age=0"},
		"content-type":              []string{"application/json"},
		"cookie":                    []string{cookies},
		"dnt":                       []string{"1"},
		"origin":                    []string{"https://snkrdunk.com"},
		"referer":                   []string{"https://snkrdunk.com/en/account/phone-number?slide=right"},
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
		log.Println("Error sending request")
		sentry.CaptureException(err)
		return o, err
	}

	if res.StatusCode >= 400 {
		log.Println("Bad request, status code: ", res.StatusCode)
		if res.StatusCode == 409 {
			sentry.CaptureMessage("Number already used on site")
			log.Println("Phone number already used on site")
		} else {
			msg := fmt.Sprintf("Bad request, status code: %v", res.StatusCode)
			sentry.CaptureMessage(msg)
		}
		return o, errors.New("bad request")
	}

	return o, nil
}

func (t *SnkrDunkTask) submitCode(options sms.Options) error {
	var err error
	code := ""
	log.Println("Waiting for code..")

	if len(t.Options.SmsAPIKey) > 4 {
		if strings.Contains(t.Options.SmsAPIKey, "1_") {
			for code == "" {
				code, err = sms.GetCode2(options)
				if err != nil {
					log.Println("Error getting code")
					return err
				}
				time.Sleep(2500 * time.Millisecond)
			}
		} else {
			for code == "" {
				code, err = sms.GetCode(options)
				if err != nil {
					log.Println("Error getting code")
					return err
				}
				time.Sleep(2000 * time.Millisecond)
			}
		}
	} else {
		for code == "" {
			code, err = sms.GetCode5Sim(options)
			if err != nil {
				log.Println("Error getting code")
				return err
			}
			time.Sleep(2000 * time.Millisecond)
		}
	}

	log.Println("Got code, submitting to snkr dunk", code)
	type VerificationPayload struct {
		PinCode string `json:"pinCode"`
	}

	payload, err := json.Marshal(VerificationPayload{PinCode: code})
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	req, err := http.NewRequest("PATCH", "https://snkrdunk.com/en/v1/account/sms-verification", bytes.NewBuffer(payload))
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	cookies := fmt.Sprintf("csrf=%s", t.csrfHeaderToken)
	req.Header = http.Header{
		"Accept":             []string{"application/json, text/plain, */*"},
		"accept-encoding":    []string{"gzip, deflate, br"},
		"accept-language":    []string{"en-US,en;q=0.9"},
		"cache-control":      []string{"max-age=0"},
		"content-type":       []string{"application/json"},
		"cookie":             []string{cookies},
		"dnt":                []string{"1"},
		"origin":             []string{"https://snkrdunk.com"},
		"referer":            []string{"https://snkrdunk.com/en/account/phone-number?slide=right"},
		"sec-ch-ua":          []string{"\".Not/A)Brand\";v=\"99\", \"Google Chrome\";v=\"103\", \"Chromium\";v=\"103\""},
		"sec-ch-ua-mobile":   []string{"?0"},
		"sec-ch-ua-platform": []string{"\"Windows\""},
		"sec-fetch-dest":     []string{"empty"},
		"sec-fetch-mode":     []string{"cors"},
		"sec-fetch-site":     []string{"same-origin"},
		"user-agent":         []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36"},
	}

	res, err := t.client.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	if res.StatusCode != 200 {
		log.Println("Request error sending code: ", res.StatusCode)
		msg := fmt.Sprintf("Request error sending code: %v", res.StatusCode)
		sentry.CaptureMessage(msg)
		return errors.New("request error sending code")
	}

	return nil
}
