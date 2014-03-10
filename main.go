package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "github.com/alexjlockwood/gcm"
    "io"
    "log"
    "net/http"
    "runtime"
    "strconv"
)

type procSettings struct {
    IpAddress string
    Port      string
    GCMAPIKey string
}

//=====================
//
// Utility functions
//
//=====================

func sendMessageToGCM(token, jsonStr string) {
    if token == "" {
        log.Println("Token was empty, exiting")
        return
    }

    if jsonStr == "" {
        log.Println("Payload was empty, exiting")
        return
    }

    // Unpack the JSON payload
    var payload map[string]interface{}
    err := json.Unmarshal([]byte(jsonStr), &payload)
    if err != nil {
        log.Println("Can't unmarshal the json: " + err.Error())
        log.Println("Original: " + jsonStr)
        return
    }

    // All is well, make & send the message
    msg := gcm.NewMessage(payload, token)
    sender := &gcm.Sender{ApiKey: settings.GCMAPIKey}
    result, err := sender.Send(msg, 2)
    if err != nil {
        log.Println("Failed to send message:")
        log.Println(err.Error())
    }
    if result != nil {
        log.Printf("Message sent: %s\n", payload["title"])
    }
}

//=====================
//
// Handlers
//
//=====================

func Send(w http.ResponseWriter, r *http.Request) {
    token := r.PostFormValue("token")
    jsonStr := r.PostFormValue("payload")

    // Push the long-running work to a new goroutine
    go sendMessageToGCM(token, jsonStr)

    // Return immediately
    output := "ok\n"
    w.Header().Set("Content-Type", "text/html")
    w.Header().Set("Content-Length", strconv.Itoa(len(output)))
    io.WriteString(w, output)
}

//=====================
//
// Main method
//
//=====================

var settings procSettings

func main() {
    // Set max parallelism
    i := runtime.NumCPU()
    log.Printf("Running on %d CPUs\n", i)
    runtime.GOMAXPROCS(i)

    // Read in flag args
    flag.StringVar(&settings.IpAddress, "ipaddress", "0.0.0.0", "IP address to listen on (default: 0.0.0.0)")
    flag.StringVar(&settings.Port, "port", "5601", "TCP port to listen on (default: 5601)")
    flag.StringVar(&settings.GCMAPIKey, "apikey", "", "GCM API key (required)")
    flag.Parse()

    // Make sure there's an API key
    if settings.GCMAPIKey == "" {
        log.Fatal("No GCM API key provided (--apikey) - GCM API key is required")
    }

    // Start listening
    listenAddress := fmt.Sprintf("%s:%s", settings.IpAddress, settings.Port)
    log.Println("Listening on " + listenAddress)
    http.HandleFunc("/gcm/send", Send)
    log.Fatal(http.ListenAndServe(listenAddress, nil))
}
