FROM golang:1.10-alpine
ARG MOUTHFUL_VER
ENV CGO_ENABLED=${CGO_ENABLED:-1} \
    GOOS=${GOOS:-linux} \
    MOUTHFUL_VER=${MOUTHFUL_VER:-1.0.2}
RUN set -ex; \
    apk add --no-cache bash build-base curl git && \
    echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    apk add --no-cache upx nodejs nodejs-npm && \
    go get -d github.com/vkuznecovas/mouthful && \
    curl -sSL https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 -o /usr/bin/dep && \
    chmod +x /usr/bin/dep
WORKDIR /go/src/github.com/vkuznecovas/mouthful
RUN git checkout $MOUTHFUL_VER && \
    ./build.sh && \
    cd dist/ && \
    upx --best mouthful

FROM alpine:3.7
COPY --from=0 /go/src/github.com/vkuznecovas/mouthful/dist/ /app/
WORKDIR /app/
VOLUME [ "/app/data" ]
EXPOSE 8080
CMD ["/app/mouthful"]
