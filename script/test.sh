#!/usr/bin/env bash

./script/build.sh
export APP_SERVER_PORT=9010
export APP_DATABASE_NAME=iot-server-version-2-test
export APP_SERVER_ENV=test
./build/server-iot --create-db
./build/server-iot &
postman collection run 14947205-55a4862a-6046-4c4f-90a4-0553b554fbd6 -e 14947205-fe8d1501-53df-40b2-a36c-21aac0eb8220 --env-var "baseUrl=http://0.0.0.0:$APP_SERVER_PORT"
fuser -k $APP_SERVER_PORT/tcp
