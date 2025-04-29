#!/bin/bash

if [ "${MODE}" != "TEST" ]; then
    echo "Tests will not be performed"
    exit 0
fi

# Waiting for app

APP_URL="http://app:8080/api"
TIMEOUT=60
INTERVAL=2
TIMEPASSED=0

echo "Waiting for swift app..."

while true; do
    res=$(curl -s -o /dev/null -w "%{http_code}" "$APP_URL")
    if [ "$res" -eq 200 ]; then
        echo "Application is ready!"
        break
    fi

    echo "Waiting for swift app - time : ${TIMEPASSED}s"
    sleep "$INTERVAL"
    TIMEPASSED=$((TIMEPASSED + INTERVAL))

    if [ "$TIMEPASSED" -ge "$TIMEOUT" ]; then
        echo "Timeout reached ${TIMEOUT}! App did not launch."
        exit 1
    fi
done

# Launching tests

echo "Starting tests..."
go test -v ./tests/...

if [ $? -eq 0 ]; then
    echo "All tests passed!"
else
    echo "Some tests failed."
    exit 1
fi