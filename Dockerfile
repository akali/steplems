FROM golang:1.21-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /bin/steplems

FROM alpine:3

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/steplems /bin/
COPY static /bin/static

WORKDIR /bin

CMD ["/bin/steplems"]