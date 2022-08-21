package cli

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/getsentry/sentry-go"
	"log"
	"main/src/api/gmail"
	"main/src/settings"
	"os"
	"time"
)

func StartCli() {
	answers := struct {
		Selection string
	}{}

	var qs = []*survey.Question{
		{
			Name: "Selection",
			Prompt: &survey.Select{
				Message: "What would you like to do?",
				Options: []string{"Start Bot", "Configure gmail auth", "Edit Settings", "Exit"},
			},
		},
	}

	err := survey.Ask(qs, &answers)
	if err != nil {
		log.Println("Error prompting questions")
		sentry.CaptureException(err)
		time.Sleep(3000)
		os.Exit(0)
	}

	if answers.Selection == "Start Bot" {
		taskEngine()
		StartCli()
	}
	if answers.Selection == "Configure gmail auth" {
		gmail.GetUserToken()
		StartCli()
	}
	if answers.Selection == "Edit Settings" {
		answers = struct {
			Selection string
		}{}
		var qs = []*survey.Question{
			{
				Name: "Selection",
				Prompt: &survey.Select{
					Message: "What would you like to edit?",
					Options: []string{"Catchall", "Referral Code", "SMS Api key", "SMS Price max", "Discord webhook", "return to main page"},
				},
			},
		}

		err = survey.Ask(qs, &answers)
		if err != nil {
			log.Println("Error prompting question, exiting")
			sentry.CaptureException(err)
			time.Sleep(3000)
			os.Exit(1)
		}

		if answers.Selection == "Catchall" {
			settings.EditCatchall()
			ClearConsole()
			log.Println("Settings saved!")
			StartCli()
		}
		if answers.Selection == "Referral Code" {
			settings.EditReferral()
			ClearConsole()
			log.Println("Settings saved!")
			StartCli()
		}
		if answers.Selection == "SMS Api key" {
			settings.EditSMSApi()
			ClearConsole()
			log.Println("Settings saved!")
			StartCli()
		}
		if answers.Selection == "SMS Price max" {
			settings.EditSMSPriceMax()
			ClearConsole()
			log.Println("Settings saved!")
			StartCli()
		}
		if answers.Selection == "Discord webhook" {
			settings.EditWebhook()
			ClearConsole()
			log.Println("Settings saved!")
			StartCli()
		}
		if answers.Selection == "return to main page" {
			ClearConsole()
			StartCli()
		}
	}
}
