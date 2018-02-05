package utilities

import (
	"log"
	"net/url"
	"strings"
)

type node struct {
	hostname string
}

func (n *node) Grep(commandstrings []string) <- chan string {
	if len(commandstrings) <= 0 {
		log.Fatal("Grep command invalid")
	}
	v := url.Values{}
	v.Add("grep", strings.Join(commandstrings, " "))
	return RemoteGrep(n.hostname, v)
}

type grepper interface {
	Grep(commandstring ...string) string
}

type Cluster struct {
	local node
	nodes []node
}


func (c *Cluster) NewCluster(configFileName string) {
	l := ReadConfig(configFileName)
	
	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		if str, ok := e.Value.(string); ok {
			node := node{hostname: str}
			c.nodes = append(c.nodes, node)
		} else {
			log.Fatal("The server names file doesn't have strings")
		}
	}

}

func (c *Cluster) Grep(commandstring []string) string {
	lg := ""
	for _, node := range(c.nodes) {
		nodeGrepResult := <-node.Grep(commandstring)
		//fmt.Println("Remote grep", nodeGrepResult)
		lg += "\n\n"
		lg += node.hostname
		lg += ":\n"
		lg += nodeGrepResult
	}
	return lg
}