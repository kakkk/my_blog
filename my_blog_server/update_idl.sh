#!/usr/bin/env bash

hz update -idl idl/page.thrift
hz update -idl idl/api.thrift --snake_tag
thriftgo -g go -o biz/model idl/common.thrift