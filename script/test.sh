#!/usr/bin/env bash

./script/build.sh
export APP_SERVER_PORT=9010
export APP_DATABASE_NAME=iot-server-test
./build/server-iot --create-db
./build/server-iot &
postman collection run 14947205-ced8c886-6ab5-4f9d-ae5f-e17b27c3737e -e 14947205-fe8d1501-53df-40b2-a36c-21aac0eb8220 --env-var "baseUrl=http://0.0.0.0:$APP_SERVER_PORT"
fuser -k $APP_SERVER_PORT/tcp
