package module

import (
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"main/src/webhooks"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (t *SnkrDunkTask) submitNewAccount() error {
	t.genEmail()
	t.genPassword()
	username := genUser()

	form := url.Values{}
	form.Add("username", username)
	form.Add("email", t.email)
	form.Add("password", t.password)
	form.Add("agreement", "on")
	form.Add("csrf_token", t.csrfLoginToken)
	form.Add("tzDatabaseName", "America/New_York")

	req, err := http.NewRequest("POST", "https://snkrdunk.com/en/signup", strings.NewReader(form.Encode()))
	if err != nil {
		sentry.CaptureException(err)
		log.Println("Error making request")
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
		"referer":                   []string{"https://snkrdunk.com/en/signup"},
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
		log.Println("Error sending request")
		return err
	}

	if res.StatusCode >= 400 {
		log.Println("Request error, status code:", res.StatusCode)
		if res.StatusCode == 409 {
			sentry.CaptureMessage("Email already exists")
			log.Println("Email already exists")
		} else {
			msg := fmt.Sprintf("Request error, status code: %v", res.StatusCode)
			sentry.CaptureMessage(msg)
		}
		return errors.New("request error")
	}

	webhooks.NewAccountWebhook(username, t.email, t.password, t.Options.TaskID)
	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("1234567890")
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func genUser() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (t *SnkrDunkTask) genEmail() {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	email := fmt.Sprintf("%s@%s", string(b), t.Options.Catchall)
	t.email = email

	return
}

func (t *SnkrDunkTask) genPassword() {
	b := make([]rune, 7)
	for i := range b {
		b[i] = letterRunes[seededRand.Intn(len(letterRunes))]
	}
	a := make([]rune, 7)
	for i := range a {
		a[i] = numberRunes[seededRand.Intn(len(numberRunes))]
	}
	t.password = fmt.Sprintf("%s%s", string(b), string(a))
	return
}
