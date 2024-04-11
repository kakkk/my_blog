#!/usr/bin/env bash
BUILD_DIR=$(cd $(dirname $0); pwd)
rm -rf $BUILD_DIR/output
cd $BUILD_DIR/my_blog_pages && yarn build
cd $BUILD_DIR/my_blog_admin && yarn build
cd $BUILD_DIR/my_blog_server && bash build.sh


mkdir -p $BUILD_DIR/output/assets
mkdir -p $BUILD_DIR/output/templates
mkdir -p $BUILD_DIR/output/admin
mkdir -p $BUILD_DIR/output/conf
cp -r $BUILD_DIR/conf/* $BUILD_DIR/output/conf/
cp -r $BUILD_DIR/my_blog_pages/src/templates/* $BUILD_DIR/output/templates/
cp -r $BUILD_DIR/my_blog_pages/dist/* $BUILD_DIR/output/assets/
cp -r $BUILD_DIR/my_blog_admin/dist/* $BUILD_DIR/output/admin/
cp -r $BUILD_DIR/my_blog_server/output/* $BUILD_DIR/output/