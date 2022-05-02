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
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
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
