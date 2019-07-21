#!/bin/sh
if [ -e /app/config.yml ]
then
    echo "using user config"
else
    cp /config.yml /app/config.yml
fi
cd /app
erproxy -c config.yml