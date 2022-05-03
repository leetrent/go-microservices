package main

import (
	"log-service/data"
	"net/http"
)

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// Read JSON in requetPayload
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	// Insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "entry logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
