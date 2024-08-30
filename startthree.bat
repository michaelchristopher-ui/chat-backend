go build ./cmd/app/main.go
start /B main.exe -configpath="./cmd/app/config.yaml" -ipport="127.0.0.1:8008" -servicename="dummyservice" -nodename="dummynode0"
start /B main.exe -configpath="./cmd/app/config.yaml" -ipport="127.0.0.1:8009" -servicename="dummyservice" -nodename="dummynode1"
start /B main.exe -configpath="./cmd/app/config.yaml" -ipport="127.0.0.1:8010" -servicename="dummyservice" -nodename="dummynode2"