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
This program returns grep results from local machine. There is only one log file that is grepped
the name of the log file that gets grepped is machine2.log for local and machine1.log for remote
*/
func main() {

	/*
		Launch a goroutine to act as a server for peer's grep requests
	*/
	/**
	1. Launch a go routine to get local grep
	2. Launch a go routine that uses a FanIn function to get peer grep
	3. We have all the grep results, send it to client
	**/

	lc := localGrep("grep", "506901129", "machine2.log")
	//Get grep result from remote machines
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

	c := remoteGrep("grepservice1", v)
	c1 := remoteGrep("grepservice2", v1)
	c2 := remoteGrep("grepclient", v2)
	c3 := remoteGrep("grepservice4", v3)
	fmt.Println("Response from local machine:", <-lc)
	fmt.Println("Response from grepservice1:", <-c)
	fmt.Println("Response from grepservice2:", <-c1)
	fmt.Println("Response from grepclient:", <-c2)
	fmt.Println("Response from grepservice4:", <-c3)
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

/*
Name: remoteGrep
Input: machine's hostname to be grepped, grep command
Output: A channel that receives remote grep output
*/
func remoteGrep(machine string, cmd url.Values) <-chan string {
	c := make(chan string)
	go func() {
		resp, err := http.PostForm("http://"+machine+":8080/", cmd)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			log.Println("ERROR: sending request to remote http server", machine)
			c <- "Error connecting to remote host"
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading response from remote")
		}
		c <- string(body)
	}()
	return c
}
