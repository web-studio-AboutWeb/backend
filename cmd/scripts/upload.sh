#!/usr/bin/env bash

set -e
cd "$(dirname "$0")/../.."
pwd

TARGET=$1
if [ -z "$TARGET" ]; then
    printf "Syntax: \n\t$0 user@addr\n"
    exit 1
fi

echo "Compiling application..."
GOOS=linux GOARCH=amd64 go build -o app -v -ldflags="-w -s" cmd/app/main.go

echo "Compiling apidocs..."
swag init -o web/static/apidocs --ot json -q -g cmd/app/main.go

echo "Enter password several times if asked, that's ok."

scp -r web app cmd/scripts/update.sh "$TARGET:/opt/ws/"

ssh -t "$TARGET" "cd /opt/ws &&
    (tmux kill-session -t about-web || pkill tmux);
    chmod +x ./update.sh && sudo ./update.sh &&
    tmux new-session -s about-web -d \"cd /opt/ws && ./app \""

rm ./app
