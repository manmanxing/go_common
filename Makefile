#!/usr/bin/env bash

all: build

fmt:
	goimports -l -w  ./

install:  clean

clean:

	rm -rf output/conf/

build: install
	go build  ./...