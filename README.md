# gdata

A gamedata package that provides a convenient cross-platform for games.

This package was made with [Ebitengine](https://github.com/hajimehoshi/ebiten/) in mind, but it should be usable with any kind of a game engine for Go.

Platforms supported:

* Windows
* Linux
* MacOS
* Browser/wasm (local storage)
* Android

This library tries to use the most conventional app data folder for every platform.

It provides a simple key-value style API. It can be considered to be a platform-agnostic localStorage.

## Installation

```bash
go get github.com/quasilyte/gdata
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/quasilyte/gdata"
)

func main() {
	// m is a data manager; treat it as a connection to a filesystem.
	m, err := gdata.Open(gdata.Config{
		AppName: "my_game",
	})
	if err != nil {
		panic(err)
	}

	if err := m.SaveItem("save.data", []byte("mydata")); err != nil {
		panic(err)
	}

	result, err := m.LoadItem("save.data")
	if err != nil {
		panic(err)
	}
	fmt.Println("=>", string(result)) // "mydata"

	fmt.Println("exists?", m.ItemExists("save.data")) // true
}
```
