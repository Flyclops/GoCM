package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/alexjlockwood/gcm"
)

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
	go appendAttempts()

	msg := gcm.NewMessage(payload, token)
	sender := &gcm.Sender{ApiKey: settings.GCMAPIKey}
	result, err := sender.Send(msg, 2)
	if err != nil {
		log.Println("Failed to send message:")
		log.Println(err.Error())

		go appendFailures()
		return false, err
	}
	canonicalsBack := 0
	if result != nil {
		canonicalsBack = result.CanonicalIDs
		//log.Printf("Message sent: %s\n", payload["title"])
		if result.CanonicalIDs > 0 {
			go appendCanonicals()
			handleCanonicalsInResult(token, result.Results)
		}
	}

	log.Printf("Message sent, canonicals: %d", canonicalsBack)

	return true, nil
}

func handleCanonicalsInResult(original string, results []gcm.Result) {
	for _, r := range results {
		canonicalReplacements = append(canonicalReplacements, canonicalReplacement{original, r.RegistrationID})
	}
}

func appendAttempts() {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Attempts++
}

func appendFailures() {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Failures++
}

func appendCanonicals() {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Canonicals++
}
