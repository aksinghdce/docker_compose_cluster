Read a more dynamic document of this README : https://docs.google.com/document/d/1LneRwSc1iNG_EZPtJHI6JIRKjVAt3nt08X6PgevXOic/edit?usp=sharing

# docker_compose_cluster
A cluster of docker containers one of which is a server and rest are clients

In order to test the system you need to have docker installed on your machine.
If you have a windows machine where docker can not be run (if you have a non Windows 10 Pro), then install docker on a linux VM

Run the following command from the parent directory where docker-compose.yml file is located

run:
"docker-compose up"

You might need to install docker-compose separately.

In this setup services are initialized on separate containers that communicate with each other over a network that is created by 
docker-compose.

