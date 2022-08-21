package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"main/src/cli"
	"main/src/security"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main() {
	log.Println("Starting services")

	go func() {
		for {
			security.CheckDebuggers()
			time.Sleep(500 * time.Millisecond)
		}
	}()
	log.Println("Anti-crack service started")
	if !security.CheckApi() {
		log.Println("Error meeting security param 1 at launch")
		time.Sleep(5000 * time.Millisecond)
		os.Exit(3)
	}
	log.Println("Security authorization services started")

	log.Println("Starting error logging service")

	if !startSentry() {
		log.Println("Error starting error logging api connection")
		log.Println("No fatal error. Continuing")
	}

	log.Println("Starting Cli..")
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/C", "title", fmt.Sprintf("Snkr Dunk Gen"))
		cmd.Run()
	default:

	}
	time.Sleep(1000 * time.Millisecond)
	cli.ClearConsole()
	cli.StartCli()
}

func startSentry() bool {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "", // removed
		Debug:            false,
		TracesSampleRate: 1.0,
		//TracesSampler: sentry.TracesSamplerFunc(func(ctx sentry.SamplingContext) sentry.Sampled {
		//	return sentry.SampledTrue
		//}),
	})
	if err != nil {
		log.Println(err)
		return false
	}
	defer sentry.Flush(2 * time.Second)

	return true
}
