var gulp = require('gulp');
var concat = require('gulp-concat');                            //- 多个文件合并为一个；
var minifyCss = require('gulp-minify-css');                     //- 压缩CSS为一行；
var notify = require('gulp-notify');

gulp.task('default', function (done) {                                //- 创建一个名为 concat 的 task
    gulp.src(['./common/*.css', './core/*.css', './extended/*.css', './hljs/*.css', './includes/*.css'])    //- 需要处理的css文件，放到一个字符串数组里
        .pipe(concat('stylesheet.min.css'))                            //- 合并后的文件名
        .pipe(minifyCss())                                      //- 压缩处理成一行
        .pipe(gulp.dest('./'))                               //- 输出文件本地
        .pipe(notify({ message: 'finish' }));
    done();
});
