# gdata

![Build Status](https://github.com/quasilyte/gdata/workflows/Go/badge.svg)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/quasilyte/gdata)](https://pkg.go.dev/mod/github.com/quasilyte/gdata)

A gamedata package that provides a convenient cross-platform storage for games.

Some examples of such gamedata that you might want to store:

* Game settings
* Save states
* Replays
* Pluging/mods metadata

This package was made with [Ebitengine](https://github.com/hajimehoshi/ebiten/) in mind, but it should be usable with any kind of a game engine for Go.

Platforms supported:

* Windows (file system, AppData)
* Linux (file system, ~/.local/share)
* MacOS (file system, ~/.local/share)
* Android (file system, app data directory)
* Browser/wasm (local storage)

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
