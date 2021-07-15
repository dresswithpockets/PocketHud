# PocketHud

A WYSIWYG editor for VGUI interfaces. This project primarily targets Team Fortress 2 hud creation.

While this tool is to be used as the primary tool for editing HUDs, it cannot do everything; you may need to make manual edits to the VGUI resource files yourself.

N.B. **this software only works in readonly mode**. There are no edit features as of July 13th, 2021. Edit features will be introduced when the render pipeline is stable.

## Usage

Open up a HUD in readonly mode:
```shell
$ go run github.com/dresswithpockets/PocketHud -- --root "path/to/my/hud" --readonly
```

N.B. as noted in the foreword, this software only worked in readonly mode. Attempting to run PocketHud in editmode (or without the --readonly flag) will panic and close immediately.

N.B.B. **this software is not production ready and is subject to sweeping changes at any moment.** Versions are likely to be incompatible between each other. Compatibility and long term support will be the end goal after v1.0.

## Features

- [ ] Opening HUD for readonly viewing - https://github.com/dresswithpockets/PocketHud/issues/1
    - [ ] VGUI Backend
        - [ ] Surface (Target Painting & Layout)
        - [ ] Panel
        - [ ] Drawing Targets (Image, TextImage)
    - [ ] VGUI Base Controls
        - [ ] Label
        - [ ] URLLabel
        - [ ] Button
        - [ ] ImageButton
        - [ ] ImagePanel
        - [ ] EditablePanel
        - [ ] and more...
    - [ ] TF2 Extended Controls
        - [ ] CExLabel
        - [ ] CExButton
        - [ ] CExImageButton
        - [ ] CExRichText
        - [ ] CExplanationPopup
        - [ ] and more...
    - [ ] Schemes
    - [ ] Menus
        - [ ] Menu/view selection
        - [ ] Scene backgrounds
    - [ ] Reactive layout
        - [ ] Aspect-ratio selection
    - [ ] Simulation
- [ ] Opening HUD for editing - https://github.com/dresswithpockets/PocketHud/issues/2

N.B. the roadmap has not been fleshed out, so this section is lacking for now.

## Building
Built with Go 1.14.2. Not guaranteed to work on any other versions.

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

## Supporting Development

This project takes a substantial amount of time to maintain; if you find this tool useful and want to support its development, please consider contributing or tipping.

The best way to support the development of this project is to contribute to the codebase.

Otherwise, shoot me a tip at [ko-fi.com/dresswithpockets](https://ko-fi.com/dresswithpockets). Anything at all really helps!