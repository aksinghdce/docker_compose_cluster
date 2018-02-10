package utilities

import (
	"bufio"
	"container/list"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

/*
Name: localGrep
Input: command, search pattern, filename
Output: Channel of strings that carries grep command output
*/
func LocalGrep(arguments []string) string {
	cmd := exec.Command(arguments[0], arguments[1:]...)
	stdOutStdErr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	return string(stdOutStdErr)
}

/*
Name: remoteGrep
Input: machine's hostname to be grepped, grep command
Output: A channel that receives remote grep output
*/
func RemoteGrep(machine string, cmd url.Values) <-chan string {
	c := make(chan string)
	go func() {
		req, err := http.NewRequest("POST", "http://"+machine+":8080/grep", strings.NewReader(cmd.Encode()))
		ctx := context.Background()
		// Don't wait for more than a second to get the grep result from remote server
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		req2 := req.WithContext(ctx)
		resp, err := http.DefaultClient.Do(req2)
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
		l.PushBack(scanner.Text())
	}
	return l
}
