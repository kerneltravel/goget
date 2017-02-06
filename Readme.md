[![Build status][travis-img]][travis-url]
[![Test coverage][coveralls-img]][coveralls-url]
[![License][license-img]][license-url]
[![GoDoc][doc-img]][doc-url]

### goget

* `wget` by golang
  - auto decompression
  - process visualization
  - resume download from break point

### usage

```
go get github.com/coderhaoxin/goget

goget http://example.com/download.mp4
(ctrl+c to break download)

#resume download
goget http://example.com/download.mp4 -c
```

### todo:
* multi-thread download (e.g.  /usr/bin/axel -n 60 ...)
* more-vorbose printing info.
* parameter before <uri>

### License
MIT

[travis-img]: https://img.shields.io/travis/coderhaoxin/goget.svg?style=flat-square
[travis-url]: https://travis-ci.org/coderhaoxin/goget
[coveralls-img]: https://img.shields.io/coveralls/coderhaoxin/goget.svg?style=flat-square
[coveralls-url]: https://coveralls.io/r/coderhaoxin/goget?branch=master
[license-img]: http://img.shields.io/badge/license-MIT-green.svg?style=flat-square
[license-url]: http://opensource.org/licenses/MIT
[doc-img]: http://img.shields.io/badge/GoDoc-reference-blue.svg?style=flat-square
[doc-url]: http://godoc.org/github.com/coderhaoxin/goget
