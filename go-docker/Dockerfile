FROM golang:alpine

RUN apk update 

RUN apk add --no-cache git && apk add make

WORKDIR /app

COPY . .

EXPOSE 8008

CMD "sh" "-c" "make build-and-run-linux"



