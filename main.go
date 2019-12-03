package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/ardacetinkaya/FirstGO/config"
	"github.com/ardacetinkaya/FirstGO/token"
)

//Log mode
type Log struct {
	Text    string    `json:"text"`
	LogType string    `json:"type"`
	Date    time.Time `json:"date"`
}

//Logs variable
type Logs []Log

var ConfigurationSettings config.Config

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

		authToken, err := token.GetRequestToken(req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if authToken != ConfigurationSettings.Token {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Invalid message body", http.StatusBadRequest)
			return
		}
		bodyString := string(body)
		defer req.Body.Close()

		credential, err := azqueue.NewSharedKeyCredential(ConfigurationSettings.AzureQueueAccountName, ConfigurationSettings.AzureQueueAccountKey)
		if err != nil {
			http.Error(w, "Invalid credential settings", http.StatusBadRequest)
			return
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
		_, err = messagesURL.Enqueue(ctx, bodyString, time.Second*0, time.Minute)
		if err != nil {
			panic(err)
		}

	default:
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}

}

func main() {
	var configFile string
	if len(os.Getenv("APP_ENV")) <= 0 {
		panic(string("Unknown Environment"))
	} else {
		if os.Getenv("APP_ENV") == "Development" {
			configFile = fmt.Sprintf("config.%s.json", "development")
		} else {
			configFile = "config.json"
		}
	}

	ConfigurationSettings = config.LoadConfiguration(configFile)

	http.HandleFunc("/logs", logs)

	if err := http.ListenAndServe(ConfigurationSettings.Port, nil); err != nil {
		panic(err)
	}
}
