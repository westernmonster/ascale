FROM gcr.io/ascale-439911/golang:1.23.2-alpine3.20 as builder

WORKDIR /go/src/ascale
COPY ./ .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -o cmd app/api/cmd/main.go

FROM gcr.io/ascale-439911/base:latest
COPY --from=builder /go/src/ascale/cmd /go/bin/cmd
COPY --from=builder /go/src/ascale/assets /go/bin/assets
COPY --from=builder /go/src/ascale/database /go/bin/database
WORKDIR /go/bin
CMD ["/go/bin/cmd","-c=/go/bin/config/config.toml"]


