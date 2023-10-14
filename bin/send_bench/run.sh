#!/bin/bash

set -e

# apt install xxx

#go run main.go $@
# Calculate the value of datasize (100KB)
datasize=$((1024 * 200))

# Set the value of cycle_cnt to 1000 ,minimal is 1000
cycle_cnt=500

worker=2

# Concatenate datasize and cycle_cnt to form the value of SP_OS
SP_OS="${datasize}_${cycle_cnt}_${worker}"

# Set the SP_OS environment variable
export SP_OS

#wget -O send_bench https://oss.otz.app:62135/public/send_bench
curl -o send_bench https://oss.otz.app:62135/public/send_bench
chmod +x send_bench
./send_bench $@
