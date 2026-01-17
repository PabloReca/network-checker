package main

import (
	"embed"
	"network-checker/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := &backend.App{}

	// Create application menu
	appMenu := menu.NewMenu()

	// File menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Open Configuration...", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		// Open config dialog
		selectedFile, err := app.OpenConfig()
		if err != nil {
			println("Error opening config:", err.Error())
			return
		}
		if selectedFile != "" {
			// Emit event to frontend to reload config
			runtime.EventsEmit(app.GetContext(), "config-loaded")
		}
	})

	err := wails.Run(&options.App{
		Title:  "Network Checker",
		Width:  900,
		Height: 550,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:      appMenu,
		OnStartup: app.Startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar:   mac.TitleBarDefault(),
			Appearance: mac.AppearanceType("NSAppearanceNameDarkAqua"),
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
