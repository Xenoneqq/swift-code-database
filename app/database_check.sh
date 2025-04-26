#!/bin/bash

echo "Checking for database $DB_HOST:$DB_PORT"

until nc -z -v -w30 $DB_HOST $DB_PORT
do
    echo "Waiting for database $DB_HOST:$DB_PORT"
    sleep 2
done

echo "Database detected. Starting up application..."
./main