# COTW Save Backup Maker

The **COTW Save Backup Maker** is a lightweight, modern, and cross-platform desktop application written in **Go** using the **Fyne** GUI library. It is designed to manage and automate backups for the game *The Hunter: Call of the Wild*, protecting your progress against save file corruption, file loss, or accidental new game initialization.

---

## Features

- **Automatic Detection**: Automatically detects your game save directory on **Windows** and **Linux** (scanning standard paths, Steam, Epic Games Store, and OneDrive folders).
- **Multiple Backup Destinations**: Configurable support for both **Local storage** (saving zip backups to a local folder of your choice) and **Microsoft OneDrive** cloud backup. You can enable either or both simultaneously!
- **OneDrive Cloud Integration**: Fully native cloud backups. Login securely in your browser via OAuth2 (using a built-in local server) and automatically store your backups under a dedicated folder (`/Apps/COTW Backup Maker`).
- **Graphical Interface**: Clean, modern UI offering interactive fields for directory selection, provider checkboxes, and visual toggles.
- **System Tray (Systray)**: The app minimizes to the system tray, running silently in the background without interrupting your gameplay.
- **Start with the PC**: Registers permanently in the operating system's startup (Linux `.desktop` and Windows Registry) to remain active in the background.
- **Smart Backup**: Recursive `.zip` compression with chronological naming and safe stream handling.
- **Automatic Cleanup (Rotation)**: Automatically keeps only the **`Defined amount` most recent backups** in your selected storage destinations (both local folder and OneDrive), deleting older ones to save space.
- **Backup on Boot**: Option to trigger a silent, automatic backup when your computer starts, hiding in the tray immediately after.


---

## How to Run from Source (Development)

### Prerequisites
1. **Go installed** (version 1.22 or superior).
2. **C Compiler (GCC/CGO)** installed on the system (required for Fyne's OpenGL bindings to compile graphical interface resources).
   - *On Ubuntu/Debian*: `sudo apt install build-essential libgl1-mesa-dev libegl1-mesa-dev libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev libxxf86vm-dev`

### Execution steps:
1. Clone the repository:
   ```bash
   git clone https://github.com/jvictorgcruz/The-Hunter-COTW-backup-maker.git
   cd The-Hunter-COTW-backup-maker
   ```
2. Download module dependencies:
   ```bash
   go mod tidy
   ```
3. Run the application directly:
   ```bash
   go run ./cmd/backup-maker
   ```

---

## How to Compile and Package (Release)

To generate standalone packages with the official icon embedded and without terminal windows associated:

### 1. Install Fyne CLI
```bash
go install fyne.io/tools/cmd/fyne@latest
```

### 2. Generate Release for Linux
1. Enter the main package directory:
   ```bash
   cd cmd/backup-maker
   ```
2. Generate the package (`.tar.xz`):
   ```bash
   ~/go/bin/fyne package -os linux -icon ../../assets/icon.png
   ```

### 3. Generate Release for Windows
Since Fyne requires CGO, to build the Windows `.exe` from Linux you will need **Docker** and the **`fyne-cross`** tool:

1. Install `fyne-cross`:
   ```bash
   go install github.com/fyne-io/fyne-cross@latest
   ```
2. From the root directory, run the cross-compilation command:
   ```bash
   ~/go/bin/fyne-cross windows -icon assets/icon.png cmd/backup-maker/main.go
   ```
   *The executable `backup-maker.exe` will be generated inside the `fyne-cross/dist/windows-amd64/` directory.*

If you are running **natively on a Windows machine**, you can automate the compilation and create a setup installer using the provided `build.bat` script:

1. Install [Inno Setup](https://jrsoftware.org/isinfo.php) and ensure `iscc` (the compiler executable) is in your system `PATH`.
2. Run the `build.bat` script from the project root:
   ```cmd
   .\build.bat
   ```
   *This packages the app and compiles the setup installer `COTWBackupMakerSetup.exe` inside the `Output/` directory.*

To package the application and compile the installer manually:
1. Enter the main package directory and run `fyne package`:
   ```cmd
   cd cmd/backup-maker
   fyne package -os windows -icon ../../assets/icon.png
   ```
2. In the root directory, compile the `installer.iss` script:
   - Either open `installer.iss` in Inno Setup Compiler and compile it (`Ctrl + F9`),
   - Or run from the terminal:
     ```cmd
     iscc installer.iss
     ```

---

## License

This project is open-source under the terms of the repository license.
