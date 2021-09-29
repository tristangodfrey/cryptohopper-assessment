#!/usr/bin/env bash

if [ "$1" == "--help" ] || [ -z "$1" ]; then
  echo "Sends a request to the cryptohopper microservice (make sure to start the container first!)"
  echo "./client.sh <exchange> <pair> <period>"
  exit 1
fi

curl "http://localhost:5000/?exchange=$1&pair=$2&period=$3"