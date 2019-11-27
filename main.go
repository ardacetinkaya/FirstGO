package main

import (
	"encoding/json"
	"net/http"
	"time"

	"./config"
)

//Log mode
type Log struct {
	Text    string    `json:"text"`
	LogType string    `json:"type"`
	Date    time.Time `json:"date"`
}

//Logs variable
type Logs []Log

func logs(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		logs := Logs{
			Log{Text: "System is running", Date: time.Now()},
			Log{Text: "Error", Date: time.Now()},
		}

		if err := json.NewEncoder(w).Encode(logs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "POST":
		var l Log
		err := json.NewDecoder(req.Body).Decode(&l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//TODO: Save
	default:
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}

}

func main() {

	configFile := config.LoadConfiguration("config.json")

	http.HandleFunc("/logs", logs)

	if err := http.ListenAndServe(configFile.Port, nil); err != nil {
		panic(err)
	}
}
