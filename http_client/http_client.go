package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {

	/**
	1. Launch a go routine to get local grep
	2. Launch a go routine that uses a FanIn function to get peer grep
	3. We have all the grep results, send it to client
	**/

	v := url.Values{}
	v.Set("ask", "grep")
	v.Add("search", "tanuki")
	v.Add("file", "machine1.log")
	resp, err := http.PostForm("http://grepservice:8080/", v)
	if err != nil {
		fmt.Println("There was an error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		fmt.Println("response:", string(body))
	}
}
