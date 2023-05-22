#!/usr/bin/env bash

cd my_blog_pages && yarn build
cd ..
rm -rf output
mkdir -p output/bin
mkdir -p output/assets
mkdir -p output/templates
mkdir -p output/conf/
cp script/* output/
cp conf/* output/conf/
cp my_blog_pages/src/templates/* output/templates/
cp my_blog_pages/dist/* output/assets/
cd my_blog_server && go build -o ../output/bin/my_blog