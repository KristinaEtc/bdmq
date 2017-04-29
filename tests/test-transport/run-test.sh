#!/bin/bash

function run {

        pids=""

        #process_name1=""
        #process_name2=""

        names=(server client)
        ./transport --ConfigPath=server.config --ShowFrames=true &
        pids+=" $!"
        #process_name1= $(ps -p "$!" -o comm)

        ./transport --ConfigPath=client.config &
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
}

#set -x
if !(go build transport.go); then
        echo "building failed"
else
        echo "starting"
        run
fi
