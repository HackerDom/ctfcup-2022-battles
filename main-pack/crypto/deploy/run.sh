#!/bin/sh

while true; do
    socat TCP-LISTEN:31337,reuseaddr,fork EXEC:"python3 -u /tmp/server.py"
done
