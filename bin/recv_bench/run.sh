#!/bin/bash

set -e

# apt install xxx

#go run main.go $@
#wget -O recv_bench https://oss.otz.app:62135/public/recv_bench
curl -o recv_bench https://oss.otz.app:62135/public/recv_bench
chmod +x recv_bench

./recv_bench $@
