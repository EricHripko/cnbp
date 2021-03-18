FROM golang:1.15

LABEL moby.buildkit.frontend.network.none="true"

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["cnbp-frontend"]
