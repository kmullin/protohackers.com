FROM golang:1.19-alpine as BUILDER

ARG problem=0

WORKDIR /usr/src/app

COPY . .
RUN go build -v -o /usr/local/bin/app ./$problem

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=BUILDER /usr/local/bin/app /usr/local/bin/app
CMD ["app"]
