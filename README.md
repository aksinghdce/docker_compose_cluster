# Overall design of the system
## How I used docker for the project
Read a more dynamic document of this README : https://tinyurl.com/ybhw2kpc
. this document contains other information about algorithms and design decisions besides the ones mentioned in this document. Please feel free to comment in this document because you can.


# Assignment 2 built on top of Assignment 1
# Features
Assignment 1 : Distributed grep
Assignment 2 : Membership service

# Distributed grep

This is a Distributed Systems project developed with Docker and GoLang. The first out of two distributed services is "distributedgrep" (distributed grep). We will describe distributed grep service while discussing Membership service.

# Membership service

The service requires a leader to manage group membership list that reflects the state of the cluster. When a leader comes up it know that it has to lead based on its hostname.

When a non-leader comes up, it tries to send udp request to a node with hostname leader.assignment2 and expect to be added to the group.

There is only one membership group for the scope of this assignment.

Nodes(other than leader.assignment2) send "ADD" Request with request code: 1, leader.assignment2 responds with request code 2 to signify success of "ADD" request

The service implements a Finite State Machine with 3 possible states:

1. State 0: Machine has just begin to run : This state is transient and trivial.
2. State 1: Machine has assumed the role of a leader : This machine has hostname leader.assignment2
3. State 2: Machine will begin requesting the leader to add to the group
4. State 3: Non leaders will try to maintain the cluster membership info independent of the leader

This work is not 100% complete. We are yet to fix a defect. We will describe the defect further down this document.

![Design Diagram](https://github.com/aksinghdce/docker_compose_cluster/blob/assignment2/doc/images/Overall%20design%20of%20membership%20service.png)

# Runtime
## Launch
![Launch](https://github.com/aksinghdce/docker_compose_cluster/blob/master/doc/images/launch.PNG)
## No Failure use case
![No Failure use case](https://github.com/aksinghdce/docker_compose_cluster/blob/master/doc/images/No_Failure_Case.PNG)

2. Plan:
Need to fix a concurrency issue in the code. The issue can be seen when we try to update the group membership information after receiving the heartbeat from any peer. The error is basically a race condition with go maps.


_*Explanation of the issue: The nodes that are running in State 3 have a problem of coupling between the layer that is responsible for maintaining the membership map and a lower layer that is responsible for sending udp packets to the peers. When the node running in State 3 is trying to consolidate the information received from a peer into it's own data structure it complains :" Attempt to read and write map at the same time "*_

3. Test Plan:
There is a technical difficulty to impersonate a peer node in the test cases. Need to figure out an alternative.

## Getting Started

Run the following command from the parent directory where docker-compose.yml file is located

run:
"docker-compose up --build"

A cluster of docker containers is launched from docker-compose.yml configuration file. The hardcoded settings is to create
25 docker containers from the same set of files located in "http_server1" directory. We began working with 5 nodes and have been testing the setup informally with 10-25 nodes.

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
to test whether the container got added to the docker-compose cluster already running

## Running the tests

To run the tests on the cluster please follow the following steps:

### In a separate terminal (Don't use the terminal where you ran "docker-compose up --build") run "docker-compose ps"


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


### The five machines would be shown in the first column of the output command. Copy the name of one of the machines on which you want to run grep. You will need to use the content in your clipboard in the next step.

### You will see the following output:
#### This output also demonstrates distributed grep service functionality.

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
```
Note: I don't choose to maintain the outputs of the programs because they keep changing too very often; it makes the process of maintaining the outputs of the commands impractical.

The following test, tests the http server's handler:
```
PS C:\Users\aksin> docker exec dockercomposecluster_grepservice1_1 go test -v ./distributedgrepserver
```

# Test of the membership service
![Test of membership service](https://github.com/aksinghdce/docker_compose_cluster/blob/master/doc/images/Test.PNG)

### And coding style tests

This is work in progress...

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
