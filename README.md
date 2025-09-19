# Backup Tool - CLI Project Backup Manager

A simple and powerful CLI tool for creating, managing and restoring project backups.

## Installation

### Automatic Installation (Recommended)

1. Download or clone the source code
2. Run the installer:
   - **PowerShell**: `.\install.ps1`
   - **Command Prompt**: `install.bat`
3. The installer will:
   - Check Go installation
   - Build the project automatically
   - Install to `%LOCALAPPDATA%\BackupTool`
   - Add to PATH for current user
   - Clean up build artifacts

### Manual Installation

1. Ensure Go 1.21+ is installed
2. Build the project: `go build -ldflags="-s -w" -o backup.exe`
3. Copy `backup.exe` to any directory in your PATH
4. Ready! Now you can use the `backup` command from any directory

### Uninstallation

Run the uninstaller: `.\uninstall.ps1`

## Commands

### `backup init`
Initialize current directory for backup management.

```bash
backup init
```

**What happens:**
- Creates `.backup-config.json` file with unique project ID
- Sets up backup directory in `%APPDATA%/ProjectBackup/{project-id}/`
- Configures default exclusions

### `backup create`
Create new project backup.

```bash
# Create backup with automatic name (date and time)
backup create

# Create backup with custom name
backup create --name "Before refactoring"
backup create -n "Version 1.0"
```

**Features:**
- Automatically excludes `node_modules/`, `.git/`, `build/`, `dist/` and other system folders
- Shows archiving progress bar
- Compresses files to ZIP format
- Supports projects up to 1GB

### `backup list`
Display interactive list of all backups.

```bash
backup list
```

**Capabilities:**
- View all backups with sizes and dates
- Display backup age (minutes, hours, days ago)
- Navigate with ↑/↓ arrows
- Quick actions:
  - `Enter` - load backup
  - `r` - rename backup
  - `d` - delete backup
  - `q` - quit

### `backup load`
Load backup into current directory.

```bash
# Interactive backup selection
backup load

# Load specific backup by name
backup load --name "Before refactoring"
backup load -n "Version 1.0"
```

**WARNING:** All files in current directory will be deleted! Operation requires confirmation.

## File Exclusions

By default excludes:
- `node_modules/` - Node.js dependencies
- `.git/` - Git repository
- `build/`, `dist/` - Build folders
- `*.tmp`, `*.log` - Temporary files
- `.env*` - Environment files
- `*.exe`, `*.dll` - Executable files
- `target/` - Rust/Java build
- `bin/`, `obj/` - .NET build
- `.vs/`, `.vscode/` - IDE files
- `__pycache__/`, `*.pyc` - Python cache

## Storage Structure

```
%APPDATA%/ProjectBackup/
├── {project-uuid-1}/
│   ├── backup_20240119_143022.zip
│   ├── backup_20240119_150315_MyFeature.zip
│   └── backup_20240120_091500_Release.zip
└── {project-uuid-2}/
    └── ...
```

## Example Usage

### Typical Workflow

1. **Initialize new project:**
```bash
cd my-project
backup init
```

2. **Create backup before important changes:**
```bash
backup create --name "Before API refactoring"
```

3. **View all backups:**
```bash
backup list
# Use arrows to navigate, Enter to load
```

4. **Restore to previous state:**
```bash
backup load --name "Before API refactoring"
```

## Technical Details

- **Language:** Go 1.21+
- **Archive format:** ZIP with basic compression
- **Supported OS:** Windows
- **Max project size:** 1GB
- **Storage:** Local in `%APPDATA%/ProjectBackup`

## Features

- Simple to use - intuitive interface
- Beautiful UI - colored output and progress bars
- Fast operation - efficient archiving
- Safe - confirmation for critical operations
- Compact - single executable without dependencies

---

**Version:** 1.0  
**Platform:** Windows