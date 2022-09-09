FROM public.ecr.aws/docker/library/golang:1.19-alpine

ARG problem=0

WORKDIR /usr/src/app

COPY . .
RUN go build -v -o /usr/local/bin/app ./$problem

CMD ["app"]
