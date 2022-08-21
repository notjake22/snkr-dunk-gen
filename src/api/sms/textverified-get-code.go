package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetCode2(options Options) (string, error) {
	uri := fmt.Sprintf("https://www.textverified.com/api/Verifications/%s", options.ID)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		sentry.CaptureException(err)
		return "", nil
	}
	req.Header.Set("Authorization", options.BearerHeader)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return "", nil
	}

	if res.StatusCode >= 400 {
		msg := fmt.Sprintf("Request error getting code TextVerified: %v", res.StatusCode)
		sentry.CaptureMessage(msg)
		return "", errors.New("request error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sentry.CaptureException(err)
		return "", nil
	}

	var response NewNumberResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		sentry.CaptureException(err)
		return "", nil
	}

	if len(response.Code) <= 2 {
		if strings.Contains(response.Sms, "Your SNKRDUNK verification code") {
			code := strings.Split(response.Sms, "Your SNKRDUNK verification code is: ")[1]
			return code, nil
		}
	} else {
		if strings.Contains(response.Sms, "Your SNKRDUNK verification code") {
			code := strings.Split(response.Sms, "Your SNKRDUNK verification code is: ")[1]
			return code, nil
		}
	}
	return "", nil
}
