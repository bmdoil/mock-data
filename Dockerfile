#   build stage
FROM golang:alpine AS dev_img
RUN apk update && apk upgrade && apk add --no-cache git && apk add netcat-openbsd
RUN go get -v \
    github.com/bmdoil/mock-data \
    github.com/icrowley/fake \
    github.com/vbauerster/mpb \
    github.com/lib/pq \
    github.com/op/go-logging

ENV APP_DIR=$GOPATH/src/mockd
RUN mkdir -p ${APP_DIR}
COPY ./ ${APP_DIR}
WORKDIR ${APP_DIR}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build -gcflags "all=-N -l" -o /mockd
#ENTRYPOINT /mockd


FROM alpine:3.7 AS prod_img

COPY --from=dev_img /mockd /usr/bin/

ARG PGPORT
ENV PGPORT $PGPORT
ARG PGUSER
ENV PGUSER $PGUSER
ARG PGPASS
ENV PGPASS $PGPASS
ARG PGDATABASE
ENV PGDATABASE $PGDATABASE
ARG PGHOST
ENV PGHOST $PGHOST
ARG TABLE
ENV TABLE $TABLE
ARG ENGINE
ENV ENGINE $ENGINE

CMD ["sh", "-c", "mockd $ENGINE -p $PGPORT -h postgres -u $PGUSER -d $PGDATABASE -w $PGPASS -x -n 10"]

