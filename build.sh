#! /bin/bash

for i in {amd64,arm,arm64}
do
    env GOOS=linux GOARCH="$i" \
        go build -o gohfs-linux_"$i" gohfs.go
done
