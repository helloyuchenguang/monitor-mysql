#!/usr/bin/env nu
let main_file = "./main/main.go"
echo "building..." $main_file
rm -rf main.exe
go build -ldflags '-w -s' $main_file
upx etcd-cli main.exe