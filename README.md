
# Chat App Backend

A web application that serves as a load-balanced backend for a Chat Application, orchestrated through Docker Compose. Each Chat Application has a Filebeat sidecar application run within the same container that retrieves logs and puts it within their own indices on elasticsearch, with their commonality being that they always start with the "chat-" prefix.

# Features

- Round Robin Load Balancing through NGINX for chat servers

- User Registration with BCrypt password encryption.

- One - to - One chat among friends that spans multiple chat server instances.

- Saving and paginated retrieval of messages within the database.

- Logging through Filebeat as a Sidecar, Elasticsearch and Kibana.

- Add and Remove Friends.

- Retrieval of old Chats.

# TODO

- Group Chat feature

- Highly Available Nginx Load Balancer + Reverse Proxy + Configuration for beyond two chat servers

- Highly Available, Sharded and Partitioned Databases

- High-Availability setup of Elasticsearch

- Implement Secure Websocket

- Better Logging

- Prove that this solution is valid for 60000 users connected per server.

- Refactoring

# Tech Stack

- Golang 1.19

- Echo

- Redis

- MySQL

- NGINX

- Docker

- Elasticsearch

- Kibana

- Filebeat

# Setup

## Installing Dependencies

The demonstration requires Docker Desktop. Install Docker Desktop here (https://www.docker.com/products/docker-desktop/)

## Running the app

Simply run this on the working directory:

- docker compose up

This will orchestrate the backend servers, the load balancer and the necessary data stores.

