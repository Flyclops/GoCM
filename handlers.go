package main

import (
    "encoding/json"
    "io"
    "net/http"
    "strconv"
)

// Send a message to GCM
func send(w http.ResponseWriter, r *http.Request) {
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

// Return a run report for this process
func getReport(w http.ResponseWriter, r *http.Request) {
    runReportWG.Wait()
    runReportWG.Add(1)
    a, _ := json.Marshal(runReport)
    runReportWG.Done()
    b := string(a)
    io.WriteString(w, b)
}

// Return all currently collected canonical reports from GCM
func getCanonicalReport(w http.ResponseWriter, r *http.Request) {
    ids := map[string][]canonicalReplacement{"canonical_replacements": canonicalReplacements}
    a, _ := json.Marshal(ids)
    b := string(a)
    io.WriteString(w, b)

    // Clear out canonicals
    go func() {
        canonicalReplacementsWG.Wait()
        canonicalReplacementsWG.Add(1)
        canonicalReplacements = nil
        canonicalReplacementsWG.Done()
    }()

    go func() {
        runReportWG.Wait()
        runReportWG.Add(1)
        runReport.Canonicals = 0
        runReportWG.Done()
    }()
}
