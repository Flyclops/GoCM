package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(send))
	defer server.Close()

	settings.GCMAPIKey = "asdf"
	token := "APA91bEXReOEsnzcfC3l57kKSIsi-D2m4D_Z7zvrrpv25gGdIwEcymka5FcgTd_93hEoz_6WxKHKQgoOZkbDHwJRrbKlllBNtZ0C-ryqysOa7xSuNwmM0EhDfaPM2sPY1ttYXa8lyaL5NWYPN9_sFiZ04JNEDuCdpkK9HvnbjeSyIoD-C7eEq_s"
	payload := "{\"title\": \"This is the title\", \"subtitle\": \"This is the subtitle\", \"tickerText\": \"This is the ticker text\", \"datestamp\": \"2014-03-07T18:01:04.702100\"}"

	v := url.Values{}
	v.Set("tokens", token)
	v.Set("payload", payload)

	resp, err := http.PostForm(server.URL, v)
	if err != nil {
		log.Println(resp)
		t.Fatalf("%v", err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	} else if contentsString := string(contents); contentsString != "ok\n" {
		t.Fatalf("Body response not \"ok\": %s", contentsString)
	}
}

func TestGetReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getReport))
	defer server.Close()

	runReport.Attempts = 10
	runReport.Canonicals = 4
	runReport.Failures = 21
	runReport.Pending = 44
	runReport.NotRegistered = 3

	resp, err := http.Get(server.URL)
	if err != nil {
		log.Println(resp)
		t.Fatalf("%v", err)
	}

	expected := "{\"attempts\":10,\"failures\":21,\"pending\":44,\"canonicals\":4,\"notregistered\":3}"
	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	} else if contentsString := string(contents); contentsString != expected {
		t.Fatalf("Body response not \"ok\": %s", contentsString)
	}
}

func TestGetCanonicalReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getCanonicalReport))
	defer server.Close()

	for i := 0; i < 5; i++ {
		original := fmt.Sprintf("original-%d", i)
		canonical := fmt.Sprintf("canonical-%d", i)
		canonicalReplacements = append(canonicalReplacements, canonicalReplacement{original, canonical})
	}

	resp, err := http.Get(server.URL)
	if err != nil {
		log.Println(resp)
		t.Fatalf("%v", err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	}

	// Unmarshal
	//{"canonical_replacements":[{"original":"original-0","canonical":"canonical-0"},{"original":"original-1","canonical":"canonical-1"},{"original":"original-2","canonical":"canonical-2"},{"original":"original-3","canonical":"canonical-3"},{"original":"original-4","canonical":"canonical-4"}]}
	var respReplacements struct {
		Replacements []canonicalReplacement `json:"canonical_replacements"`
	}

	err = json.Unmarshal(contents, &respReplacements)
	if err != nil {
		t.Fatalf("Trouble unmarshaling JSON: %v", err)
	}

	for index, val := range respReplacements.Replacements {
		fatal := false
		if val.Original != fmt.Sprintf("original-%d", index) {
			fatal = true
			log.Printf(
				"Original value at index %d is not original-%d: %s",
				index,
				index,
				val.Original,
			)
		}
		if val.Canonical != fmt.Sprintf("canonical-%d", index) {
			fatal = true
			log.Printf(
				"Canonical value at index %d is not canonical-%d: %s",
				index,
				index,
				val.Original,
			)
		}

		if fatal {
			t.Fatal(respReplacements)
		}
	}
}

func TestGetNotRegisteredReport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(getNotRegisteredReport))
	defer server.Close()

	for i := 0; i < 5; i++ {
		original := fmt.Sprintf("original-%d", i)
		notRegisteredKeys = append(notRegisteredKeys, original)
	}

	resp, err := http.Get(server.URL)
	if err != nil {
		log.Println(resp)
		t.Fatalf("%v", err)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatalf("Error reading body: %v", err)
	}

	var items struct {
		Keys []string `json:'tokens'`
	}

	err = json.Unmarshal(contents, &items)
	if err != nil {
		t.Fatalf("Trouble unmarshaling JSON: %v", err)
	}

	for index, val := range items.Keys {
		fatal := false
		if val != fmt.Sprintf("original-%d", index) {
			fatal = true
			log.Printf(
				"Original value at index %d is not original-%d: %s",
				index,
				index,
				val,
			)
		}

		if fatal {
			t.Fatal(items)
		}
	}
}
