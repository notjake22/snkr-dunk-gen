package module

import (
	"github.com/getsentry/sentry-go"
	"golang.org/x/net/publicsuffix"
	"log"
	"main/src/api/gmail"
	"main/src/settings"
	"main/src/webhooks"
	"net/http"
	"net/http/cookiejar"
	"os"
	"time"
)

type TaskOptions struct {
	settings.Settings
	TaskID string
}

type SnkrDunkTask struct {
	Options TaskOptions

	csrfLoginToken  string
	csrfHeaderToken string
	EnsID           string
	client          http.Client

	email       string
	password    string
	phoneNumber string
}

func Task(taskOptions TaskOptions) {
	// was originally going to use state constants, but decided to just do it like this instead
	t := SnkrDunkTask{
		Options: taskOptions,
	}

	log.Println("Initializing task: ", t.Options.TaskID)
	err := t.initialize()
	if err != nil {
		log.Println("Error initializing task")
		sentry.CaptureMessage("Error initializing task")
		time.Sleep(3000 * time.Millisecond)
		os.Exit(1)
	}

	log.Println("Getting session")
	err = t.getSession()
	if err != nil {
		log.Println("Error getting session")
		sentry.CaptureMessage("Error getting session")
		time.Sleep(3000 * time.Millisecond)
		os.Exit(1)
	}
	log.Println("Got session")

	log.Println("Submitting credentials for new account..")
	err = t.submitNewAccount()
	if err != nil {
		log.Println("Error submitting new account")
		sentry.CaptureMessage("Error submitting new account")
		time.Sleep(3000 * time.Millisecond)
		os.Exit(1)
	}
	log.Println("Account created")

	log.Println("Verifying email..")
	log.Println("Waiting for email from snkr dunk")

	url := ""
	for url == "" {
		url, err = gmail.ReadUserMessages()
		if err != nil {
			log.Println("Error reading email")
			sentry.CaptureMessage("Error reading user email for url")
			time.Sleep(3000 * time.Millisecond)
			os.Exit(1)
		}
		time.Sleep(1000 * time.Millisecond)
	}

	log.Println("Got verification url")
	err = t.verifyEmailReq(url)
	if err != nil {
		log.Println("Error verifying email")
		sentry.CaptureMessage("Error verifying email")
		time.Sleep(3000 * time.Millisecond)
		os.Exit(1)
	}
	log.Println("Verified email!")
	webhooks.EmailVerified(t.email, t.password, t.Options.TaskID)

	log.Println("Verifying sms")
	err = t.verifySms()
	if err != nil {
		log.Println("Error verifying sms:", err)
		sentry.CaptureMessage("Error verifying sms")
		time.Sleep(3000 * time.Millisecond)
		os.Exit(1)
	}
	webhooks.SmsVerified(t.email, t.password, t.Options.TaskID)
	log.Println("Verified sms!")

	log.Println("Adding address")
	err = t.addAddress()
	if err != nil {
		log.Println("Error adding address:", err)
		sentry.CaptureMessage("Error setting account address")
		time.Sleep(3000 * time.Millisecond)
		os.Exit(1)
	}
	log.Println("Address added!")

	log.Println("Submitting code")
	err = t.submitReferralCode()
	if err != nil {
		log.Println("Error submitting code:", err)
		sentry.CaptureMessage("Error submitting referral code")
		time.Sleep(3000 * time.Millisecond)
		os.Exit(1)
	}
	log.Println("Code submitted!")

	return
}

func (t *SnkrDunkTask) initialize() error {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	t.client.Jar = jar

	err = t.proxyHandler()
	if err != nil {
		sentry.CaptureException(err)
		log.Println("Error getting proxy on initialing")
		return err
	}

	return nil
}
