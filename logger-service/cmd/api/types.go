package main

import "log-service/data"

type Config struct {
	Models data.Models
}

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}
