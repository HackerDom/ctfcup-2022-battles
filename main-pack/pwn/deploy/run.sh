#!/bin/sh

while true; do
    socat TCP-LISTEN:31337,reuseaddr,fork EXEC:"/tmp/ld-linux-x86-64.so.2 /tmp/server"
done
