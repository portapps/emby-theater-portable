//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"fmt"
	"os"

	. "github.com/portapps/portapps"
	"github.com/portapps/portapps/pkg/dialog"
	"github.com/portapps/portapps/pkg/utl"
	"github.com/portapps/portapps/pkg/win"
)

var (
	app *App
)

func init() {
	var err error

	// Init app
	if app, err = New("emby-theater-portable", "Emby Theater"); err != nil {
		Log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "Emby.Theater.exe")

	// Check arch
	if win.Is64Arch() && utl.Exists(utl.PathJoin(app.AppPath, "x86")) {
		Log.Error().Msg("Emby Theater win32 cannot be launched on win64 system")
		if _, err := dialog.MsgBox(
			fmt.Sprintf("%s portable", app.Name),
			"Emby Theater win32 cannot be launched on win64 system.",
			dialog.MsgBoxBtnOk|dialog.MsgBoxIconError); err != nil {
			Log.Error().Err(err).Msg("Cannot create dialog box")
		}
		return
	}

	// Data
	configFile := utl.PathJoin(app.AppPath, "Emby.Theater.exe.config")
	if err := utl.Replace(configFile, `ProgramDataPath" value=""`, `ProgramDataPath" value="../data"`); err != nil {
		Log.Fatal().Err(err).Msg("Cannot change ProgramDataPath")
	}

	// Cancel cec driver install
	cecCancelFile, err := os.Create(utl.PathJoin(utl.CreateFolder(app.DataPath, "cec-driver"), "cancel"))
	if err != nil {
		Log.Error().Err(err).Msg("Cannot write cec-driver cancel file")
	}
	cecCancelFile.Close()

	// Disable auto update
	systemFile := utl.PathJoin(app.DataPath, "config", "system.xml")
	if utl.Exists(systemFile) {
		if err := utl.Replace(systemFile, `<EnableAutoUpdate>true`, `<EnableAutoUpdate>false`); err != nil {
			Log.Error().Err(err).Msg("Cannot disable auto update")
		}
	}

	app.Launch(os.Args[1:])
}
