FROM golang:1.10-alpine
ARG MOUTHFUL_VER
ENV CGO_ENABLED=${CGO_ENABLED:-1} \
    GOOS=${GOOS:-linux} \
    MOUTHFUL_VER=${MOUTHFUL_VER:-master}
RUN set -ex; \
    apk add --no-cache bash build-base curl git && \
    echo "http://dl-cdn.alpinelinux.org/alpine/edge/community" >> /etc/apk/repositories && \
    echo "http://dl-cdn.alpinelinux.org/alpine/edge/main" >> /etc/apk/repositories && \
    apk add --no-cache upx nodejs nodejs-npm && \
    go get -d github.com/vkuznecovas/mouthful && \
    go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/vkuznecovas/mouthful
RUN git checkout $MOUTHFUL_VER && \
    ./build.sh && \
    cd dist/ && \
    upx --best mouthful

FROM alpine:3.7
COPY --from=0 /go/src/github.com/vkuznecovas/mouthful/dist/ /app/
# this is needed if we're using ssl
RUN apk add --no-cache ca-certificates
WORKDIR /app/
VOLUME [ "/app/data" ]
EXPOSE 8080
CMD ["/app/mouthful"]
