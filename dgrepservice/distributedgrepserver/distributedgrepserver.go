package main

import (
	"app/utilities"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	/*
		Create a log file on every node of the cluster for logging
		user requests.
	*/
	ctx := context.Background()
	startTime := time.Now()
	utilities.Log(ctx, startTime.String())
	/**
	1. Get a grep request from peer, parse it
	2. Get a goroutine to get local grep
	3. Have the grep result from local? Send it
	**/
	http.ListenAndServe(":8080", handler())
}

func handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/grep", utilities.DecorateWithLog(commandHandler))
	// TO-DO : Add a handler function (with a specification comment)
	// similar to the commandHandler to send data from peers to MembershipManager
	// send appropriate data about the peer to the membership service to service
	// 3 kinds of events, as described in the assignment statement.
	return r
}

/*
Specification
*/
func commandHandler(resWriter http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	ctx = context.WithValue(ctx, int(42), rand.Int63())

	/* The following 5 lines of code is to test that whether I
	am using the standard way of passing values from the client.

	The code prints "Missiong request form value" when run
	We will take a look at this later to see if we need to
	get the grep command through form
	*/
	text := r.FormValue("grep")
	if text == "" {
		utilities.Log(ctx, "Missiong request form value")
	}
	utilities.Log(ctx, "test in Request Form: %s", text)
	/* The above lines of code is to test that whether I
	am using the standard way of passing values from the client.

	The code prints "Missiong request form value" when run
	We will take a look at this later to see if we need to
	get the grep command through form
	*/

	bodyBuff, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	m, err := url.ParseQuery(string(bodyBuff))
	if err != nil {
		log.Fatal(err)
	}

	grepCommand := m.Get("grep")
	grepresult := utilities.LocalGrep(strings.Split(grepCommand, " "))
	utilities.Log(ctx, grepCommand)
	fmt.Fprint(resWriter, grepresult)
}
