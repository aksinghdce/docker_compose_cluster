Read a more dynamic document of this README : https://tinyurl.com/ybhw2kpc
# Assignment 2 built on top of Assignment 1
# Features
Assignment 1 : Distributed grep
Assignment 2 : Membership service

# Distributed grep

This is a Distributed Systems project developed with Docker and GoLang. The distributed service being implemented in this project is "distributedgrep" (distributed grep).

# Membership service

![Design of Membership service](https://github.com/aksinghdce/docker_compose_cluster/blob/assignment2/doc/images/Overall%20design%20of%20membership%20service.png)

The service requires a leader to manage group membership list that reflects
the state of the cluster. When a leader comes up it know that it has to lead
based on its hostname.

When a non-leader comes up, it tries to send udp request to a node with 
hostname leader.assignment2 and expect to be added to the group.

There is only one group for the scope of this assignment.

Nodes(other than leader.assignment2) send "ADD" Request with request code: 1
leader.assignment2 responds with request code 2 to signify success of "ADD" request

The service implements a Finite State Machine with 3 possible states:

1. State 0: Machine has just begin to run
2. State 1: Machine has assumed the role of a leader : This machine has hostname leader.assignment2
3. State 2: Machine will begin requesting the leader to add to the group
4. State 3: Non leaders will try to maintain the cluster membership info independent of the leader

State 3 is not fully implemented yet.
This work is not 100% complete. 

# Runtime
```
distributedservice5_1  | STOPPING ADD REQUEST NOW
distributedservice5_1  | Running in state 3 now
distributedservice4_1  | STOPPING ADD REQUEST NOW
distributedservice4_1  | Running in state 3 now
distributedservice3_1  | STOPPING ADD REQUEST NOW
distributedservice3_1  | Running in state 3 now
distributedservice2_1  | STOPPING ADD REQUEST NOW
distributedservice2_1  | Running in state 3 now
distributedservice1_1  | Received Lower stack:{172.20.0.2 {[] 107868 1}}
distributedservice1_1  | Received ADD request:{172.20.0.2 {[] 107868 1}}
distributedservice1_1  | First time saw:172.20.0.2
distributedservice1_1  | GroupInfo:[     172.20.0.2]
distributedservice1_1  | Received Lower stack:{172.20.0.2 {[] 108167 1}}
distributedservice1_1  | Received ADD request:{172.20.0.2 {[] 108167 1}}
distributedservice1_1  | First time saw:172.20.0.2
distributedservice1_1  | GroupInfo:[     172.20.0.2]
distributedservice1_1  | Received Lower stack:{172.20.0.5 {[] 7590 1}}
distributedservice1_1  | Received ADD request:{172.20.0.5 {[] 7590 1}}
distributedservice1_1  | First time saw:172.20.0.5
distributedservice1_1  | GroupInfo:[     172.20.0.2 172.20.0.5]
distributedservice1_1  | Received Lower stack:{172.20.0.5 {[] 7597 1}}
distributedservice1_1  | Received ADD request:{172.20.0.5 {[] 7597 1}}
distributedservice1_1  | First time saw:172.20.0.5
distributedservice1_1  | GroupInfo:[     172.20.0.5 172.20.0.2]
distributedservice1_1  | Received Lower stack:{172.20.0.4 {[] 22002 1}}
distributedservice1_1  | Received ADD request:{172.20.0.4 {[] 22002 1}}
distributedservice1_1  | First time saw:172.20.0.4
distributedservice1_1  | GroupInfo:[     172.20.0.2 172.20.0.5 172.20.0.4]
distributedservice1_1  | Received Lower stack:{172.20.0.3 {[] 41411 1}}
distributedservice1_1  | Received ADD request:{172.20.0.3 {[] 41411 1}}
distributedservice1_1  | First time saw:172.20.0.3
distributedservice1_1  | GroupInfo:[     172.20.0.2 172.20.0.5 172.20.0.4 172.20.0.3]
distributedservice1_1  | Received Lower stack:{172.20.0.4 {[] 22008 1}}
distributedservice1_1  | Received ADD request:{172.20.0.4 {[] 22008 1}}
distributedservice1_1  | First time saw:172.20.0.4
distributedservice1_1  | GroupInfo:[     172.20.0.4 172.20.0.3 172.20.0.2 172.20.0.5]
distributedservice1_1  | Received Lower stack:{172.20.0.3 {[] 41416 1}}
distributedservice1_1  | Received ADD request:{172.20.0.3 {[] 41416 1}}
distributedservice1_1  | First time saw:172.20.0.3
distributedservice1_1  | GroupInfo:[     172.20.0.4 172.20.0.3 172.20.0.2 172.20.0.5]
distributedservice5_1  | Cluster Info:[     172.20.0.4 172.20.0.3 172.20.0.2 172.20.0.5]
distributedservice2_1  | Cluster Info:[     172.20.0.4 172.20.0.3 172.20.0.2 172.20.0.5]
distributedservice3_1  | Cluster Info:[     172.20.0.4 172.20.0.3 172.20.0.2 172.20.0.5]
distributedservice4_1  | Cluster Info:[     172.20.0.4 172.20.0.3 172.20.0.2 172.20.0.5]

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

## Manage running containers with docker-compose

### To add a new container 
In order to add a new container to the cluster we need to follow the follow steps:
1. Run the cluster with existing configuration in daemon mode:
```
docker-compose up --build -d
```
2. Add the changes in docker-compose.yml file to reflect a new container
```
docker-compose up --no-recreate -d
```

You can do a 
```
docker-compose ps
```
to test whether the container got added

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
