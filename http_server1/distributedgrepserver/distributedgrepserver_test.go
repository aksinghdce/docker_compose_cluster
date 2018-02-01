package main

import (
	"net/http"
	"testing"
	"strings"
	"app/utilities"
	"context"
	"net/url"
	"net/http/httptest"
)

func TestCommandHandler(t *testing.T) {
	//See cluster.go +19
	cmd := url.Values{}
	cmdstrings := []string{"grep", "-c", "8080", "/go/src/app/Dockerfile"}
	cmd.Add("grep", strings.Join(cmdstrings, " "))
	req, err := http.NewRequest("POST", "http://localhost:8080/", strings.NewReader(cmd.Encode()))
	if err != nil {
		utilities.Log(context.Background(), "Error creating test request")
	}

	//We will send this request to the handler we are testing
	/*This information is something I would need to memorize :D */
	// httptest.NewRecorder() does the same as http.ResponseWriter does 
	// They both implement the same interface
	rec := httptest.NewRecorder()
	commandHandler(rec, req)
}