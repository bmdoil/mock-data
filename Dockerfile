#   build stage
FROM golang:alpine AS test-env
RUN apk update && apk upgrade && apk add --no-cache git
RUN go get -v \
    github.com/bmdoil/mock-data \
    github.com/icrowley/fake \
    github.com/vbauerster/mpb \
    github.com/lib/pq \
    github.com/op/go-logging
ENV APP_DIR=$GOPATH/src/mock-data/
RUN mkdir -p ${APP_DIR}
COPY . ${APP_DIR}
WORKDIR ${APP_DIR}
CMD ["go", "test"]


FROM test-env AS build-env

RUN GOOS=linux GOARCH=amd64 \
        go build -o /mockd
ENTRYPOINT /mockd

FROM alpine:3.7 as prod-env
COPY --from=build-env /mockd /
ENTRYPOINT /mockd