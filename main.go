package main

import (
    "flag"
    "github.com/dresswithpockets/go-vgui"
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
)

func main() {

    // parse arguments for our app state
    var root string
    flag.StringVar(&root, "root", "", "Specify root folder for hud.")
    flag.StringVar(&root, "r", "", "Specify root folder for hud.")
    flag.Parse()

    hudSourceProvider := &vgui.HudFileSourceProvider{Root: root}

    // the state of the application
    app := &App{
        &SourceVguiProvider{map[string]*vgui.Object{}, hudSourceProvider},
        &SourcePictureProvider{map[string]pixel.Picture{}, hudSourceProvider},

        RootControl{},

        nil,
        nil,
    }

    pixelgl.Run(app.run)
}
