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

echo "Compiling migrate..."
GOOS=linux GOARCH=amd64 go build -o migrate -v -ldflags="-w -s" cmd/migrate/migrate.go

echo "Compiling apidocs..."
swag init -o web/static/apidocs --ot json --parseDependency --parseInternal --parseDepth 1 -q -g cmd/app/main.go

FILES=""

if [ "$2" == "migrate" ] && [ "$3" != "" ]; then
  echo "migrate $3 will be executed."
  MIGRATE="&& chmod +x ./migrate && ./migrate $3"
  FILES="./migrations ./migrate"
fi

echo "Enter password several times if asked, that's ok."

ssh -t "$TARGET" "sudo rm /opt/ws/app"

scp -r ./web ./app ./cmd/scripts/update.sh $FILES "$TARGET:/opt/ws"

ssh -t "$TARGET" "cd /opt/ws &&
    (sudo tmux kill-session -t about-web || pkill tmux);
    chmod +x ./update.sh && sudo ./update.sh &&
    sudo tmux new-session -s about-web -d \"sudo /opt/ws/app -config-path /opt/ws/config.yml \" $MIGRATE "

rm ./app
