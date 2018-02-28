package cluster

import (
	"app/utilities"
	"log"
	"net/url"
	"strings"
)

type node struct {
	hostname string
}

/*
Specification:
This is a grep command that runs on a node.
The node might be remote or local, eitherway
the call reaches the server via REST api
and response from the REST server comes likewise

Input: Command string for grep, as received from
the user terminal

Output: A channel of string. The information on this channel
is collected by the Cluster (plays the role of an aggregator)
*/
func (n *node) Grep(commandstrings []string) <-chan string {
	if len(commandstrings) <= 0 {
		log.Fatal("Grep command invalid")
	}
	v := url.Values{}
	v.Add("grep", strings.Join(commandstrings, " "))
	return utilities.RemoteGrep(n.hostname, v)
}


/*
Specification: Defines the functions a cluster can
perform.

We will just keep this code to remind us
that we need to get a good grasp on go interfaces.
*/
type clustercan interface {
	Grep(commandstring ...string) string
}

/*
A cluster entity on which we will run our distributed
commands. A cluster is just a collection of nodes.
*/
type Cluster struct {
	nodes []node
}

/*
Build a new cluster from a configuration file.
*/
func (c *Cluster) NewCluster(configFileName string) {
	l := utilities.ReadConfig(configFileName)

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		if str, ok := e.Value.(string); ok {
			node := node{hostname: strings.Trim(str, "\n")}
			c.nodes = append(c.nodes, node)
		} else {
			log.Fatal("The server names file doesn't have strings")
		}
	}

}

/*
Specification: This function is to grep on the cluster.
grep is one of the commands that can be run on the cluster

The function fans out the same command to the members
of the cluster, collects their responses, makes the
responses organized for user to view and sends the
response to the user.
*/
func (c *Cluster) Grep(commandstring []string) string {
	lg := ""
	/*
		Run grep on all the nodes, collect the results and send
		back in lumpsum
	*/
	for _, node := range c.nodes {
		nodeGrepResult := <-node.Grep(commandstring)
		//fmt.Println("Remote grep", nodeGrepResult)
		lg += "\n\n"
		lg += node.hostname
		lg += ":\n"
		lg += nodeGrepResult
	}
	return lg
}