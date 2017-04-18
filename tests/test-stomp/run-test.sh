#!/bin/bash

#set -x
go build stomp.go command-processer.go commands-endpoints.go

echo "starting"

pids=""
#process_name1=""
#process_name2=""

names=(server client)
./stomp --ConfigPath=server.config --ShowFrames=true &
pids+=" $!"
#process_name1= $(ps -p "$!" -o comm)

./stomp --ConfigPath=client.config &
pids+=" $!"
#process_name2= $(ps -p "$!" -o comm)

echo "$pids"
#echo "$process_name1"
#echo "$process_name2"

i=0
for p in $pids; do
        if wait $p; then
                echo "Process $p ${names[i]} success"
        else
                echo "Process $p ${names[i]} fail"
        fi
        i=$((var+1))
done

#set +x
