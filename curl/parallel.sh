#!/bin/bash
execute_curl() {
    curl --location 'http://localhost:8080/getVisits'
}

export -f execute_curl

# Run the curl command 100000 times in parallel
seq 100000 | xargs -P 10 -I {} bash -c 'execute_curl'
