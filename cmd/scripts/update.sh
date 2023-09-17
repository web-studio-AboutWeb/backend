#!/usr/bin/env bash


if [ `whoami` != root ]; then
    echo Please run this script as root or using sudo
    exit
fi

set -e
cd `dirname $0`


chown -R ws:ws /opt/ws
chmod -R 0600 /opt/ws
chmod -R u+rwX /opt/ws
chmod u+x /opt/ws/app
