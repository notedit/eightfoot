#!/usr/bin/env sh

DB_HOST=localhost
DB_USER=user
DB_NAME=database

psql -h$DB_HOST -U$DB_USER -d$DB_NAME -q < database.sql
psql -h$DB_HOST -U$DB_USER -d$DB_NAME -q < init_data.sql
