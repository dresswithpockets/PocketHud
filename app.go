package main

import (
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "golang.org/x/image/colornames"
)

type LaunchSettings struct {
    path     string
    readonly bool
    verbose  bool
}

type App struct {
    launchSettings LaunchSettings

    vguiProvider    VguiProvider
    pictureProvider PictureProvider

    rootControl RootControl

    window *pixelgl.Window
    batch  *pixel.Batch
}

type RootControl struct {
    BaseControl
}

func (r *RootControl) draw() {}

func (a *App) run() {
    // initialize graphics context and window
    cfg := pixelgl.WindowConfig{
        Title:     "PocketHud: VGUI Hud Editor",
        Bounds:    pixel.R(0, 0, 1024, 768),
        Resizable: true,
        VSync:     true,
    }

    win, err := pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    a.window = win

    for !win.Closed() {
        a.draw()
        win.Update()
    }
}

func (a *App) draw() {
    a.window.Clear(colornames.Skyblue)
    a.rootControl.drawChildren()
}
