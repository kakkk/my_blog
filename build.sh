#!/usr/bin/env bash
BUILD_DIR=$(cd $(dirname $0); pwd)
rm -rf $BUILD_DIR/output
rm -rf $BUILD_DIR/my_blog_comment/dist
rm -rf $BUILD_DIR/my_blog_pages/dist
rm -rf $BUILD_DIR/my_blog_admin/dist
rm -rf $BUILD_DIR/my_blog_server/output

# 整点花里胡哨没用的，主打一个心情愉悦 :)

# 评论构建
start_time=$(date +%s)
echo -e "🔨\033[32m Build comment... \033[0m"
cd $BUILD_DIR/my_blog_comment && yarn install >> /dev/null 2>&1 && yarn build >> /dev/null 2>&1

# 页面构建
start_time=$(date +%s)
echo -e "🔨\033[32m Build pages... \033[0m"
cd $BUILD_DIR/my_blog_pages && yarn install >> /dev/null 2>&1 && yarn build >> /dev/null 2>&1

# 后台构建
echo -e "🔨\033[32m Build admin... \033[0m"
cd $BUILD_DIR/my_blog_admin && yarn install >> /dev/null 2>&1 && yarn build >> /dev/null 2>&1

# 后端服务构建
echo -e "🔨\033[32m Build server... \033[0m"
cd $BUILD_DIR/my_blog_server && bash build.sh >> /dev/null 2>&1


# 复制结果
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

end_time=$(date +%s)
cost_time=$[ $end_time-$start_time ]
echo -e "✨ \033[32m Done in $(($cost_time))s \033[0m"