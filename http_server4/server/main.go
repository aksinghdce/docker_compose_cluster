package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
)

func main() {
	/**
	1. Get a grep request from peer, parse it
	2. Get a goroutine to get local grep
	3. Have the grep result from local? Send it
	**/
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
	c := localGrep(m.Get("ask"), m.Get("option"), m.Get("search"), m.Get("file"))
	fmt.Fprint(resWriter, <-c)
}

/*
Name: localGrep
Input: command, search pattern, filename
Output: Channel of strings that carries grep command output
*/
func localGrep(ask, option, search, file string) <-chan string {
	c := make(chan string)
	go func() {
		cmd := exec.Command(ask, option, search, file)
		stdOutStdErr, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		c <- string(stdOutStdErr)
	}()
	return c
}
