#!/bin/bash

yum install epel-release golang -y
go env -w GOPROXY=https://goproxy.cn,direct
