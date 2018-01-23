package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func main() {
	v := url.Values{}
	v.Set("ask", "grep")
	v.Add("option", "inr")
	v.Add("pattern", "amit")
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
