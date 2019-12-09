package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ardacetinkaya/FirstGO/azstoragequeue"
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

var configurationSettings config.Config
var logger *log.Logger

type key int

const (
	requestIDKey key = 0
)

func logs() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authToken, err := token.GetRequestToken(req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		if authToken != configurationSettings.Token {
			http.Error(w, "Invalid token", http.StatusBadRequest)

			return
		}
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

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				http.Error(w, "Invalid message body", http.StatusBadRequest)
				return
			}
			bodyString := string(body)
			defer req.Body.Close()

			qm := azstoragequeue.CreateQueueManager()

			if err := qm.Init(configurationSettings.AzureQueueAccountName, configurationSettings.AzureQueueAccountKey); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := qm.CreateQueue("logs"); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := qm.Put(bodyString); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		default:
			http.Error(w, "Invalid method", http.StatusBadRequest)
		}

	})
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

	configurationSettings = config.LoadConfiguration(configFile)

	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("Server is starting...")

	router := http.NewServeMux()
	router.Handle("/logs", logs())

	server := &http.Server{
		Addr:        configurationSettings.Port,
		ErrorLog:    logger,
		IdleTimeout: 25 * time.Second,
		Handler:     (logging(logger)(router)),
	}
	server.SetKeepAlivesEnabled(true)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("PORT is not available %s: %v\n", configurationSettings.Port, err)
	}
}

func logging(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}
