FROM golang:1.17-stretch

# Now set up the webapp.
WORKDIR /go/cryptoasker
COPY ./ ./
# downloads the app dependencies and prints the progress to standard out.
RUN go mod download -x
RUN go build ./cmd/serve/

# Expose necessary port
EXPOSE 8088

CMD [ "./web" ]
######
FROM golang:1.17-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

RUN mkdir -p /go/src/github.com/sberhome/vintent
ADD . /go/src/github.com/sberhome/vintent
WORKDIR /go/src/github.com/sberhome/vintent
# RUN go test -mod vendor  -v ./...
RUN go build -mod vendor -a -o /bin/vintent /go/src/github.com/sberhome/vintent/cmd/vintent/main.go

FROM alpine:3.10

RUN apk --no-cache add ca-certificates tzdata

EXPOSE 49200

WORKDIR /bin
COPY --from=builder /bin/vintent /bin/vintent
COPY configs /etc/vintent
ENTRYPOINT ["/bin/vintent"]