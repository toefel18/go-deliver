#!/usr/bin/env bash

go build main.go
sudo PORT=443 APP_SOURCES=/home/ubuntu/src/github.com/toefel18/go-deliver/DHLApp ./main
