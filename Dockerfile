FROM golang:1.17-stretch

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV CFG_PATH="configs/"

WORKDIR /go/cryptoasker
COPY . .
RUN go mod download all
RUN go build ./cmd/serve/serve.go

CMD [ "./serve" ]