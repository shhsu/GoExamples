FROM alpine
ADD . /builds
RUN apk add --update bash && rm -rf /var/cache/apk/*
ENTRYPOINT echo "hello world"
