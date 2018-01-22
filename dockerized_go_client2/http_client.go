package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.Get("http://grepservice:8080/")
	if err != nil {
		fmt.Println("There was an error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		fmt.Println("response:", string(body))
	}
}
