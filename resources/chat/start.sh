#!/bin/sh

# Start the Go application
echo "Starting the Go application..."
./main -configpath=./config.yaml -ipport=${HOSTNAME}:8008 -servicename=${HOSTNAME}service -nodename=${HOSTNAME}node &

# Wait a few seconds for the Go application to start
sleep 5

# Start Filebeat
echo "Starting Filebeat..."
filebeat -e -c ${FILEBEAT_CONFIG_PATH} -d "*"
