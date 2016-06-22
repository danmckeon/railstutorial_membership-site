var $ = require('gulp-load-plugins')();
var browser = require('browser-sync');
var gulp = require('gulp');
var rimraf = require('rimraf');
var sequence = require('run-sequence');

gulp.task('clean', function(done) {
    //Delete our old css files
    rimraf('static/css/**/*', done);
});

// Compile SCSS files to CSS
gulp.task('sass', function() {
    return gulp.src('assets/sass/styles.scss')
        .pipe($.sourcemaps.init())
        .pipe($.sass())
        .on('error', $.sass.logError)
        .pipe($.sourcemaps.write())
        .pipe(gulp.dest('static/css'));
});

gulp.task('build', function(done) {
  sequence('clean', 'sass', done);
});

// Start a server with LiveReload to preview the site in
gulp.task('server', function() {
  browser.init({
    server: '.',
    port: 8000
  });
});

gulp.task('default', ['build', 'server'], function() {
  gulp.watch('assets/sass/**/*', ['sass', browser.reload])
});
