package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
)

func main() {
	http.HandleFunc("/", commandHandler)
	//nil as second argument meand we are using DefaultServeMux
	http.ListenAndServe(":8080", nil)
}

func commandHandler(resWriter http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBuff, err := ioutil.ReadAll(r.Body)
	if err == nil {
		m, err := url.ParseQuery(string(bodyBuff))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(reflect.TypeOf(m))
		fmt.Println(m)
		fmt.Println(m.Get("ask"))
		fmt.Println(m.Get("option"))
		fmt.Println(m.Get("pattern"))
		//fmt.Println("Request form:", bodyBuff)
	}
	fmt.Fprint(resWriter, "Hello from commandHandler")
}
