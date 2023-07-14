#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE domain_hq_test;
	GRANT ALL PRIVILEGES ON DATABASE domain_hq_test TO "$POSTGRES_USER";
EOSQL