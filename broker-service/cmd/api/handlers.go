package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker (using helper functions)",
	}

	err := app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		fmt.Println(err)
	}
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	logSnippet := "[broker-service][handlers][HandleSubmission] =>"

	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("\n%s (requestPayload.Action): %s", logSnippet, requestPayload.Action)

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown submission request"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, ap AuthPayload) {
	/////////////////////////////////////////////////////////
	// Create request JSON to send to authentication-service
	/////////////////////////////////////////////////////////
	jsonData, err := json.MarshalIndent(ap, "", "\t")
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	/////////////////////////////////////////////////////////
	// Call authentication-service
	/////////////////////////////////////////////////////////
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	/////////////////////////////////////////////////////////
	// Make sure we get back the correct status code
	/////////////////////////////////////////////////////////
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling authentication service"))
		return
	}

	/////////////////////////////////////////////////////////////
	// Decode the JSON sent back from the authentication service
	/////////////////////////////////////////////////////////////
	var jsonFromAuthService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromAuthService)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	if jsonFromAuthService.Error {
		app.errorJSON(w, errors.New(jsonFromAuthService.Message), http.StatusUnauthorized)
	}

	var responsePayload jsonResponse
	responsePayload.Error = false
	responsePayload.Message = "Authenticated"
	responsePayload.Data = jsonFromAuthService.Data

	app.writeJSON(w, http.StatusAccepted, responsePayload)
}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {

	logSnippet := "\n[broker-service][handlers][logItem] =>"

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	log.Printf("\n%s (response.StatusCode): %d", logSnippet, response.StatusCode)

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("logging attempt failed"), response.StatusCode)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "entry was successfully logged",
	}

	app.writeJSON(w, http.StatusAccepted, responsePayload)
}

func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	logSnippet := "\n[broker-service][handlers][sendMail] =>"

	jsonData, err := json.MarshalIndent(msg, "", "\t")
	if err != nil {
		log.Printf("%s (ERROR-json.MarshalIndent): %s", logSnippet, err.Error())
		app.errorJSON(w, err)
		return
	}

	mailServiceURL := "http://mailer-service/send"

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("%s (ERROR-http.NewRequest): %s", logSnippet, err.Error())
		app.errorJSON(w, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("%s (ERROR-client.Do): %s", logSnippet, err.Error())
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	log.Printf("\n%s (response.StatusCode): %d", logSnippet, response.StatusCode)

	if response.StatusCode != http.StatusAccepted {
		log.Printf("%s (ERROR-response.StatusCode): %d", logSnippet, response.StatusCode)
		app.errorJSON(w, errors.New("error calling mail service"), response.StatusCode)
		return
	}

	responsePayload := jsonResponse{
		Error:   false,
		Message: "Email message successfully sent to " + msg.To,
	}

	app.writeJSON(w, http.StatusAccepted, responsePayload)
}
