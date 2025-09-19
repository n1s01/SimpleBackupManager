# Backup Tool - CLI Project Backup Manager

A simple and powerful CLI tool for creating, managing and restoring project backups.

## Installation

### One-Line Install (Fastest)

```powershell
iex (irm https://raw.githubusercontent.com/n1s01/SimpleBackupManager/v1.0.0/install.ps1)
```

### Manual Install

1. **Download installer**: [install.ps1](https://github.com/n1s01/SimpleBackupManager/blob/v1.0.0/install.ps1)
2. **Run installer**: `.\install.ps1`
3. **Restart terminal** and use `backup` commands

### Direct Download

1. Download `backup.exe` from [latest release](https://github.com/n1s01/SimpleBackupManager/releases/latest)
2. Copy to any directory in your PATH

### Uninstallation

```powershell
iex (irm https://raw.githubusercontent.com/n1s01/SimpleBackupManager/v1.0.0/uninstall.ps1)
```

Or download and run: [uninstall.ps1](https://github.com/n1s01/SimpleBackupManager/blob/v1.0.0/uninstall.ps1)

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