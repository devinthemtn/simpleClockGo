# simpleclock

A simple cross-platform desktop clock built with Go and [Fyne](https://fyne.io/). Displays the current time and date, updating every second.

## Features

- Large, readable time display
- Full date shown below the time
- Optional secondary clocks for additional timezones
- Day offset indicator when a secondary clock is on a different calendar day
- Respects system light/dark theme
- Borderless mode with draggable window
- Cross-platform: Linux, Windows, macOS

## Requirements

### Go

Go 1.22 or later.

### Linux

```bash
sudo apt-get install -y libgl1-mesa-dev libx11-dev libxcursor-dev \
  libxrandr-dev libxinerama-dev libxi-dev libxxf86vm-dev
```

### Windows (cross-compiling from Linux)

```bash
sudo apt-get install -y gcc-mingw-w64-x86-64
```

### macOS (cross-compiling from Linux)

Requires [osxcross](https://github.com/tpoechtrager/osxcross) with `o64-clang` on your `PATH`. Building natively on a Mac requires no extra tools beyond Go.

## Building

```bash
# Linux
make

# Windows
make windows

# macOS
make mac

# Clean build artifacts
make clean
```

Binaries are output to the `bin/` directory.

## Running

```bash
# After building
./bin/simpleclock

# Or directly with Go
go run .
```

## Options

| Flag | Description |
|------|-------------|
| `--no-titlebar` | Launch without a window title bar. The window can still be dragged by clicking and dragging anywhere on it. |

```bash
./bin/simpleclock --no-titlebar
```

## Configuration

Secondary clocks are configured via a YAML file. The file is optional — the app runs fine without it.

**Config file location:**

| Platform | Path |
|----------|------|
| Linux / macOS (XDG) | `~/.config/simpleclock/config.yaml` |
| macOS (standard) | `~/Library/Application Support/simpleclock/config.yaml` |
| Windows | `%APPDATA%\simpleclock\config.yaml` (e.g. `C:\Users\<user>\AppData\Roaming\simpleclock\config.yaml`) |

Create the `simpleclock` directory if it does not exist, then add a `config.yaml`:

```yaml
clocks:
  - timezone: "America/New_York"
    label: "New York"
  - timezone: "Europe/London"
    label: "London"
  - timezone: "Asia/Tokyo"
    label: "Tokyo"
```

- `timezone` — IANA timezone name (e.g. `America/Chicago`, `UTC`). Required.
- `label` — display name shown next to the clock. Defaults to the timezone name if omitted.

Changes to the config file take effect on the next app launch.

## Dependencies

- [fyne.io/fyne/v2](https://github.com/fyne-io/fyne) v2.7.3

## Author

**Thomas Koefod** — [devinthemtn@gmail.com](mailto:devinthemtn@gmail.com)
