package cli

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"main/src/module"
	"main/src/settings"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func taskEngine() {
	log.Println("Starting bot")

	settingsCheck := settings.ReadSettings()
	if len(settingsCheck.SmsAPIKey) < 3 {
		log.Println("Please set an sms api key in settings, free mode disabled")
		//StartCli()
	}
	if len(settingsCheck.Catchall) < 3 {
		log.Println("Please set a catchall in settings")
		StartCli()
	}
	if len(settingsCheck.ReferralCode) < 3 {
		log.Println("Please set your referral code in settings")
		StartCli()
	}

	fmt.Printf("Enter how many tasks you'd like to run: ")
	var taskAmount string
	if _, err := fmt.Scan(&taskAmount); err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Unable to read task number code: %v", err)
	}
	taskNum, err := strconv.Atoi(strings.TrimSpace(taskAmount))
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Unable to read task number code: %v", err)
	}
	if taskNum > 25 {
		log.Println("Task limit 25, returning to cli page")
		StartCli()
	}

	var taskIds []string
	for i := 0; i < taskNum; i++ {
		slice := []string{strconv.Itoa(i + 1)}
		taskIds = append(taskIds, slice...)
	}
	taskOptions := module.TaskOptions{Settings: settings.ReadSettings()}
	go func() {
		for _, taskId := range taskIds {
			taskOptions.TaskID = taskId
			module.Task(taskOptions)
		}
		fmt.Println("\nTo return to main cli page, press control and c at the same time :)")
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
