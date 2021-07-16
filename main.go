package main

import (
    "flag"
    "fmt"
    "github.com/dresswithpockets/go-vgui"
    "github.com/faiface/pixel"
    "github.com/faiface/pixel/pixelgl"
    "log"
    "os"
)

var (
    warningLogger *log.Logger
    infoLogger    *log.Logger
    errorLogger   *log.Logger
)

func getPrintUsage(flagSet *flag.FlagSet) func() {
    return func() {
        fmt.Printf("Usage: %s [options...] <hudpath>\n", os.Args[0])
        flagSet.PrintDefaults()
    }
}

func main() {

    launchSettings := LaunchSettings{}

    // parse arguments for our app state
    flagSet := flag.NewFlagSet("", flag.ExitOnError)
    flagSet.Usage = getPrintUsage(flagSet)

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
        &launchSettings,

        &SourceVguiProvider{map[string]*vgui.Object{}, hudSourceProvider},
        &SourcePictureProvider{map[string]pixel.Picture{}, hudSourceProvider},

        &ControlProvider{},
        &RootControl{},

        nil,
        nil,
    }

    file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
    errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

    pixelgl.Run(app.run)
}
