package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "Shodan GUI",
		Width:            1400,
		Height:           900,
		MinWidth:         1100,
		MinHeight:        700,
		BackgroundColour: &options.RGBA{R: 13, G: 17, B: 23, A: 1},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		Frameless:        false,
		CSSDragProperty:  "--wails-draggable",
		CSSDragValue:     "drag",
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
