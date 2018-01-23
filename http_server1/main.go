package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
)

/*
A http server that listens for grep requests from peers
Calls a handler to fetch the local grep and return the results
*/
func main() {
	http.HandleFunc("/", commandHandler)
	//nil as second argument meand we are using DefaultServeMux
	http.ListenAndServe(":8080", nil)
}

/*
the http server's request handler for "/" endpoint
*/
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
	c := localGrep(m.Get("ask"), m.Get("search"), m.Get("file"))

	fmt.Fprint(resWriter, <-c)
}

/*
Name: localGrep
Input: command, search pattern, filename
Output: Channel of strings that carries grep command output
*/
func localGrep(ask, search, file string) <-chan string {
	c := make(chan string)
	go func() {
		cmd := exec.Command(ask, search, file)
		stdOutStdErr, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		c <- string(stdOutStdErr)
	}()
	return c
}
