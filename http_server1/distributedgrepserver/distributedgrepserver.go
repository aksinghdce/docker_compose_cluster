package main

import (
	"app/utilities"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	/**
	1. Get a grep request from peer, parse it
	2. Get a goroutine to get local grep
	3. Have the grep result from local? Send it
	**/
	// Declare a context object
	http.HandleFunc("/", commandHandler)
	//nil as second argument meand we are using DefaultServeMux
	http.ListenAndServe(":8080", nil)
}

func commandHandler(resWriter http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBuff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	m, err := url.ParseQuery(string(bodyBuff))
	if err != nil {
		log.Fatal(err)
	}

	grepresult := utilities.LocalGrep(strings.Split(m.Get("grep"), " "))
	fmt.Fprint(resWriter, grepresult)
}
