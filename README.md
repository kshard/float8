<p align="center">
  <img src="./doc/hnsw.svg" height="240" />
  <h3 align="center">Float8</h3>
  <p align="center"><strong>minifloat in Golang</strong></p>

  <p align="center">
    <!-- Version -->
    <a href="https://github.com/kshard/float8/releases">
      <img src="https://img.shields.io/github/v/tag/kshard/float8?label=version" />
    </a>
    <!-- Documentation -->
    <a href="https://pkg.go.dev/github.com/kshard/float8">
      <img src="https://pkg.go.dev/badge/github.com/kshard/float8" />
    </a>
    <!-- Build Status -->
    <a href="https://github.com/kshard/float8/actions/">
      <img src="https://github.com/kshard/float8/workflows/build/badge.svg" />
    </a>
    <!-- GitHub -->
    <a href="http://github.com/kshard/float8">
      <img src="https://img.shields.io/github/last-commit/kshard/float8.svg" />
    </a>
    <!-- Coverage -->
    <a href="https://coveralls.io/github/kshard/float8?branch=main">
      <img src="https://coveralls.io/repos/github/kshard/float8/badge.svg?branch=main" />
    </a>
    <!-- Go Card -->
    <a href="https://goreportcard.com/report/github.com/kshard/float8">
      <img src="https://goreportcard.com/badge/github.com/kshard/float8" />
    </a>
  </p>
</p>

--- 

In computing, [minifloats](https://en.wikipedia.org/wiki/Minifloat) are floating-point values represented with very few bits. The library implements `float8` (8-bit `uint8`). It ideal for applications where memory and storage efficiency are crucial but lossy precision is acceptable (e.g. computer graphics, manche learning, etc).

## Features

- IEEE 754 and FP8 E4M3 compatible format.
- Fast conversion from/to float32.
- Fast algebraic operations (+, -, *, /).

## Getting Started

The latest version of the module is available at `main` branch. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

Use `go get` to retrieve the library and add it as dependency to your application.

```bash
go get -u github.com/kshard/float8
```

### Quick example

```go
package main

import (
  "fmt"
  "math"

  "github.com/kshard/float8"
)

func main() {
  phi08 := float8.ToFloat8(math.Phi)
  phi32 := float8.ToFloat32(phi08)
  fmt.Printf("𝝅(08) %f : %08b \n", phi32, phi08)
  fmt.Printf("𝝅(32) %f : %032b \n", float32(math.Phi), math.Float32bits(math.Phi))
}
```

The conversion is lossy and supported range of values is limited.


### Benchmark

The library implement public api using code books so that its pure Go code is efficient:

```
go test -run=^$ -bench=. -benchtime=1s -cpu 1
BenchmarkToFloat8     546344320         1.9410 ns/op
BenchmarkToFloat32   1000000000         0.5903 ns/op
BenchmarkAdd         1000000000         0.6267 ns/op
BenchmarkMul         1000000000         0.6284 ns/op
BenchmarkAddFloat32  1000000000         0.2874 ns/op
BenchmarkMulFloat32  1000000000         0.2872 ns/op
```

The internal package `math8` implements float-point algebra with focus on correctness, which is used to build code books.


## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

The build and testing process requires [Go](https://golang.org) version 1.21 or later.


### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/kshard/float8/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/kshard/float8.svg?style=for-the-badge)](LICENSE)


## References

1. https://en.wikipedia.org/wiki/Minifloat
2. https://en.wikipedia.org/wiki/Exponent_bias
3. https://en.wikipedia.org/wiki/Floating-point_arithmetic
4. https://docs.nvidia.com/deeplearning/transformer-engine/user-guide/examples/fp8_primer.html
5. https://github.com/x448/float16