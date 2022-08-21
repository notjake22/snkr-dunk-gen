package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"main/src/settings"
	"net/http"
	"strings"
)

type CodeResponse struct {
	Response string `json:"response"`
	Number   string `json:"number"`
	Sms      string `json:"sms"`
	Text     string `json:"text"`
}

func GetCode(options Options) (string, error) {
	uri := fmt.Sprintf("https://smspva.com/priemnik.php?metod=get_sms&country=us&service=opt%s&id=%s&apikey=%s", options.OptionNumber, options.ID, settings.ReadSettings().SmsAPIKey)
	//log.Println("Url: ", uri)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	if res.StatusCode != 200 {
		msg := fmt.Sprintf("Request error getting code SmsPva: %v", res.StatusCode)
		sentry.CaptureMessage(msg)
		return "", errors.New("request error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sentry.CaptureException(err)
		return "", err
	}

	var response CodeResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("%+v \n", string(body))
		sentry.CaptureException(err)
		bodyResp := fmt.Sprintf("%+v", string(body))
		sentry.CaptureMessage(bodyResp)
		return "", err
	}

	if response.Response != "1" {
		if response.Response == "2" {
			return "", nil
		} else {
			log.Println("Error getting code")
			msg := fmt.Sprintf("Error getting code SmsPva, api error code: %s", response.Response)
			sentry.CaptureMessage(msg)
			return "", errors.New("code error")
		}
	}

	if len(response.Sms) <= 2 {
		if strings.Contains(response.Text, "Your SNKRDUNK verification code") {
			code := strings.Split(response.Text, "Your SNKRDUNK verification code is: ")[1]
			return code, nil
		}
	}
	if strings.Contains(response.Sms, "Your SNKRDUNK verification code") {
		code := strings.Split(response.Sms, "Your SNKRDUNK verification code is: ")[1]
		return code, nil
	} else {
		return response.Sms, nil
	}
}

//TODO: Your SNKRDUNK verification code is: 463049
