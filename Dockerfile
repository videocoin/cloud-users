FROM golang:1.12.4 as builder
WORKDIR /go/src/github.com/videocoin/cloud-users
COPY . .
RUN make build

FROM bitnami/minideb:jessie
COPY --from=builder /go/src/github.com/videocoin/cloud-users/bin/users /opt/videocoin/bin/users
CMD ["/opt/videocoin/bin/users"]
