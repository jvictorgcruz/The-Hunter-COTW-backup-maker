package ui

import (
	_ "embed"
	"errors"
	"strconv"

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

	backupOnStartup := cfg.BackupOnStartup
	if cfg.MaxBackups <= 0 {
		cfg.MaxBackups = 3
	}
	maxBackups := cfg.MaxBackups

	checkIfChanged := func() {
		if saveBtn == nil {
			return
		}
		hasChanges := sourceEntry.Text != cfg.SourceDir ||
			destinationEntry.Text != cfg.DestinationDir ||
			backupOnStartup != cfg.BackupOnStartup ||
			maxBackups != cfg.MaxBackups
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
		backupOnStartup = checked
		checkIfChanged()
	})
	startupCheck.SetChecked(backupOnStartup)

	maxBackupsSelect := widget.NewSelect(
		[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
		func(selected string) {
			if val, err := strconv.Atoi(selected); err == nil {
				maxBackups = val
				checkIfChanged()
			}
		},
	)
	maxBackupsSelect.SetSelected(strconv.Itoa(maxBackups))

	backupBtn := widget.NewButton("Fazer Backup Agora", func() {
		sourcePath := sourceEntry.Text
		destinationPath := destinationEntry.Text

		if sourcePath == "" || destinationPath == "" {
			dialog.ShowError(errors.New("Por favor, selecione as pastas de origem e destino."), window)
			return
		}

		go func() {
			err := backup.CreateBackup(sourceEntry.Text, destinationEntry.Text, maxBackups)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			dialog.ShowInformation("Aviso", "Backup realizado com sucesso!", window)
		}()
	})

	resetBtn := widget.NewButton("Limpar Configurações", func() {
		dialog.ShowConfirm(
			"Confirmar Ação",
			"Deseja realmente apagar todas as configurações salvas?",
			func(confirmed bool) {
				if !confirmed {
					return
				}

				err := config.ClearConfig()
				if err != nil {
					dialog.ShowError(err, window)
					return
				}
				sourceEntry.SetText("")
				destinationEntry.SetText("")
				startupCheck.SetChecked(false)
				maxBackupsSelect.SetSelected("3")
				dialog.ShowInformation("Sucesso", "Configurações limpas com sucesso!", window)
			},
			window,
		)
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
		cfg.BackupOnStartup = backupOnStartup

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

	backupBtn.Importance = widget.HighImportance
	resetBtn.Importance = widget.DangerImportance
	saveBtn.Importance = widget.SuccessImportance

	content := container.NewVBox(
		widget.NewLabel("Pasta de Origem (Save do Jogo):"),
		sourceRow,
		widget.NewLabel("Pasta de Destino (Onde salvar o ZIP):"),
		destinationRow,
		widget.NewLabel("Limite Máximo de Backups:"),
		maxBackupsSelect,
		startupCheck,
		backupBtn,
		resetBtn,
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
