package ui

import (
	_ "embed"
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"backup-maker/internal/backup"
	"backup-maker/internal/config"
	"backup-maker/internal/paths"
	"backup-maker/internal/startup"
)

//go:embed icon.png
var iconBytes []byte

var appIcon = fyne.NewStaticResource("icon.png", iconBytes)

func SetupUI(window fyne.Window) {
	var saveBtn *widget.Button
	cfg, err := config.LoadConfig()
	if err != nil {
		dialog.ShowError(err, window)
	}

	sourceEntry, sourceRow := createFolderSelector("Selecione a pasta de origem...", window)
	if cfg.SourceDir != "" {
		sourceEntry.SetText(cfg.SourceDir)
	} else if defaultPath := paths.DetectSavePath(); defaultPath != "" {
		sourceEntry.SetText(defaultPath)
	}

	destinationEntry, destinationRow := createFolderSelector("Selecione a pasta de destino...", window)
	destinationEntry.SetText(cfg.DestinationDir)

	BackupOnStartup := cfg.BackupOnStartup

	checkIfChanged := func() {
		if saveBtn == nil {
			return
		}
		hasChanges := sourceEntry.Text != cfg.SourceDir ||
			destinationEntry.Text != cfg.DestinationDir ||
			BackupOnStartup != cfg.BackupOnStartup
		if hasChanges {
			saveBtn.Enable()
		} else {
			saveBtn.Disable()
		}
	}

	sourceEntry.OnChanged = func(text string) {
		checkIfChanged()
	}

	destinationEntry.OnChanged = func(text string) {
		checkIfChanged()
	}

	startupCheck := widget.NewCheck("Fazer backup ao iniciar o computador", func(checked bool) {
		BackupOnStartup = checked
		checkIfChanged()
	})
	startupCheck.SetChecked(BackupOnStartup)

	backupBtn := widget.NewButton("Fazer Backup Agora", func() {
		sourcePath := sourceEntry.Text
		destinationPath := destinationEntry.Text

		if sourcePath == "" || destinationPath == "" {
			dialog.ShowError(errors.New("Por favor, selecione as pastas de origem e destino."), window)
			return
		}

		go func() {
			if err := backup.CreateBackup(sourceEntry.Text, destinationEntry.Text); err != nil {
				dialog.ShowError(err, window)
				return
			}

			dialog.ShowInformation("Aviso", "Backup realizado com sucesso!", window)
		}()
	})

	saveBtn = widget.NewButton("Salvar Configurações", func() {
		sourcePath := sourceEntry.Text
		destinationPath := destinationEntry.Text

		if sourcePath == "" || destinationPath == "" {
			dialog.ShowError(errors.New("Por favor, selecione as pastas de origem e destino."), window)
			return
		}

		cfg.SourceDir = sourcePath
		cfg.DestinationDir = destinationPath
		cfg.BackupOnStartup = BackupOnStartup

		if err := config.SaveConfig(cfg); err != nil {
			dialog.ShowError(err, window)
			return
		}

		if err = startup.SetAutostart(); err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation("Sucesso", "Configurações salvas com sucesso!", window)
		saveBtn.Disable()
	})

	content := container.NewVBox(
		widget.NewLabel("Pasta de Origem (Save do Jogo):"),
		sourceRow,
		widget.NewLabel("Pasta de Destino (Onde salvar o ZIP):"),
		destinationRow,
		startupCheck,
		backupBtn,
		saveBtn,
	)

	saveBtn.Disable()
	window.SetContent(content)
}

func createFolderSelector(placeholder string, window fyne.Window) (*widget.Entry, fyne.CanvasObject) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	btn := widget.NewButton("Buscar...", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				entry.SetText(uri.Path())
			}
		}, window)
	})
	row := container.NewBorder(nil, nil, nil, btn, entry)
	return entry, row
}

func SetupSystray(app fyne.App, window fyne.Window) {
	window.SetIcon(appIcon)

	window.SetCloseIntercept(func() {
		window.Hide()
	})

	if desk, ok := app.(desktop.App); ok {
		desk.SetSystemTrayIcon(appIcon)

		menu := fyne.NewMenu("COTW Backup Maker",
			fyne.NewMenuItem("Configurações", func() {
				window.Show()
			}),
			fyne.NewMenuItem("Sair", func() {
				app.Quit()
			}),
		)
		desk.SetSystemTrayMenu(menu)
	}
}
