package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/ardacetinkaya/FirstGO/config"
)

//Log mode
type Log struct {
	Text    string    `json:"text"`
	LogType string    `json:"type"`
	Date    time.Time `json:"date"`
}

//Logs variable
type Logs []Log

//Configurations variable for settings
type ConfigurationSettings config.Config

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
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		bodyString := string(bodyBytes)

		ConfigurationSettings := config.LoadConfiguration("config.json")
		credential, err := azqueue.NewSharedKeyCredential(ConfigurationSettings.AzureQueueAccountName, ConfigurationSettings.AzureQueueAccountKey)
		if err != nil {
			panic(err)
		}
		p := azqueue.NewPipeline(credential, azqueue.PipelineOptions{})
		u, _ := url.Parse(fmt.Sprintf("https://%s.queue.core.windows.net", ConfigurationSettings.AzureQueueAccountName))
		serviceURL := azqueue.NewServiceURL(*u, p)

		queueURL := serviceURL.NewQueueURL("logs")
		ctx := context.TODO()
		_, err = queueURL.Create(ctx, azqueue.Metadata{})
		if err != nil {
			panic(err)
		}

		messagesURL := queueURL.NewMessagesURL()
		println(bodyString)
		_, err = messagesURL.Enqueue(ctx, bodyString, time.Second*0, time.Minute)
		if err != nil {
			panic(err)
		}

	default:
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}

}

func main() {

	ConfigurationSettings := config.LoadConfiguration("config.json")

	http.HandleFunc("/logs", logs)

	if err := http.ListenAndServe(ConfigurationSettings.Port, nil); err != nil {
		panic(err)
	}
}
