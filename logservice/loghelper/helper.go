package loghelper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"net/http"
	"log"

	"bitbucket.org/asnegovoy-dataart-projects/jaeger-rd/entity"
)

var logserviceURL = flag.String("logservice", "http://logservice:5000",
	"Address of the logging service")

var tr = http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}
var client = &http.Client{Transport: &tr}

func WriteEntry(entry *entity.LogEntry) {
    var buf bytes.Buffer

    // Encode the entry into the buffer
    encoder := json.NewEncoder(&buf)
    if err := encoder.Encode(entry); err != nil {
        log.Printf("Failed to encode entry: %v", err)
        return
    }

    // Create the HTTP request
    req, err := http.NewRequest(http.MethodPost, *logserviceURL, &buf)
    if err != nil {
        log.Printf("Failed to create request: %v", err)
        return
    }

    // Make the HTTP request
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Failed to send request: %v", err)
        return
    }
    defer resp.Body.Close()

    // Optionally, you can check the response status
    if resp.StatusCode != http.StatusOK {
        log.Printf("Request failed with status: %v", resp.Status)
    }
}