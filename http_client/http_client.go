package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {

	/**
	1. Launch a go routine to get local grep
	2. Launch a go routine that uses a FanIn function to get peer grep
	3. We have all the grep results, send it to client
	**/

	fmt.Println("Local Grep: ", localgrep("grep", "506901129", "machine2.log"))
	v := url.Values{}
	v.Set("ask", "grep")
	v.Add("search", "tanuki")
	v.Add("file", "machine1.log")

	v1 := url.Values{}
	v1.Set("ask", "grep")
	v1.Add("search", "0.149471")
	v1.Add("file", "machine1.log")

	v2 := url.Values{}
	v2.Set("ask", "grep")
	v2.Add("search", "1.113414")
	v2.Add("file", "machine1.log")

	v3 := url.Values{}
	v3.Set("ask", "grep")
	v3.Add("search", "3161")
	v3.Add("file", "machine1.log")

	c := remoteGrep("grepservice", v)
	c1 := remoteGrep("grepservice1", v1)
	c2 := remoteGrep("grepservice2", v2)
	c3 := remoteGrep("grepservice3", v3)
	fmt.Println("Response from server:", <-c)
	fmt.Println("Response from server:", <-c1)
	fmt.Println("Response from server:", <-c2)
	fmt.Println("Response from server:", <-c3)
}

func remoteGrep(machine string, cmd url.Values) <-chan string {
	c := make(chan string)
	go func() {
		resp, err := http.PostForm("http://"+machine+":8080/", cmd)
		if err != nil {
			log.Fatal("ERROR: sending request to remote http server")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading response from remote")
		}
		c <- string(body)
	}()
	return c
}
