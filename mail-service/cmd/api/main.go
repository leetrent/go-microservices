package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

func main() {
	logSnippet := "\n[mail-service][main][main] =>"
	app := Config{}

	log.Printf("%s (Starting mail service on port %s)", logSnippet, webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Printf("%s (ERROR||srv.ListenAndServe): %s", logSnippet, err.Error())
		log.Panic(err)
	}
}
