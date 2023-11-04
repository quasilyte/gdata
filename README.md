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
