#!/usr/bin/env bash

echo -e '🔨\033[32m build comment... \033[0m'
cd my_blog_comment && yarn build
cd ..
echo -e '🔨\033[32m build pages... \033[0m'
cd my_blog_pages && yarn build
cd ..
echo -e '🔨\033[32m build admin... \033[0m'
cd my_blog_admin && yarn build
cd ..
rm -rf output
mkdir -p output/bin
mkdir -p output/assets
mkdir -p output/templates
mkdir -p output/admin
mkdir -p output/conf
cp script/* output/
cp conf/* output/conf/
cp my_blog_pages/src/templates/* output/templates/
cp my_blog_comment/dist/comment.min.js output/assets/comment.min.js
cp my_blog_pages/dist/* output/assets/
cp my_blog_admin/dist/* output/admin/
echo -e '🔨\033[32m build server... \033[0m'
cd my_blog_server && go build -o ../output/bin/my_blog
echo -e '🔨\033[32m build success!!! \033[0m'