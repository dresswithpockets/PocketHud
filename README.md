# PocketHud

A WYSIWYG editor for VGUI interfaces. This project primarily targets Team Fortress 2 hud creation.

While this tool is to be used as the primary tool for editing HUDs, it cannot do everything; you may need to make manual edits to the VGUI resource files yourself.

N.B. **this software only works in readonly mode**. There are no edit features as of July 13th, 2021. Edit features will be introduced when the render pipeline is stable.

## Features

- [ ] Opening HUD for viewing
- [ ] Opening HUD for editing

## Usage

Open up a HUD in readonly mode:
```shell
$ go run github.com/dresswithpockets/PocketHud -- --root "path/to/my/hud" --readonly
```

N.B. as noted in the foreword, this software only worked in readonly mode. Attempting to run PocketHud in editmode (or without the --readonly flag) will panic and close immediately.

## Building
Pull dependencies:

```sh
$ go get github.com/dresswithpockets/go-vgui \
$     github.com/faiface/pixel \
$     github.com/ahmetb/go-linq \
$     golang.org/x/image
```

Get:
```shell
$ go get github.com/dresswithpockets/PocketHud
```

Build:
```shell
$ go build github.com/dresswithpockets/PocketHud
```

Run:
```shell
$ go run github.com/dresswithpockets/PocketHud
```

## Downloading/Installing

There are no pre-built releases for download at the moment. This software is still WIP and is not ready for general use.
