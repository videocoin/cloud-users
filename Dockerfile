FROM alpine:3.7

COPY bin/vc-user /opt/videocoin/bin/vc-user

CMD ["/opt/videocoin/bin/vc-user"]
