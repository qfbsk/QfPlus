package main

import (
	"context"
	"embed"
	"log"

	qfapp "QfPlus/internal/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := qfapp.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "QfPlus",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		StartHidden:      true,
		OnStartup: func(ctx context.Context) {
			qfapp.Startup(app, ctx)
		},
		OnShutdown: func(ctx context.Context) {
			qfapp.Shutdown(app)
		},
		OnDomReady: func(ctx context.Context) {
			runtime.WindowShow(ctx)
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "com.qfplus.app",
			OnSecondInstanceLaunch: func(_ options.SecondInstanceData) {
				qfapp.ShowMainWindow(app)
			},
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Println("Error:", err.Error())
	}
}
