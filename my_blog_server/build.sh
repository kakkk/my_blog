#!/bin/bash
RUN_NAME=kakkk.my_blog.service
mkdir -p output/bin
cp script/* output 2>/dev/null
chmod +x output/bootstrap.sh
go build -o output/bin/${RUN_NAME}