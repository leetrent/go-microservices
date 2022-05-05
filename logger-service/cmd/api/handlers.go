package main

import (
	"log"
	"log-service/data"
	"net/http"
)

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	logSnippet := "\n[logger-service][handlers][WriteLog] =>"

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
		log.Printf("%s (ERROR-app.Models.LogEntry.Insert): %s", logSnippet, err.Error())
		app.errorJSON(w, err)
		return
	}
	log.Printf("%s (SUCCESS-app.Models.LogEntry.Insert)", logSnippet)

	resp := jsonResponse{
		Error:   false,
		Message: "entry logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
