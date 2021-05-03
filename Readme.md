[![](https://img.shields.io/github/license/masagroup/soft.go.svg)](https://github.com/masagroup/soft.go/blob/master/LICENSE)
![](https://github.com/masagroup/soft.go/actions/workflows/build_and_test.yaml/badge.svg)
[![Go](https://img.shields.io/github/go-mod/go-version/masagroup/soft.go)](https://github.com/masagroup/soft.go)
![](https://img.shields.io/github/v/release/masagroup/soft.go)
# Soft Go #

Soft Go is an implementation of the EMF Ecore library in Golang. This library is partially generated and referenced by the [Soft Go Generator](https://github.com/masagroup/soft.gen)
 
Soft Go is part of [Soft](https://github.com/masagroup/soft) project

# Installation #
To install Ecore, use go get:
```shell
go get github.com/masagroup/soft.go
```

This will then make the following packages available to you:
```
github.com/masagroup/soft.go/ecore
```

Import the masagroup/soft.go/ecore package into your code using this template:
```Golang
package yours

import (
  "github.com/masagroup/soft.go/ecore"
)
```

# Supported go versions #
We support Go v1.12

