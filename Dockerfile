FROM golang:1.23.2-alpine3.20 as builder
RUN apk --update add --no-cache ca-certificates openssl git tzdata && \
update-ca-certificates

WORKDIR /go/src/ascale
COPY ./ .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o cmd app/api/cmd/main.go

FROM scratch
LABEL maintainer="Jack"
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/src/ascale/cmd /go/bin/cmd
COPY --from=builder /go/src/ascale/assets /go/bin/assets
COPY --from=builder /go/src/ascale/database /go/bin/database
WORKDIR /go/bin
CMD ["/go/bin/cmd","-c=/go/bin/config/config.toml"]


