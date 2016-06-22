var $ = require('gulp-load-plugins')();
var gulp         = require('gulp');
var rimraf          = require('rimraf');
var sequence = require('run-sequence');


gulp.task('clean',function(done) {
  //Delete our old css files
  rimraf('static/css/**/*', done);
});

// Compile SCSS files to CSS
gulp.task('scss', function() {
    return gulp.src('assets/sass/styles.scss')
        .pipe($.sourcemaps.init())
        .pipe($.sass())
        .on('error', $.sass.logError)
        .pipe($.sourcemaps.write())
        .pipe(gulp.dest('static/css'));
});

gulp.task('default', function(done) {
  sequence('clean', 'scss', done);
});
