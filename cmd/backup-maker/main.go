package main

import (
	"backup-maker/internal/backup"
	"backup-maker/internal/instance"
	"backup-maker/internal/ui"
	"flag"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

func main() {
	autostartMode := flag.Bool("autostart", false, "Executa o backup em background e fecha")
	flag.Parse()

	isFirst, ch, err := instance.TryLock()
	if err != nil {
		return
	}
	if !isFirst {
		if !*autostartMode {
			instance.NotifyExisting()
		}
		os.Exit(0)
	}

	os.Setenv("FYNE_SYSTEM_DIALOGS", "1")

	app := app.NewWithID("com.jvictorgcruz.cotw-backup-maker")
	window := app.NewWindow("The Hunter: Call of the Wild Backup Maker")
	ui.SetupUI(window)
	ui.SetupSystray(app, window)
	window.Resize(fyne.NewSize(400, 200))

	go func() {
		for msg := range ch {
			if msg == "show" {
				window.Show()
				window.RequestFocus()
			}
		}
	}()

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
