FROM golang:alpine AS builder

RUN apk update 

RUN apk add --no-cache git && apk add make

WORKDIR /app

COPY . .

RUN make build

CMD ["sh", "-c", "make build"]

#CMD ["sh", "-c", "ls"]


# Step 2: Create the final image with Filebeat and the Golang application
FROM docker.elastic.co/beats/filebeat:8.15.0

# Set environment variables
ENV FILEBEAT_CONFIG_PATH=/usr/share/filebeat/filebeat.yml

# Copy the Golang application from the builder image
COPY --from=builder /app/main /app/main

# Copy the Config for the Golang Application from the builder image
COPY --from=builder /app/cmd/app/config.yaml /app/config.yaml

# Set the working directory to the app folder
WORKDIR /app

# Copy Filebeat configuration file (assumes you have a `filebeat.yml` file in your project)
COPY ./resources/go/filebeat.yml ${FILEBEAT_CONFIG_PATH}

# Expose necessary ports (update if your app requires specific ports)
EXPOSE 8008

# Set the user to root to allow for the creation of log files
USER root

CMD ["sh", "-c", "./main -configpath=./config.yaml -ipport=${HOSTNAME}:8008 -servicename=${HOSTNAME}service -nodename=${HOSTNAME}node & sleep 5 && filebeat -e -c ${FILEBEAT_CONFIG_PATH} -d '*'"]

