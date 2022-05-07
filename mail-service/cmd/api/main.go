package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const webPort = "80"

func main() {
	logSnippet := "\n[mail-service][main][main] =>"

	mailer, err := createMail()
	if err != nil {
		log.Printf("%s (ERROR||createMail): %s", logSnippet, err.Error())
		log.Panic(err)
	}

	app := Config{
		Mailer: mailer,
	}

	log.Printf("%s (Starting mail service on port %s)", logSnippet, webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("%s (ERROR||srv.ListenAndServe): %s", logSnippet, err.Error())
		log.Panic(err)
	}
}

func createMail() (Mail, error) {
	logSnippet := "\n[mail-service][main][createMail] =>"

	port, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil {
		log.Printf("%s (ERROR-strconv.Atoi(MAIL_PORT): %s", logSnippet, err.Error())
		return Mail{}, err
	}

	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
	}

	return m, nil
}
