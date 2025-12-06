FROM golang:1.25-alpine AS builder

ARG problem=0

WORKDIR /usr/src/app

COPY . .
RUN go build -v -o /usr/local/bin/app ./$problem

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
CMD ["app"]
