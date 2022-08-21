package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetNumber5Sim() (string, Options, error) {
	type Number5Sim struct {
		ID               int       `json:"id"`
		Phone            string    `json:"phone"`
		Operator         string    `json:"operator"`
		Product          string    `json:"product"`
		Price            float64   `json:"price"`
		Status           string    `json:"status"`
		Expires          time.Time `json:"expires"`
		Sms              string    `json:"sms"`
		CreatedAt        time.Time `json:"created_at"`
		Forwarding       bool      `json:"forwarding"`
		ForwardingNumber string    `json:"forwarding_number"`
		Country          string    `json:"country"`
	}
	o := Options{
		ID: "",
	}
	status, err := checkProductStatus()
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}
	if !status {
		sentry.CaptureMessage("Numbers out of stock on 5sim")
		return "", o, errors.New("numbers out of stock")
	}

	req, err := http.NewRequest("GET", "https://5sim.net/v1/user/buy/activation/usa/any/other", nil)
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}

	req.Header = http.Header{
		"Authorization": []string{"Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTAxMzU3ODYsImlhdCI6MTY1ODU5OTc4NiwicmF5IjoiOTZhNGYzZTc3OWM3ODMzNmJjNGM0YmNkODE0NWZjYTUiLCJzdWIiOjExNjk0ODJ9.xvdWH371psDfEc4nKlFxoPP0TUehKWgnfZe7GgX7Z9blyWWUdQvh0UlGgP9ke3DfAF4Q2CPQoI7635WildgdMxC1lZVnDkJE2-W-QM_L7kYyC_Aty0qWzOVM7lKsa4t_SO0pc2-KMJzB8rz-9XGLVBlfWOQQBFyASGkInsKBa7p1UnwzAM_ttIu-1znFN_dKHj5clsayUzvreIXxcYDuRHfkpEZZG2IFYpt4gD0o-HHyt86fHOzKdcyBKeUcxNKbjtpJYpaWeiheU8MpL18eQIU2dtVEhf65sCMlRcZVW8FlPWOYNHkMRdwI6Ls8aWLt7_2-tVD7E80GdmnV_f3t4Q"},
		"Accept":        []string{"application/json"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}

	if res.StatusCode >= 400 {
		msg := fmt.Sprintf("Request error getting number 5sim: %v", res.StatusCode)
		sentry.CaptureMessage(msg)
		return "", o, errors.New("request error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}

	var response Number5Sim
	err = json.Unmarshal(body, &response)
	if err != nil {
		sentry.CaptureException(err)
		return "", o, err
	}

	o.ID = strconv.Itoa(response.ID)
	return strings.Split(response.Phone, "+1")[1], o, nil
}

func checkProductStatus() (bool, error) {
	type ProductStatus5Sim struct {
		Usa struct {
			Other struct {
				Virtual8 struct {
					Cost  float64 `json:"cost"`
					Count int     `json:"count"`
				} `json:"virtual8"`
			} `json:"other"`
		} `json:"usa"`
	}
	req, err := http.NewRequest("GET", "https://5sim.net/v1/guest/prices?country=usa&product=other", nil)
	if err != nil {
		sentry.CaptureException(err)
		return false, err
	}
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	if res.StatusCode >= 400 {
		msg := fmt.Sprintf("Request error checking code 5sim: %v", res.StatusCode)
		sentry.CaptureMessage(msg)
		return false, errors.New("request error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var response ProductStatus5Sim
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	if response.Usa.Other.Virtual8.Count <= 1 {
		return false, nil
	}
	if response.Usa.Other.Virtual8.Cost >= 10.0 {
		return false, nil
	}

	return true, nil
}
