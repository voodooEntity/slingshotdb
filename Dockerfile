FROM golang:1.13

WORKDIR /go/src/app
ADD . .

RUN mv src/* ../
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
