#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE DATABASE vanderbot;
	GRANT ALL PRIVILEGES ON DATABASE vanderbot TO postgres;
	CREATE DATABASE vanderbot_shadow;
	GRANT ALL PRIVILEGES ON DATABASE vanderbot_shadow TO postgres;
EOSQL