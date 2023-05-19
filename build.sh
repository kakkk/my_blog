#!/usr/bin/env bash

rm -rf output
mkdir -p output/bin
mkdir -p output/assets/css
mkdir -p output/assets/js
mkdir -p output/assets/katex
mkdir -p output/templates
mkdir -p output/conf/
mkdir -p output/templates/
cp script/* output/
cp conf/* output/conf/
cp my_blog_front/templates/* output/templates/
cp my_blog_front/assets/css/stylesheet.min.css output/assets/css/stylesheet.min.css
cp my_blog_front/assets/js/highlight.min.js output/assets/js/highlight.min.js
cp -r my_blog_front/assets/katex/* output/assets/katex/
cd my_blog_server && go build -o ../output/bin/my_blog