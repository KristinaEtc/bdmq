#!/bin/bash

#set -x
go build stomp.go

echo "starting"

pids=""

./stomp --ConfigPath=server.config --FileWithFrames="frames.txt" --ShowFrames=true &
pids+=" $!"
#./stomp --ConfigPath=client.config &
#pids+=" $!"

echo "$pids"

for p in $pids; do
        if wait $p; then
                echo "Process $p success"
        else
                echo "Process $p fail"
        fi
done

#set +x
