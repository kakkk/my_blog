#!/usr/bin/env bash

rm -rf output
mkdir -p output/bin
mkdir -p output/assets/css
mkdir -p output/assets/js
mkdir -p output/templates
mkdir -p output/conf/
mkdir -p output/templates/
cp script/* output/
cp conf/* output/conf/
cp my_blog_front/templates/* output/templates/
cp my_blog_front/assets/css/stylesheet.min.css output/assets/css/stylesheet.min.css
cp my_blog_front/assets/css/katex/katex.min.css output/assets/css/katex.min.css
cp my_blog_front/assets/js/highlight.min.js output/assets/js/highlight.min.js
cp my_blog_front/assets/js/katex.min.js output/assets/js/katex.min.js
cp my_blog_front/assets/js/auto-render.min.js output/assets/js/auto-render.min.js
cp my_blog_server/biz/mock/mock_post.md output/mock_post.md
cd my_blog_server && go build -o ../output/bin/my_blog