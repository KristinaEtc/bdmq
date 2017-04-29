#!/bin/bash

set -x
set -e
set -o pipefail

if !(exec test-transport/run-test.sh); then
        echo "transport failed"
else
        (exec test-stomp/run-test.sh)
fi

set +x
