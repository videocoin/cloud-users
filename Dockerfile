FROM alpine:3.7

COPY bin/users /opt/videocoin/bin/users

CMD ["/opt/videocoin/bin/users"]
