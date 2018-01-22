package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", commandHandler)
	//nil as second argument meand we are using DefaultServeMux
	http.ListenAndServe(":8080", nil)
}

func commandHandler(resWriter http.ResponseWriter, r *http.Request) {
	fmt.Fprint(resWriter, "Hello from commandHandler")
}
