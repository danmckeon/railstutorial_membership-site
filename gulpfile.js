var $ = require('gulp-load-plugins')();
var rimraf = require('rimraf');
var sequence = require('run-sequence');

var browser = require('browser-sync');
var gulp = require('gulp');

gulp.task('clean', function(done) {
    //Delete our old css files
    rimraf('{public,static/css}/**/*', done);
});

// Compile SCSS files to CSS
gulp.task('styles', function() {
    return gulp.src('src/scss/styles.scss')
        .pipe($.sourcemaps.init())
        .pipe($.sass())
        .on('error', $.sass.logError)
        .pipe($.sourcemaps.write())
        .pipe(gulp.dest('static/css'));
});

gulp.task('build', function(done) {
  sequence('clean', 'styles', done);
});

// Start a server with LiveReload to preview the site in
gulp.task('server', function() {
  browser.init({
    server: 'public/',
    port: 8000
  });
});

gulp.task('default', ['build', 'server'], function() {
  gulp.watch('src/scss/**/*', ['styles', browser.reload]);
  gulp.watch('public/**/*', ['styles', browser.reload]);
});
