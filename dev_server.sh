BUILD_DIR=$(cd $(dirname $0); pwd)
rm -rf ./output/bin
cd $BUILD_DIR/my_blog_server && bash build.sh
cp -r $BUILD_DIR/my_blog_server/output/* $BUILD_DIR/output/