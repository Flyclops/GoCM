package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"sync"
)

type procSettings struct {
	Host      string
	Port      string
	GCMAPIKey string
	Logto     string
}

type canonicalReplacement struct {
	Original  string `json:"original"`
	Canonical string `json:"canonical"`
}

type report struct {
	Attempts   int `json:"attempts"`
	Failures   int `json:"failures"`
	Canonicals int `json:"canonicals"`
}

var settings procSettings

// Reporting, with matching sync mutexes
var runReport report
var runReportMutex sync.Mutex

var canonicalReplacements []canonicalReplacement
var canonicalReplacementsMutex sync.Mutex

//=====================
//
// Main method
//
//=====================

func main() {
	// Set max parallelism
	i := runtime.NumCPU()
	log.Printf("Running on %d CPUs\n", i)
	runtime.GOMAXPROCS(i)

	// Read in flag args
	flag.StringVar(&settings.Host, "host", "0.0.0.0", "IP address to listen on (default: 0.0.0.0)")
	flag.StringVar(&settings.Port, "port", "5601", "TCP port to listen on (default: 5601)")
	flag.StringVar(&settings.GCMAPIKey, "apikey", "", "GCM API key (required)")
	flag.StringVar(&settings.Logto, "logto", "", "Path to log file (default: stdout)")
	flag.Parse()

	// Set up logging
	if settings.Logto != "" {
		fmt.Println("Log path: ", settings.Logto)
		f, err := os.OpenFile(settings.Logto, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("Log file won't open", err)
		} else {
			log.SetOutput(f)
			defer f.Close()
		}
	}

	// Make sure there's an API key
	if settings.GCMAPIKey == "" {
		log.Fatal("No GCM API key provided (--apikey) - GCM API key is required")
	}

	// API Key in log
	log.Printf("Using GCM API key %s\n", settings.GCMAPIKey)

	// Start listening
	listenAddress := fmt.Sprintf("%s:%s", settings.Host, settings.Port)
	log.Println("Listening on " + listenAddress)

	// Set up handlers
	http.HandleFunc("/gcm/send", send)
	http.HandleFunc("/gcm/report", getReport)
	http.HandleFunc("/gcm/report/canonical", getCanonicalReport)

	// Start the listener
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
