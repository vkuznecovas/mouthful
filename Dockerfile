# build
FROM golang:alpine AS build

ENV APP mouthful
ENV REPO vkuznecovas/$APP

RUN apk update
RUN apk add --update build-base nodejs nodejs-npm git openssh 
RUN go get -u github.com/golang/dep/cmd/dep

ADD . /${GOPATH}/src/github.com/${REPO}/
WORKDIR /${GOPATH}/src/github.com/${REPO}/
RUN sh ./build.sh

# run
FROM alpine:latest

ENV APP mouthful
ENV REPO vkuznecovas/$APP

COPY --from=build /go/src/github.com/${REPO}/dist /${APP}
WORKDIR /${APP}
VOLUME /${APP}/_mouthful_db/mouthful.db
RUN ls -la
CMD ["./mouthful"] 