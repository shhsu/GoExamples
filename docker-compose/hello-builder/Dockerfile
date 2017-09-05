FROM alpine
ADD . /builds
RUN apk add --update bash && rm -rf /var/cache/apk/* && /bin/bash /builds/run-me.sh
ENTRYPOINT echo "hello world"
