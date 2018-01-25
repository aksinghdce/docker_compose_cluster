package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
	cmd := "grep"
	search := "tanuki"
	logFile := "machine1.log"
	lc := localGrep(cmd, search, logFile)
	fmt.Println("Response from local machine:", <-lc)
	//Get grep result from remote machines
	v := url.Values{}
	v.Set("ask", cmd)
	v.Add("search", search)
	v.Add("file", logFile)
	l := readLine("text.txt")
	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		//fmt.Println(e.Value)
		if str, ok := e.Value.(string); ok {
			/* act on str */
			c := remoteGrep(str, v)
			fmt.Println("Response from grepservice5:", <-c)
		} else {
			/* not string */
			fmt.Println("The server names file doesn't have strings")
		}

	}

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

func readLine(path string) *list.List {
	l := list.New()
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		l.PushBack(scanner.Text())
	}
	return l
}
