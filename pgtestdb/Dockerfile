FROM postgres:9.4.24-alpine


RUN mkdir -p "$PGDATA" && chmod 700 "$PGDATA"
ADD create.sql /docker-entrypoint-initdb.d