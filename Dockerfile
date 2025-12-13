FROM golang:1.25-trixie AS builder

ARG problem=0

WORKDIR /usr/src/app

# Download Go modules
COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -race -o /usr/local/bin/app ./$problem

FROM debian:trixie
# RUN apk --no-cache add ca-certificates
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
CMD ["app"]
