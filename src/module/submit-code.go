package module

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"main/src/webhooks"
	"net/http"
)

func (t *SnkrDunkTask) submitReferralCode() error {
	type ReferralPayload struct {
		Code string `json:"code"`
	}
	payload, err := json.Marshal(ReferralPayload{Code: t.Options.ReferralCode})
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	req, err := http.NewRequest("POST", "https://snkrdunk.com/en/v1/invitation", bytes.NewBuffer(payload))
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	cookies := fmt.Sprintf("csrf=%s", t.csrfHeaderToken)
	req.Header = http.Header{
		"Accept":                    []string{"application/json, text/plain, */*"},
		"accept-encoding":           []string{"gzip, deflate, br"},
		"accept-language":           []string{"en-US,en;q=0.9"},
		"content-type":              []string{"application/json"},
		"cookie":                    []string{cookies},
		"dnt":                       []string{"1"},
		"origin":                    []string{"https://snkrdunk.com"},
		"referer":                   []string{"https://snkrdunk.com/en/account/coupons?slide=right"},
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
		sentry.CaptureException(err)
		return err
	}

	if res.StatusCode >= 400 {
		log.Println("Request error, status: ", res.StatusCode)
		msg := fmt.Sprintf("Request error, status: %v", res.StatusCode)
		sentry.CaptureMessage(msg)
		body, _ := ioutil.ReadAll(res.Body)
		bodyMsg := fmt.Sprintf("Error response body: %s \n referal code attempting to be used: %s", string(body), t.Options.ReferralCode)
		sentry.CaptureMessage(bodyMsg)
		return errors.New("request error")
	}
	webhooks.ReferralCodeSubmitted(t.email, t.password, t.Options.TaskID)

	return nil
}
