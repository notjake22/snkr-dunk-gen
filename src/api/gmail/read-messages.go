package gmail

import (
	"context"
	"encoding/base64"
	"github.com/getsentry/sentry-go"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"strings"
)

// most of this was copy and pasted from googles docs

func ReadUserMessages() (string, error) {
	ctx := context.Background()
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	//log.Println("Checking login token")

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Unable to parse client secret file to config: %v\n", err)
		return "", err
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Unable to retrieve Gmail client: %v\n", err)
		return "", err
	}

	user := "me"
	m, err := srv.Users.Messages.List(user).Do()
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Unable to retrieve messages: %v\n", err)
		return "", err
	}
	if len(m.Messages) == 0 {
		log.Println("No messages found")
		return "", err
	}

	checks := 0
	for _, l := range m.Messages {
		//fmt.Printf("- %+v\n", l.Id)
		if checks == 1 {
			return "", nil
		}

		msg, err := srv.Users.Messages.Get(user, l.Id).Do()
		if err != nil {
			sentry.CaptureException(err)
			log.Println("Error getting messages")
		}

		//log.Println(msg.Payload.Parts[0].Body.Data)
		if len(msg.Payload.Parts) == 0 {
			//log.Println("No body in email, continuing")
			checks = checks + 1
			continue
		}
		str := msg.Payload.Parts[0].Body.Data
		str = strings.TrimSpace(str)

		dst := make([]byte, base64.StdEncoding.DecodedLen(len(str)))
		n, err := base64.StdEncoding.Decode(dst, []byte(str))
		if err != nil {
			//fmt.Println("decode error:", err)
			//log.Println("Err decoding email, continuing")
			checks = checks + 1
			continue
		}
		dst = dst[:n]
		//fmt.Printf("%q\n", string(dst))

		check := strings.Contains(string(dst), "Please access the following URL")
		if !check {
			//log.Println("Email does not match")
			checks = checks + 1
			continue
		}

		msgSplit := strings.Split(string(dst), "\n")[4]
		url := strings.Split(msgSplit, "\r")[0]

		//err = srv.Users.Messages.Delete(user, l.Id).Do()
		//if err != nil {
		//	log.Println("Error deleting email: ", err)
		//	return "", err
		//}

		return url, nil
	}

	return "", nil
}
