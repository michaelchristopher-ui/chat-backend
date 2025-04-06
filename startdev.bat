go build ./cmd/app/main.go
start /B main.exe -configpath="./cmd/app/config-dev.yaml" -ipport="127.0.0.1:8008" -servicename="dummyservice" -nodename="dummynode0"
