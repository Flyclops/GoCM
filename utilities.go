package main

import (
    "encoding/json"
    "errors"
    "github.com/alexjlockwood/gcm"
    "log"
)

//=====================
//
// Utility functions
//
//=====================

func sendMessageToGCM(token, jsonStr string) (bool, error) {
    if token == "" {
        errText := "Token was empty, exiting"
        log.Println(errText)
        return false, errors.New(errText)
    }

    if jsonStr == "" {
        errText := "Payload was empty, exiting"
        log.Println(errText)
        return false, errors.New(errText)
    }

    // Unpack the JSON payload
    var payload map[string]interface{}
    err := json.Unmarshal([]byte(jsonStr), &payload)
    if err != nil {
        log.Println("Can't unmarshal the json: " + err.Error())
        log.Println("Original: " + jsonStr)
        return false, err
    }

    // All is well, make & send the message
    msg := gcm.NewMessage(payload, token)
    sender := &gcm.Sender{ApiKey: settings.GCMAPIKey}
    go appendAttempts()
    result, err := sender.Send(msg, 2)
    if err != nil {
        log.Println("Failed to send message:")
        log.Println(err.Error())

        go appendFailures()

        return false, err
    }
    if result != nil {
        //log.Printf("Message sent: %s\n", payload["title"])
        if result.CanonicalIDs > 0 {
            go appendCanonicals()
            handleCanonicalsInResult(token, result.Results)
        }
    }

    return true, nil
}

func handleCanonicalsInResult(original string, results []gcm.Result) {
    for _, r := range results {
        canonicalReplacements = append(canonicalReplacements, canonicalReplacement{original, r.RegistrationID})
    }
}

func appendAttempts() {
    runReport.Attempts++
}

func appendFailures() {
    runReport.Failures++
}

func appendCanonicals() {
    runReport.Canonicals++
}
