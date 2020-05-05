#!/bin/bash

# Create a platform code for opensolar on openx
curl --location --request POST 'http://localhost:8080/admin/platform/new' \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --data-urlencode 'username=admin' \
    --data-urlencode 'token=pmkjMEnyeUpdTyhdHElkBExEKeLIlYft' \
    --data-urlencode 'name=opensolar' \
    --data-urlencode 'code=CODE' \
    --data-urlencode 'timeout=false'