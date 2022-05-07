package main

import (
	"log"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	logSnippet := "\n[mail-service][handlers][SendMail] =>"

	var requestPayload mailMessage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Printf("%s (ERROR-app.readJSON): %s", logSnippet, err.Error())
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Printf("%s (ERROR-app.Mailer.SendSMTPMessage): %s", logSnippet, err.Error())
		app.errorJSON(w, err)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "sent to " + requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, responsePayload)
}
