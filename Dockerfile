FROM alpine:latest

ADD shorten /bin/
RUN apk -Uuv add ca-certificates
ENTRYPOINT /bin/shorten
