FROM golang:1.17-alpine as builder

RUN  apk add --update \
     make

ENV CGO_ENABLED 0

RUN mkdir -p /out
RUN mkdir -p /go/src/github.com/vultr/vultr-cli
ADD . /go/src/github.com/vultr/vultr-cli

RUN cd /go/src/github.com/vultr/vultr-cli && \
    make builds/vultr-cli_linux_amd64


FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/github.com/vultr/vultr-cli/builds/* /
ENTRYPOINT ["/vultr-cli_linux_amd64"]