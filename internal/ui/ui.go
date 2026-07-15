package ui

import (
	"context"
	_ "embed"
	"errors"
	"os"
	"slices"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"backup-maker/internal/backup"
	"backup-maker/internal/config"
	"backup-maker/internal/onedrive"
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

	localEnabled := slices.Contains(cfg.EnabledProviders, "local")
	onedriveEnabled := slices.Contains(cfg.EnabledProviders, "onedrive")

	sourceEntry, _, sourceRow := createFolderSelector("Selecione a pasta de origem...", window)
	if cfg.SourceDir != "" {
		sourceEntry.SetText(cfg.SourceDir)
	} else if defaultPath := paths.DetectSavePath(); defaultPath != "" {
		sourceEntry.SetText(defaultPath)
	}

	destinationEntry, destinationBtn, destinationRow := createFolderSelector("Selecione a pasta de destino...", window)
	destinationEntry.SetText(cfg.DestinationDir)

	backupOnStartup := cfg.BackupOnStartup
	if cfg.MaxBackups <= 0 {
		cfg.MaxBackups = 3
	}
	maxBackups := cfg.MaxBackups

	var localCheck *widget.Check
	var onedriveCheck *widget.Check

	checkIfChanged := func() {
		if saveBtn == nil {
			return
		}

		var currentProviders []string
		if localEnabled {
			currentProviders = append(currentProviders, "local")
		}
		if onedriveEnabled {
			currentProviders = append(currentProviders, "onedrive")
		}

		if len(currentProviders) == 0 {
			saveBtn.Disable()
			return
		}

		providersChanged := !slices.Equal(currentProviders, cfg.EnabledProviders)

		hasChanges := sourceEntry.Text != cfg.SourceDir ||
			destinationEntry.Text != cfg.DestinationDir ||
			backupOnStartup != cfg.BackupOnStartup ||
			maxBackups != cfg.MaxBackups ||
			providersChanged

		if hasChanges {
			saveBtn.Enable()
		} else {
			saveBtn.Disable()
		}
	}

	toggleLocalVisibility := func(enabled bool) {
		if enabled {
			destinationEntry.Enable()
			destinationBtn.Enable()
		} else {
			destinationEntry.Disable()
			destinationBtn.Disable()
		}
	}

	sourceEntry.OnChanged = func(text string) {
		checkIfChanged()
	}

	destinationEntry.OnChanged = func(text string) {
		checkIfChanged()
	}

	localCheck = widget.NewCheck("Backup Local (Pasta de destino do .zip)", func(checked bool) {
		localEnabled = checked
		toggleLocalVisibility(checked)
		checkIfChanged()
	})
	localCheck.SetChecked(localEnabled)

	onedriveCheck = widget.NewCheck("Backup no OneDrive", func(checked bool) {
		onedriveEnabled = checked
		if checked {
			if !onedrive.IsAuthenticated() {
				go func() {
					_, err := onedrive.GetClient(context.Background())
					if err != nil {
						dialog.ShowError(err, window)
						onedriveCheck.SetChecked(false)
					} else {
						dialog.ShowInformation("OneDrive", "Autenticação realizada com sucesso!", window)
						checkIfChanged()
					}
				}()
			}
		}
		checkIfChanged()
	})
	onedriveCheck.SetChecked(onedriveEnabled)

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

		if sourcePath == "" {
			dialog.ShowError(errors.New("Por favor, selecione a pasta de origem."), window)
			return
		}

		if localEnabled && destinationPath == "" {
			dialog.ShowError(errors.New("Por favor, selecione a pasta de destino local."), window)
			return
		}

		if onedriveEnabled && !onedrive.IsAuthenticated() {
			dialog.ShowError(errors.New("O OneDrive está ativado, mas você não está autenticado. Por favor, marque e desmarque o OneDrive para realizar o login."), window)
			return
		}

		var activeProviders []string
		if localEnabled {
			activeProviders = append(activeProviders, "local")
		}
		if onedriveEnabled {
			activeProviders = append(activeProviders, "onedrive")
		}

		go func() {
			err := backup.CreateBackup(sourceEntry.Text, activeProviders, maxBackups)
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
				localCheck.SetChecked(true)
				onedriveCheck.SetChecked(false)
				dialog.ShowInformation("Sucesso", "Configurações limpas com sucesso!", window)
			},
			window,
		)
	})

	saveBtn = widget.NewButton("Salvar Configurações", func() {
		sourcePath := sourceEntry.Text
		destinationPath := destinationEntry.Text

		if sourcePath == "" {
			dialog.ShowError(errors.New("Por favor, selecione a pasta de origem."), window)
			return
		}

		if localEnabled && destinationPath == "" {
			dialog.ShowError(errors.New("Por favor, selecione a pasta de destino local."), window)
			return
		}

		cfg.SourceDir = sourcePath
		cfg.DestinationDir = destinationPath
		cfg.BackupOnStartup = backupOnStartup
		cfg.MaxBackups = maxBackups

		var activeProviders []string
		if localEnabled {
			activeProviders = append(activeProviders, "local")
		}
		if onedriveEnabled {
			activeProviders = append(activeProviders, "onedrive")
		}
		cfg.EnabledProviders = activeProviders

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
		widget.NewSeparator(),
		widget.NewLabel("Destinos de Backup:"),
		localCheck,
		destinationRow,
		onedriveCheck,
		widget.NewSeparator(),
		widget.NewLabel("Limite Máximo de Backups:"),
		maxBackupsSelect,
		startupCheck,
		widget.NewSeparator(),
		backupBtn,
		resetBtn,
		saveBtn,
	)

	saveBtn.Disable()
	window.SetContent(content)
	toggleLocalVisibility(localEnabled)
}

func createFolderSelector(placeholder string, window fyne.Window) (*widget.Entry, *widget.Button, fyne.CanvasObject) {
	entry := widget.NewEntry()
	entry.SetPlaceHolder(placeholder)
	btn := widget.NewButton("Buscar...", func() {
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			if err == nil && uri != nil {
				entry.SetText(uri.Path())
			}
		}, window)

		defer folderDialog.Show()

		currentPath := entry.Text
		if currentPath == "" {
			return
		}

		info, err := os.Stat(currentPath)
		if err != nil || !info.IsDir() {
			return
		}

		fileURI := storage.NewFileURI(currentPath)
		folderURI, err := storage.ListerForURI(fileURI)
		if err != nil {
			return
		}

		folderDialog.SetLocation(folderURI)
	})
	row := container.NewBorder(nil, nil, nil, btn, entry)
	return entry, btn, row
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
