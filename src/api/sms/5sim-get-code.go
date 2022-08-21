package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetCode5Sim(o Options) (string, error) {
	type CodeResp5Sim struct {
		ID        int       `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Phone     string    `json:"phone"`
		Product   string    `json:"product"`
		Price     float64   `json:"price"`
		Status    string    `json:"status"`
		Expires   time.Time `json:"expires"`
		Sms       []struct {
			ID        int       `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			Date      time.Time `json:"date"`
			Sender    string    `json:"sender"`
			Text      string    `json:"text"`
			Code      string    `json:"code"`
		} `json:"sms"`
		Forwarding       bool   `json:"forwarding"`
		ForwardingNumber string `json:"forwarding_number"`
		Country          string `json:"country"`
	}
	uri := fmt.Sprintf("https://5sim.net/v1/user/check/%s", o.ID)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}
	req.Header = http.Header{
		"Authorization": []string{"Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTAxMzU3ODYsImlhdCI6MTY1ODU5OTc4NiwicmF5IjoiOTZhNGYzZTc3OWM3ODMzNmJjNGM0YmNkODE0NWZjYTUiLCJzdWIiOjExNjk0ODJ9.xvdWH371psDfEc4nKlFxoPP0TUehKWgnfZe7GgX7Z9blyWWUdQvh0UlGgP9ke3DfAF4Q2CPQoI7635WildgdMxC1lZVnDkJE2-W-QM_L7kYyC_Aty0qWzOVM7lKsa4t_SO0pc2-KMJzB8rz-9XGLVBlfWOQQBFyASGkInsKBa7p1UnwzAM_ttIu-1znFN_dKHj5clsayUzvreIXxcYDuRHfkpEZZG2IFYpt4gD0o-HHyt86fHOzKdcyBKeUcxNKbjtpJYpaWeiheU8MpL18eQIU2dtVEhf65sCMlRcZVW8FlPWOYNHkMRdwI6Ls8aWLt7_2-tVD7E80GdmnV_f3t4Q"},
		"Accept":        []string{"application/json"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	if res.StatusCode >= 400 {
		return "", errors.New("request error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	var response CodeResp5Sim
	err = json.Unmarshal(body, &response)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}
	//log.Println(response.Status)
	if response.Sms != nil {
		log.Println(response.Sms)
		if len(response.Sms) >= 1 {
			if len(response.Sms[0].Code) <= 3 {
				if strings.Contains(response.Sms[0].Text, "Your SNKRDUNK verification code") {
					code := strings.Split(response.Sms[0].Text, "Your SNKRDUNK verification code is: ")[1]
					return code, nil
				} else {
					return "", nil
				}
			} else {
				return response.Sms[0].Code, nil
			}
		} else {
			return "", nil
		}
	} else {
		return "", nil
	}
}
