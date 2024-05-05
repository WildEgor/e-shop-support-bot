#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER dbuser WITH ENCRYPTED PASSWORD 'dbuser';
	CREATE DATABASE dbuser OWNER dbuser;

	CREATE USER postgres WITH ENCRYPTED PASSWORD 'postgres';
	CREATE DATABASE e_shop_support_db OWNER postgres;
EOSQL