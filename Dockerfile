FROM golang:1.8

WORKDIR /go/src/go-example-api
COPY . .

EXPOSE 8080

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-example-api"]

