FROM alpine:3.12

RUN apk add --no-cache ca-certificates
COPY vultr-cli .
ENTRYPOINT ["./vultr-cli"]