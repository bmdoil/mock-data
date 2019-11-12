#!/bin/sh

while ! nc -z postgres 5432; do
  sleep 0.1
done

echo "PostgreSQL started"

/mockd postgres -h ${PGHOST} -p ${PGPORT} -d ${PGUSER} -t ${TABLE} -n 10
