#!/usr/bin/env bash

kubectl create secret generic cloudsql-instance-credentials --from-file=credentials.json=[PROXY_KEY_FILE_PATH]

kubectl create secret generic cloudsql-db-credentials --from-literal=dsn=[DSN]

