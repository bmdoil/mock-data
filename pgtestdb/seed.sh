#!/bin/bash

set -e

psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE foo (a int, b text);
EOSQL
    