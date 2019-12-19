FROM golang:1.13

WORKDIR /go/src/app
#COPY . .
ADD . .

#RUN rm -R docs
RUN mv src/* ../
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
