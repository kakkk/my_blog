#!/usr/bin/env bash
BUILD_DIR=$(cd $(dirname $0); pwd)
rm -rf $BUILD_DIR/output

echo -e '🔨\033[32m build comment... \033[0m'
cd $BUILD_DIR/my_blog_comment && yarn build
echo -e '🔨\033[32m build pages... \033[0m'
cd $BUILD_DIR/my_blog_pages && yarn build
echo -e '🔨\033[32m build admin... \033[0m'
cd $BUILD_DIR/my_blog_admin && yarn build
echo -e '🔨\033[32m build server... \033[0m'
cd $BUILD_DIR/my_blog_server && bash build.sh
echo -e '🔨\033[32m copy... \033[0m'
mkdir -p $BUILD_DIR/output/assets
mkdir -p $BUILD_DIR/output/templates
mkdir -p $BUILD_DIR/output/admin
mkdir -p $BUILD_DIR/output/conf
cp -r $BUILD_DIR/conf/* $BUILD_DIR/output/conf/
cp $BUILD_DIR/my_blog_comment/dist/comment.min.js $BUILD_DIR/output/assets/comment.min.js
cp -r $BUILD_DIR/my_blog_pages/src/templates/* $BUILD_DIR/output/templates/
cp -r $BUILD_DIR/my_blog_pages/dist/* $BUILD_DIR/output/assets/
cp -r $BUILD_DIR/my_blog_admin/dist/* $BUILD_DIR/output/admin/
cp -r $BUILD_DIR/my_blog_server/output/* $BUILD_DIR/output/

echo -e '🔨\033[32m build success!!! \033[0m'