package webhooks

import (
	"bytes"
	"fmt"
	"log"
	"main/src/settings"
	"net/http"
)

func NewAccountWebhook(user string, email string, password string, taskId string) {
	webhookUrl := settings.ReadSettings().Webhook

	if len(webhookUrl) < 2 {
		return
	}

	payload := fmt.Sprintf(`{
  "content": null,
  "embeds": [
    {
      "title": "Snkr Dunk Account generated",
		"description": "task id: %s",
      "color": 8567148,
      "fields": [
        {
          "name": "Username",
          "value": "%s",
          "inline": true
        },
        {
          "name": "Email",
          "value": "%s",
          "inline": true
        },
        {
          "name": "Password",
          "value": "||%s||"
        }
      ]
    }
  ],
  "attachments": []
}`, taskId, user, email, password)

	req, err := http.NewRequest("POST", webhookUrl, bytes.NewReader([]byte(payload)))
	if err != nil {
		log.Println("Error sending webhook")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Error sending webhook")
		return
	}
	return
}
