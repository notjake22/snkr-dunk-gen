package settings

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func EditCatchall() {
	settings := ReadSettings()

	var newCatchall string
	fmt.Printf("Enter your new catchall here: ")
	_, err := fmt.Scanln(&newCatchall)
	if err != nil {
		log.Println("Error reading catchall")
		sentry.CaptureException(err)
		os.Exit(0)
	}

	newSettings, err := json.Marshal(Settings{
		Catchall:     strings.TrimSpace(newCatchall),
		ReferralCode: settings.ReferralCode,
		SmsAPIKey:    settings.SmsAPIKey,
		SmsPriceMax:  settings.SmsPriceMax,
		Webhook:      settings.Webhook,
	})
	_, err = os.Create("./settings.json")
	if err != nil {
		log.Println(err)
		sentry.CaptureException(err)
		os.Exit(3)
	}
	err = ioutil.WriteFile("./settings.json", newSettings, os.ModePerm)
	if err != nil {
		log.Println("Error writing new settings to file")
		sentry.CaptureException(err)
		os.Exit(4)
	}
	log.Println("Settings changed!")
	return
}

func EditSMSApi() {
	settings := ReadSettings()

	var newSMSApiKey string
	fmt.Printf("Enter your new SMS API key here: ")
	_, err := fmt.Scanln(&newSMSApiKey)
	if err != nil {
		log.Println("Error reading api key")
		sentry.CaptureException(err)
		os.Exit(0)
	}

	newSettings, err := json.Marshal(Settings{
		Catchall:     settings.Catchall,
		ReferralCode: settings.ReferralCode,
		SmsAPIKey:    strings.TrimSpace(newSMSApiKey),
		SmsPriceMax:  settings.SmsPriceMax,
		Webhook:      settings.Webhook,
	})
	_, err = os.Create("./settings.json")
	if err != nil {
		log.Println(err)
		sentry.CaptureException(err)
		os.Exit(3)
	}
	err = ioutil.WriteFile("./settings.json", newSettings, os.ModePerm)
	if err != nil {
		log.Println("Error writing new settings to file")
		sentry.CaptureException(err)
		os.Exit(4)
	}
	log.Println("Settings changed!")
	return
}

func EditSMSPriceMax() {
	settings := ReadSettings()

	var newPrice string
	fmt.Printf("Enter the max price per number you like to limit: ")
	_, err := fmt.Scanln(&newPrice)
	if err != nil {
		log.Println("Error reading api key")
		sentry.CaptureException(err)
		os.Exit(0)
	}

	newSettings, err := json.Marshal(Settings{
		Catchall:     settings.Catchall,
		ReferralCode: settings.ReferralCode,
		SmsAPIKey:    settings.SmsAPIKey,
		SmsPriceMax:  strings.TrimSpace(newPrice),
		Webhook:      settings.Webhook,
	})
	_, err = os.Create("./settings.json")
	if err != nil {
		log.Println(err)
		sentry.CaptureException(err)
		os.Exit(3)
	}
	err = ioutil.WriteFile("./settings.json", newSettings, os.ModePerm)
	if err != nil {
		log.Println("Error writing new settings to file")
		sentry.CaptureException(err)
		os.Exit(4)
	}
	log.Println("Settings changed!")
	return
}

func EditWebhook() {
	settings := ReadSettings()

	var newWebhook string
	fmt.Printf("Enter your new Webhook here: ")
	_, err := fmt.Scanln(&newWebhook)
	if err != nil {
		log.Println("Error reading catchall")
		sentry.CaptureException(err)
		os.Exit(0)
	}

	newSettings, err := json.Marshal(Settings{
		Catchall:     settings.Catchall,
		ReferralCode: settings.ReferralCode,
		SmsAPIKey:    settings.SmsAPIKey,
		SmsPriceMax:  settings.SmsPriceMax,
		Webhook:      strings.TrimSpace(newWebhook),
	})
	_, err = os.Create("./settings.json")
	if err != nil {
		log.Println(err)
		sentry.CaptureException(err)
		os.Exit(3)
	}
	err = ioutil.WriteFile("./settings.json", newSettings, os.ModePerm)
	if err != nil {
		log.Println("Error writing new settings to file")
		sentry.CaptureException(err)
		os.Exit(4)
	}
	log.Println("Settings changed!")
	return
}

func EditReferral() {
	settings := ReadSettings()

	var newCode string
	fmt.Printf("Enter your new Referral code here: ")
	_, err := fmt.Scanln(&newCode)
	if err != nil {
		log.Println("Error reading code")
		sentry.CaptureException(err)
		os.Exit(0)
	}

	newSettings, err := json.Marshal(Settings{
		Catchall:     settings.Catchall,
		ReferralCode: strings.TrimSpace(newCode),
		SmsAPIKey:    settings.SmsAPIKey,
		SmsPriceMax:  settings.SmsPriceMax,
		Webhook:      settings.Webhook,
	})
	_, err = os.Create("./settings.json")
	if err != nil {
		log.Println(err)
		sentry.CaptureException(err)
		os.Exit(3)
	}
	err = ioutil.WriteFile("./settings.json", newSettings, os.ModePerm)
	if err != nil {
		log.Println("Error writing new settings to file")
		sentry.CaptureException(err)
		os.Exit(4)
	}
	log.Println("Settings changed!")
	return
}

func ReadSettings() Settings {
	jsonFile, err := os.ReadFile("./settings.json")
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
		os.Exit(10)
	}

	var settings Settings
	err = json.Unmarshal(jsonFile, &settings)
	if err != nil {
		fmt.Println(err)
		sentry.CaptureException(err)
		os.Exit(10)
	}

	return settings
}

type Settings struct {
	Catchall     string `json:"catchall"`
	ReferralCode string `json:"referral-code"`
	SmsAPIKey    string `json:"sms-api-key"`
	SmsPriceMax  string `json:"sms-price-max"`
	Webhook      string `json:"webhook"`
}
