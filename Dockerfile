FROM alpine:3 as certs
RUN apk --no-cache add ca-certificates

FROM golang:1.21-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/akali/steplems

COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0

RUN go get -v all

COPY . .

RUN go build -o /bin/steplems

FROM alpine

WORKDIR /bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/steplems /bin/steplems
CMD ["/bin/steplems"]
