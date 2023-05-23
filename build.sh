#!/usr/bin/env bash

cd my_blog_pages && yarn build
cd ..
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
cp my_blog_pages/dist/* output/assets/
cp my_blog_admin/dist/* output/admin/
cd my_blog_server && go build -o ../output/bin/my_blog