# ProgettoDP

![alt text](architecture.png "Architecture")

## Prerequisites

- docker

- docker compose

- make

## Steps to run

- make update_proto (compile protobuf definitions)

- make build (build containers)

- make up (or make upd for detached version)

- make down (to stop containers and clean the system from volumes and networks)

## Environment Variables

In the .env file it is possible to edit the environment variables (such as the number of replicas per service).