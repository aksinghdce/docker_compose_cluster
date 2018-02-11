Read a more dynamic document of this README : https://tinyurl.com/ybhw2kpc

# Features
1.Distributed grep
2.Membership service

# Distributed grep

This is a Distributed Systems project developed with Docker and GoLang. The distributed service being implemented in this project is "distributedgrep" (distributed grep).

# Membership service

The service requires a leader to manage group membership list that reflects
the state of the cluster. When a leader comes up it know that it has to lead
based on its hostname.

When a non-leader comes up, it tries to ping(send udp request) to the leader
and expect to be added to the group.

# TO-DO design

1. Now: 
  1. Run the leader: 
  ```
  > docker exec dockercomposecluster_grepservice1_1 go run ./multicastheartbeatser
ver/multicastheartbeatserver.go
  
  My hostname:leader.assignment2
  Interface: eth0
  Interface Flag: up|broadcast|multicast
  Network: ip
  multicast address 0 : 224.0.0.1
  ```

  2. Run a non-leader: 
  ```
  docker exec dockercomposecluster_grepservice2_1 go run ./multicastheartbeatser
ver/multicastheartbeatserver.go
My hostname:node2.assignment2
I am not a leader. I am too old to serve. I will just die
  ```
  3. Run a udp client at a non-leader
  ```
  docker exec dockercomposecluster_grepservice5_1 multicastheartbeater
  ```
  
  The effect can be seen at the console of the multicastheartbeatserver:
  ```
  From 172.20.0.6:10002 Data received: 0
Network: ip
multicast address 0 : 224.0.0.1
From 172.20.0.6:10002 Data received: 1
Network: ip
multicast address 0 : 224.0.0.1
From 172.20.0.6:10002 Data received: 2
Network: ip
multicast address 0 : 224.0.0.1
From 172.20.0.6:10002 Data received: 3
Network: ip
multicast address 0 : 224.0.0.1
From 172.20.0.6:10002 Data received: 4
Network: ip
multicast address 0 : 224.0.0.1
From 172.20.0.6:10002 Data received: 5
Network: ip
multicast address 0 : 224.0.0.1
From 172.20.0.6:10002 Data received: 6
Network: ip
multicast address 0 : 224.0.0.1
From 172.20.0.6:10002 Data received: 7
Network: ip
multicast address 0 : 224.0.0.1
From 172.20.0.6:10002 Data received: 8
Network: ip
multicast address 0 : 224.0.0.1
From 172.20.0.6:10002 Data received: 9
  ```
2. Plan:
3. Test Plan:
In response the leader checks if the node is a new node.

## Getting Started

Run the following command from the parent directory where docker-compose.yml file is located

run:
"docker-compose up --build"

A cluster of docker containers is launched from docker-compose.yml configuration file. The hardcoded settings is to create
5 docker containers from the same set of files located in "http_server1" directory.

### Prerequisites

In order to test the system you need to have docker installed on your machine.
If you have a windows machine where docker can not be run (if you have a non Windows 10 Pro), then install docker on a linux VM
OR
Install older version of docker for your windows. You will be using a virtual machine hypervisor instead of using OS level container features in windows.

You WILL need to install docker-compose separately after installing docker.

## Running the tests

To run the tests on the cluster please follow the following steps:

1. In a separate terminal (Don't use the terminal where you ran "docker-compose up --build") run "docker-compose ps"


```
PS C:\Users\aksin\go\src\docker_compose_cluster> docker-compose ps
               Name                          Command          State    Ports
------------------------------------------------------------------------------
dockercomposecluster_grepservice1_1   distributedgrepserver   Up      8080/tcp
dockercomposecluster_grepservice2_1   distributedgrepserver   Up      8080/tcp
dockercomposecluster_grepservice3_1   distributedgrepserver   Up      8080/tcp
dockercomposecluster_grepservice4_1   distributedgrepserver   Up      8080/tcp
dockercomposecluster_grepservice5_1   distributedgrepserver   Up      8080/tcp
```


2. The five machines would be shown in the first column of the output command. Copy the name of one of the machines on which you want to run grep. You will need to use the content in your clipboard in the next step.

3. You will see the following output:


```
PS C:\Users\aksin\go\src\docker_compose_cluster> docker exec dockercomposecluster_grepservice1_1 distributedgrep grep -c 8080 Dockerfile
Node's hostname: {grepservice1}
Node's hostname: {grepservice2}
Node's hostname: {grepservice3}
Node's hostname: {grepservice4}
Node's hostname: {grepservice5}
Commandstring: [grep -c 8080 Dockerfile]
Remote grep 1

Commandstring: [grep -c 8080 Dockerfile]
Remote grep 1

Commandstring: [grep -c 8080 Dockerfile]
Remote grep 1

Commandstring: [grep -c 8080 Dockerfile]
Remote grep 1

Commandstring: [grep -c 8080 Dockerfile]
Remote grep 1

1
```

So, effectively you have grepped for "8080" in a file called "Dockerfile" and got a count with "-c" option of grep command. The command was run on the machine "dockercomposecluster_grepservice1_1" and it showed the grep result from the whole cluster.

### Break down into end to end tests

The test cases are written for all the packages developed for this project. [Some test cases are work in progress]
The following test, tests the local node's grepping responsibility:
```
PS C:\Users\aksin\go\src\docker_compose_cluster> docker exec dockercomposecluster_grepservice1_1 go test -v ./utilities
=== RUN   TestCluster
--- PASS: TestCluster (0.00s)
=== RUN   TestLocalGrep
=== RUN   TestLocalGrep/exporting_8080_grepped
=== RUN   TestLocalGrep/local_log_file_creation_grepped
--- PASS: TestLocalGrep (0.00s)
    --- PASS: TestLocalGrep/exporting_8080_grepped (0.00s)
    --- PASS: TestLocalGrep/local_log_file_creation_grepped (0.00s)
=== RUN   ExampleLocalGrep
--- PASS: ExampleLocalGrep (0.00s)
PASS
ok      app/utilities   0.008s
```

The following test, tests the http server's handler:
```
PS C:\Users\aksin> docker exec dockercomposecluster_grepservice1_1 go test -v ./distributedgrepserver
=== RUN   TestCommandHandler
--- PASS: TestCommandHandler (0.00s)
PASS
ok      app/distributedgrepserver       0.005s
```


### And coding style tests

Explain what these tests test and why

```
Give an example
```

## Built With

* [Docker](https://www.docker.com/) - Docker is world's leading software containerization platform
* [Go](https://golang.org/) - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.

## Authors

* **Amit Kumar Singh** - *Initial work* - [Aksinghdce](https://github.com/Aksinghdce)
* **Vimal Philip** - *Initial work* - [Vimalphilip](https://github.com/vimalphilip)

See also the list of [contributors](https://github.com/aksinghdce/docker_compose_cluster/contributors) who participated in this project.

## License

This project is copyrighted. Please write to me before you plan to use this project.

## Acknowledgments

* Thanks Prof. Sathish for providing the guiding words and support
* I have used https://tinyurl.com/yaxrl5ea template for organizing this document
