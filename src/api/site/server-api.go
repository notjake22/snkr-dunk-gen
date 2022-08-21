package site

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"log"
	"net/http"
	"os"
	"strings"
)

func StartServer() {
	log.Println("Starting local server for site")

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		redirectHandler(writer, request)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error starting local server")
		sentry.CaptureException(err)
		os.Exit(222)
	}

}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		sentry.CaptureMessage("Method not supported on local server for google login")
		return
	}

	if strings.Contains(r.URL.RawQuery, "state=state-token") {
		handle := strings.Split(r.URL.RawQuery, "code=")[1]
		code := strings.Split(handle, "&scope=")[0]

		message := fmt.Sprintf("Gmail code: %s", code)
		_, err := fmt.Fprintf(w, message)
		if err != nil {
			sentry.CaptureException(err)
			os.Exit(333)
		}
	} else {
		_, err := fmt.Fprintf(w, "ERROR GETTING CODE 404")
		if err != nil {
			sentry.CaptureException(err)
			os.Exit(333)
		}
	}
}
