package main

import (
    "flag"
    "fmt"
    "github.com/dresswithpockets/go-vgui"
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "os"
)

func getPrintUsage(flagSet *FlagSet)
func printUsage() {
    fmt.Printf("Usage: %s [options...] <hudpath>\n", os.Args[0])
    flag.PrintDefaults()
}

func main() {

    launchSettings := LaunchSettings{}

    // parse arguments for our app state
    flagSet := flag.NewFlagSet("", flag.ExitOnError)
    flagSet.Usage = printUsage

    flagSet.BoolVar(&launchSettings.readonly, "readonly", false, "Launch editor in readonly mode. Huds cannot be altered or saved in this mode, only viewed.")
    flagSet.BoolVar(&launchSettings.verbose, "verbose", false, "Detailed logging.")
    _ = flagSet.Parse(os.Args[1:])

    if flagSet.NArg() == 0 {
        fmt.Fprintln(os.Stderr, "Missing positional argument 'hudpath'")
        flagSet.Usage()
        os.Exit(1)
    }

    launchSettings.path = flagSet.Arg(0)

    if !launchSettings.readonly {
        // TODO: give link to issue with more information
        fmt.Println("Edit mode is not supported yet. Launch with -readonly flag.")
        return
    }

    hudSourceProvider := &vgui.HudFileSourceProvider{Root: launchSettings.path}

    // the state of the application
    app := &App{
        launchSettings,

        &SourceVguiProvider{map[string]*vgui.Object{}, hudSourceProvider},
        &SourcePictureProvider{map[string]pixel.Picture{}, hudSourceProvider},

        RootControl{},

        nil,
        nil,
    }

    pixelgl.Run(app.run)
}
