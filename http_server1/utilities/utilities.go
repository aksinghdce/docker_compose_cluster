package utilities

import (
	"bufio"
	"container/list"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

/*
Name: localGrep
Input: command, search pattern, filename
Output: Channel of strings that carries grep command output
*/
func LocalGrep(ask, search, file string) <-chan string {
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
func RemoteGrep(machine string, cmd url.Values) <-chan string {
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

func ReadConfig(path string) *list.List {
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
