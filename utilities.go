package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/alexjlockwood/gcm"
)

func sendMessageToGCM(tokens []string, payloadAsString string) (bool, error) {
	// At any exit, decrement pending
	defer func() {
		go decrementPending()
	}()

	if len(tokens) == 0 {
		errText := "No tokens were supplied, exiting"
		log.Println(errText)
		return false, errors.New(errText)
	}

	if payloadAsString == "" {
		errText := "Payload was empty, exiting"
		log.Println(errText)
		return false, errors.New(errText)
	}

	// Unpack the JSON payload
	var payload map[string]interface{}
	err := json.Unmarshal([]byte(payloadAsString), &payload)
	if err != nil {
		log.Println("Can't unmarshal the json: " + err.Error())
		log.Println("Original: " + payloadAsString)
		return false, err
	}

	// All is well, make & send the message
	go appendAttempts(len(tokens))

	msg := gcm.NewMessage(payload, tokens...)
	sender := &gcm.Sender{ApiKey: settings.GCMAPIKey}
	response, err := sender.Send(msg, 2)
	if err != nil {
		log.Println("Failed to send message:")
		log.Println(err.Error())

		go appendFailures(1)
		return false, err
	}

	numCan := 0
	numErr := 0
	if response != nil {
		for i, result := range response.Results {
			// Canonicals
			if result.RegistrationID != "" {
				numCan++
				canonicalReplacements = append(canonicalReplacements, canonicalReplacement{tokens[i], result.RegistrationID})
			}
			if result.Error != "" {
				numErr++
				log.Printf("Error sending: %s", result.Error)

				if result.Error == "NotRegistered" {
					handleNotRegisteredError(tokens[i])
					go appendNotRegistered(1)
				}
			}
		}

		go appendCanonicals(numCan)
		go appendFailures(numErr)
	}

	log.Printf("Message sent. Attempts: %d, Errors: %d, Successful: %d (Canonicals: %d)", len(tokens), numErr, len(tokens)-numErr, numCan)

	return true, nil
}

func handleCanonicalsInResult(original string, results []gcm.Result) {
	for _, r := range results {
		canonicalReplacements = append(canonicalReplacements, canonicalReplacement{original, r.RegistrationID})
	}
}

func handleNotRegisteredError(original string) {
	notRegisteredMutex.Lock()
	notRegisteredKeys = append(notRegisteredKeys, original)
	notRegisteredMutex.Unlock()
}

func appendAttempts(numToAppend int) {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Attempts += numToAppend
}

func appendFailures(numToAppend int) {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Failures += numToAppend
}

func appendCanonicals(numToAppend int) {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Canonicals += numToAppend
}

func appendNotRegistered(numToAppend int) {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.NotRegistered += numToAppend
}

func incrementPending() {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Pending++
}

func decrementPending() {
	runReportMutex.Lock()
	defer runReportMutex.Unlock()
	runReport.Pending--
}
