rm -rf ./output/bin
mkdir -p output/bin
cd my_blog_server && go build -o ../output/bin/my_blog
cd ../output/bin && ./my_blog