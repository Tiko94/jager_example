package loghelper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"net/http"

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
	enc := json.NewEncoder(&buf)

	if err := enc.Encode(entry); err != nil {
		// Handle the error here, e.g., log it or return it
		log.Printf("Failed to encode entry: %v", err)
	}
	req, _ := http.NewRequest(http.MethodPost, *logserviceURL, &buf)
	client.Do(req)
	resp, err := client.Do(req)
	if err != nil {
    	// Handle the error here, e.g., log it or return it
    	log.Printf("Failed to send request: %v", err)
	}
}
