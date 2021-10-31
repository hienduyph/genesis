FROM golang:1.17-buster AS compiler

WORKDIR /builder

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/tbb ./cmd/tbb

# final
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata tini

# make golang follow /etc/hosts https://github.com/golang/go/issues/22846
RUN echo "hosts: files dns" > /etc/nsswitch.conf

RUN rm -rf /etc/localtime\
  && cp /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.5 && \
  wget -qO /bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe

COPY --from=compiler /builder/build/* /bin/

WORKDIR /app

ENTRYPOINT ["tini", "--"]
