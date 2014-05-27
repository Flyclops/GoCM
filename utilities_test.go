package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/alexjlockwood/gcm"
)

func TestAppendAttempts(t *testing.T) {
	a := runReport.Attempts
	appendAttempts(2)
	if runReport.Attempts != a+2 {
		t.Fatalf("Append attempts should be %d, is %d", a+2, runReport.Attempts)
	}
}

func TestAppendFailures(t *testing.T) {
	a := runReport.Failures
	appendFailures(3)
	if runReport.Failures != a+3 {
		t.Fatalf("Append failures should be %d, is %d", a+3, runReport.Failures)
	}
}

func TestAppendCanonicals(t *testing.T) {
	a := runReport.Canonicals
	appendCanonicals(1)
	if runReport.Canonicals != a+1 {
		t.Fatalf("Append canonicals should be %d, is %d", a+1, runReport.Canonicals)
	}
}

func TestAppendNotRegistered(t *testing.T) {
	a := runReport.NotRegistered
	appendNotRegistered(1)
	if runReport.NotRegistered != a+1 {
		t.Fatalf("Append notregistered should be %d, is %d", a+1, runReport.NotRegistered)
	}
}

func TestIncrementPending(t *testing.T) {
	a := runReport.Pending
	incrementPending()
	if runReport.Pending != a+1 {
		t.Fatalf("Increment pending should be %d, is %d", a+1, runReport.Pending)
	}
}

func TestDecrementPending(t *testing.T) {
	a := runReport.Pending
	decrementPending()
	if runReport.Pending != a-1 {
		t.Fatalf("Decrement pending should be %d, is %d", a-1, runReport.Pending)
	}
}

func TestHandleCanonicalsInResult(t *testing.T) {
	canonicalReplacements = nil
	var results []gcm.Result

	for i := 0; i < 4; i++ {
		g := gcm.Result{"asdf", fmt.Sprintf("%d-%d-%d-%d", i, i, i, i), ""}
		results = append(results, g)
		//gcm.Response{1, 1, 0, 1, gs}
	}

	// "Handle" the results
	handleCanonicalsInResult("asdf", results)

	// Now parse through them
	for i, r := range canonicalReplacements {
		if r.Original != "asdf" {
			t.Fatal("Original is not \"asdf\" as it should be")

			replacement := fmt.Sprintf("%d-%d-%d-%d", i, i, i, i)
			if r.Canonical != replacement {
				t.Fatalf("Canonical is wrong. Expecting: %s, got: %s", replacement, r.Canonical)
			}
		}
	}
}

func TestSendMessageToGCM(t *testing.T) {
	// Test empty token
	ok, err := sendMessageToGCM([]string{}, "")
	if ok {
		t.Fatal("ok should be false")
	}
	if err.Error() != "No tokens were supplied, exiting" {
		t.Fatalf("Expecting 'No tokens were supplied, exiting', got: %s", err.Error())
	}

	// Test empty payload
	ok, err = sendMessageToGCM([]string{"asdf"}, "")
	if ok {
		t.Fatal("ok should be false")
	}
	if err.Error() != "Payload was empty, exiting" {
		t.Fatalf("Expecting 'Payload was empty, exiting', got: %s", err.Error())
	}

	// Test bad json
	ok, err = sendMessageToGCM([]string{"asdf"}, "asdf")
	if ok {
		t.Fatal("ok should be false")
	}
	if !strings.HasPrefix(err.Error(), "invalid character") {
		t.Fatalf("Unexpected error string: %s", err.Error())
	}

	// Test bad send
	aOrig := runReport.Attempts
	fOrig := runReport.Failures

	ok, err = sendMessageToGCM([]string{"asdf"}, "{\"key\": \"value\"}")
	time.Sleep(1 * time.Second)
	if ok {
		t.Fatal("ok should be false")
	}
	if runReport.Attempts != aOrig+1 {
		t.Fatal("Attempts not incremented by 1")
	}
	if runReport.Failures != fOrig+2 {
		// Plus 2 because of the default retry
		t.Fatalf("Failures not incremented by 2 (orig: %d, new: %d)", fOrig, runReport.Failures)
	}
	if !strings.HasPrefix(err.Error(), "401 error") {
		t.Fatalf("Unexpected error string: %s", err.Error())
	}
}
