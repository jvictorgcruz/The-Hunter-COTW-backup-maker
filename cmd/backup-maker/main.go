package main

import (
	"backup-maker/internal/backup"
	"backup-maker/internal/ui"
	"flag"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

func main() {
	autostartMode := flag.Bool("autostart", false, "Executa o backup em background e fecha")
	flag.Parse()

	app := app.NewWithID("com.jvictorgcruz.cotw-backup-maker")
	window := app.NewWindow("The Hunter: Call of the Wild Backup Maker")
	ui.SetupUI(window)
	ui.SetupSystray(app, window)
	window.Resize(fyne.NewSize(400, 200))

	if *autostartMode {
		if err := backup.BackupOnStartup(); err != nil {
			window.Show()
			dialog.ShowError(err, window)
		}
		app.Run()
	} else {
		window.ShowAndRun()
	}
}
