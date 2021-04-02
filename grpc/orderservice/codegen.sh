#!/bin/bash

# This script serves as an example to demonstrate how to generate the gRPC-Go
# interface and the related messages from .proto file.
#
# It assumes the installation of i) Google proto buffer compiler at
# https://github.com/google/protobuf (after v2.6.1) and ii) the Go codegen
# plugin at https://github.com/golang/protobuf (after 2015-02-20). If you have
# not, please install them first.
#
# We recommend running this script at $GOPATH/src.
#
# If this is not what you need, feel free to make your own scripts. Again, this
# script is for demonstration purpose.

WORKSPACE=$(cd $(dirname $0)/; pwd)
cd $WORKSPACE
files=`ls ./*.proto`
echo $files

#--proto_path 等同于 -I
#当前目录下或 GOPATH/src/目录下 寻找待编译的 .proto 文件
#protoc -I ./ -I $GOPATH/src/ --go_out=plugins=grpc:. $files
#protoc --proto_path ./ --proto_path $GOPATH/src/ --go_out=plugins=grpc:. $files
#只在当前目录下编译
protoc --go_out=plugins=grpc:. $files