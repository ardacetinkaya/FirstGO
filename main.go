package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Log struct {
	Text    string    `json:"text"`
	LogType string    `json:"type"`
	Date    time.Time `json:"date"`
}

type Logs []Log

func logs(w http.ResponseWriter, req *http.Request) {
	logs := Logs{
		Log{Text: "System is running", LogType: "Info", Date: time.Now()},
		Log{Text: "Error", LogType: "Error", Date: time.Now()},
	}

	if err := json.NewEncoder(w).Encode(logs); err != nil {
		panic(err)
	}
}

func main() {

	http.HandleFunc("/logs", logs)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		panic(err)
	}
}
