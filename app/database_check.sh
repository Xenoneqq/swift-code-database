#!/bin/bash

until nc -z -v -w30 $DB_HOST $DB_PORT
do
    echo "Waiting for database $DB_HOST:$DB_PORT"
    sleep 5
done

echo "Database detected. Starting up application..."
./main