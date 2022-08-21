package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func GetNewNumber2(apiKey string) (string, Options, error) {
	o := Options{
		BearerHeader: "",
	}

	bearer, err := authApi(apiKey)
	if err != nil {
		return "", o, err
	}
	time.Sleep(500 * time.Millisecond)
	status, err := checkStatus(bearer)
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}
	if !status {
		return "", o, errors.New("api error")
	}
	time.Sleep(500 * time.Millisecond)

	type GetNumberPayload struct {
		ID                     int32  `json:"id"`
		AreaCode               string `json:"areaCode"`
		RequestedTimeAllotment string `json:"requestedTimeAllotment"`
	}
	payload, err := json.Marshal(GetNumberPayload{
		ID:                     0,
		RequestedTimeAllotment: "Not implemented",
	})
	req, err := http.NewRequest("POST", "https://www.textverified.com/api/Verifications", bytes.NewBuffer(payload))
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}

	bearerHeader := fmt.Sprintf("Bearer %s", bearer)
	req.Header = http.Header{
		"Accept":        []string{"application/json"},
		"Authorization": []string{bearerHeader},
		"Content-Type":  []string{"application/json"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}

	if res.StatusCode >= 400 {
		er := fmt.Sprintf("request error 1 %+v", res.StatusCode)
		sentry.CaptureMessage(er)
		if res.StatusCode == 402 {
			sentry.CaptureMessage("User too low of balance")
		}
		return "", o, errors.New(er)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}

	var response NewNumberResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}
	//if response.Status != "Completed" {
	//	return "", o, errors.New(response.Status)
	//}

	o.BearerHeader = bearerHeader
	o.ID = response.ID
	return response.Number, o, nil
}

func authApi(apiKey string) (string, error) {
	type Auth struct {
		BearerToken string `json:"bearer_token"`
	}
	req, err := http.NewRequest("POST", "https://www.textverified.com/Api/SimpleAuthentication", nil)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	req.Header.Set("X-SIMPLE-API-ACCESS-TOKEN", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	if res.StatusCode >= 400 {
		er := fmt.Sprintf("request error 2 %+v", res.StatusCode)
		sentry.CaptureException(errors.New(er))
		return "", errors.New(er)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	var response Auth
	err = json.Unmarshal(body, &response)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	return response.BearerToken, nil
}

func checkStatus(bearer string) (bool, error) {
	type Status struct {
		TargetID       int     `json:"targetId"`
		Name           string  `json:"name"`
		NormalizedName string  `json:"normalizedName"`
		Cost           float64 `json:"cost"`
		Status         int     `json:"status"`
		PricingMode    int     `json:"pricingMode"`
	}
	req, err := http.NewRequest("GET", "https://www.textverified.com/api/Targets/0", nil)
	if err != nil {
		sentry.CaptureException(err)
		return false, err
	}

	bearerHeader := fmt.Sprintf("Bearer %s", bearer)
	req.Header.Set("Authorization", bearerHeader)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return false, err
	}

	if res.StatusCode >= 400 {
		er := fmt.Sprintf("request error 3 %+v", res.StatusCode)
		sentry.CaptureMessage(er)
		return false, errors.New(er)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sentry.CaptureException(err)
		return false, err
	}

	var response Status
	err = json.Unmarshal(body, &response)
	if err != nil {
		sentry.CaptureException(err)
		return false, err
	}

	if response.Status != 4 {
		if response.Status == 1 {
			log.Println("Out of numbers")
			sentry.CaptureMessage("TextVerified status 1, out of numbers")
			return false, errors.New("out of numbers")
		}
		if response.Status == 128 {
			sentry.CaptureMessage("TextVerified status 128, quota exceeded")
			return false, errors.New("quota exceeded")
		} else {
			msg := fmt.Sprintf("TextVerified api error, status code: %v", response.Status)
			sentry.CaptureMessage(msg)
			return false, errors.New("text verified api error")
		}
	}

	return true, nil
}

type NewNumberResponse struct {
	ID              string  `json:"id"`
	Cost            float64 `json:"cost"`
	TargetName      string  `json:"target_name"`
	Number          string  `json:"number"`
	SenderNumber    string  `json:"sender_number"`
	TimeRemaining   string  `json:"time_remaining"`
	ReuseWindow     string  `json:"reuse_window"`
	Status          string  `json:"status"`
	Sms             string  `json:"sms"`
	Code            string  `json:"code"`
	VerificationURI string  `json:"verification_uri"`
	CancelURI       string  `json:"cancel_uri"`
	ReportURI       string  `json:"report_uri"`
	ReuseURI        string  `json:"reuse_uri"`
}
