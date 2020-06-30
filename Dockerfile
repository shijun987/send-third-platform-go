FROM alpine:3.11.3

WORKDIR $GOPATH/src/send-third-platform-henan

COPY send-third-platform .

CMD ./send-third-platform