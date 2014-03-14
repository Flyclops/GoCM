package main

import (
    "fmt"
    "github.com/alexjlockwood/gcm"
    "log"
    "strings"
    "testing"
    "time"
)

func TestAppendAttempts(t *testing.T) {
    a := runReport.Attempts
    appendAttempts()
    if runReport.Attempts != a+1 {
        log.Fatalf("Append attempts should be %d, is %d", a+1, runReport.Attempts)
    }
}

func TestAppendFailures(t *testing.T) {
    a := runReport.Failures
    appendFailures()
    if runReport.Failures != a+1 {
        log.Fatalf("Append failures should be %d, is %d", a+1, runReport.Failures)
    }
}

func TestAppendCanonicals(t *testing.T) {
    a := runReport.Canonicals
    appendCanonicals()
    if runReport.Canonicals != a+1 {
        log.Fatalf("Append canonicals should be %d, is %d", a+1, runReport.Canonicals)
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
            log.Fatal("Original is not \"asdf\" as it should be")

            replacement := fmt.Sprintf("%d-%d-%d-%d", i, i, i, i)
            if r.Canonical != replacement {
                log.Fatalf("Canonical is wrong. Expecting: %s, got: %s", replacement, r.Canonical)
            }
        }
    }
}

func TestSendMessageToGCM(t *testing.T) {
    // Test empty token
    ok, err := sendMessageToGCM("", "")
    if ok {
        log.Fatal("ok should be false")
    }
    if err.Error() != "Token was empty, exiting" {
        log.Fatalf("Unexpected error string: %s", err.Error())
    }

    // Test empty payload
    ok, err = sendMessageToGCM("asdf", "")
    if ok {
        log.Fatal("ok should be false")
    }
    if err.Error() != "Payload was empty, exiting" {
        log.Fatalf("Unexpected error string: %s", err.Error())
    }

    // Test bad json
    ok, err = sendMessageToGCM("asdf", "asdf")
    if ok {
        log.Fatal("ok should be false")
    }
    if !strings.HasPrefix(err.Error(), "invalid character") {
        log.Fatalf("Unexpected error string: %s", err.Error())
    }

    // Test bad send
    aOrig := runReport.Attempts
    fOrig := runReport.Failures

    ok, err = sendMessageToGCM("asdf", "{\"key\": \"value\"}")
    time.Sleep(5 * time.Second)
    if ok {
        log.Fatal("ok should be false")
    }
    if runReport.Attempts != aOrig+1 {
        log.Fatal("Attempts not incremented by 1")
    }
    if runReport.Failures != fOrig+1 {
        log.Fatal("Failures not incremented by 1")
    }
    if !strings.HasPrefix(err.Error(), "401 error") {
        log.Fatalf("Unexpected error string: %s", err.Error())
    }
}
