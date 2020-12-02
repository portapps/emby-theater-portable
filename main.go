//go:generate go install -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
//go:generate goversioninfo -icon=res/papp.ico -manifest=res/papp.manifest
package main

import (
	"fmt"
	"os"

	"github.com/portapps/portapps/v3"
	"github.com/portapps/portapps/v3/pkg/log"
	"github.com/portapps/portapps/v3/pkg/utl"
	"github.com/portapps/portapps/v3/pkg/win"
)

var (
	app *portapps.App
)

func init() {
	var err error

	// Init app
	if app, err = portapps.New("emby-theater-portable", "Emby Theater"); err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize application. See log file for more info.")
	}
}

func main() {
	utl.CreateFolder(app.DataPath)
	app.Process = utl.PathJoin(app.AppPath, "Emby.Theater.exe")

	// Check arch
	if win.Is64Arch() && utl.Exists(utl.PathJoin(app.AppPath, "x86")) {
		log.Error().Msg("Emby Theater win32 cannot be launched on win64 system")
		if _, err := win.MsgBox(
			fmt.Sprintf("%s portable", app.Name),
			"Emby Theater win32 cannot be launched on win64 system.",
			win.MsgBoxBtnOk|win.MsgBoxIconError); err != nil {
			log.Error().Err(err).Msg("Cannot create dialog box")
		}
		return
	}

	// Data
	configFile := utl.PathJoin(app.AppPath, "Emby.Theater.exe.config")
	if err := utl.Replace(configFile, `ProgramDataPath" value=""`, `ProgramDataPath" value="../data"`); err != nil {
		log.Fatal().Err(err).Msg("Cannot change ProgramDataPath")
	}

	// Cancel cec driver install
	cecCancelFile, err := os.Create(utl.PathJoin(utl.CreateFolder(app.DataPath, "cec-driver"), "cancel"))
	if err != nil {
		log.Error().Err(err).Msg("Cannot write cec-driver cancel file")
	}
	cecCancelFile.Close()

	// Disable auto update
	systemFile := utl.PathJoin(app.DataPath, "config", "system.xml")
	if utl.Exists(systemFile) {
		if err := utl.Replace(systemFile, `<EnableAutoUpdate>true`, `<EnableAutoUpdate>false`); err != nil {
			log.Error().Err(err).Msg("Cannot disable auto update")
		}
	}

	defer app.Close()
	app.Launch(os.Args[1:])
}
