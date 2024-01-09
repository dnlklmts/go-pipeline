FROM golang:1.21.4-alpine AS compiling_stage
RUN mkdir -p /go/src/pipeline
WORKDIR /go/src/pipeline
COPY *.go ./
COPY go.mod ./
RUN go install .
 
FROM alpine:latest
LABEL version="1.0.0"
LABEL maintainer="Kolomeets Daniil<kolomeetz.dany@yandex.ru>"
WORKDIR /root/
COPY --from=compiling_stage /go/bin/pipeline .
ENTRYPOINT ./pipeline
