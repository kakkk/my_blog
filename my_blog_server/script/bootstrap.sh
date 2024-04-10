#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=kakkk.my_blog.service
export ROOT_PATH=$CURDIR
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}