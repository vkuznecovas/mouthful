FROM node:14-alpine AS node

FROM golang:1.15.2-alpine AS golang
ENV CGO_ENABLED=${CGO_ENABLED:-1} \
    GOOS=${GOOS:-linux} \
    MOUTHFUL_VER=${MOUTHFUL_VER:-master}
COPY --from=node /usr/lib /usr/lib
COPY --from=node /usr/local/share /usr/local/share
COPY --from=node /usr/local/lib /usr/local/lib
COPY --from=node /usr/local/include /usr/local/include
COPY --from=node /usr/local/bin /usr/local/bin

RUN set -ex; \
    apk add --no-cache bash build-base curl git && \
    echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    apk add --no-cache python3 python2 

WORKDIR /go/src/github.com/vkuznecovas/mouthful
COPY . /go/src/github.com/vkuznecovas/mouthful
RUN ./build.sh

FROM alpine:3.20.3
COPY --from=golang /go/src/github.com/vkuznecovas/mouthful/dist/ /app/
# this is needed if we're using ssl
RUN apk add --no-cache ca-certificates
WORKDIR /app/
VOLUME [ "/app/data" ]
EXPOSE 8080
CMD ["/app/mouthful"]
