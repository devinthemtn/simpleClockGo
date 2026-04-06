# simpleclock

A simple cross-platform desktop clock built with Go and [Fyne](https://fyne.io/). Displays the current time and date, updating every second.

## Features

- Large, readable time display
- Full date shown below the time
- Respects system light/dark theme
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

## Dependencies

- [fyne.io/fyne/v2](https://github.com/fyne-io/fyne) v2.7.3

## Author

**Thomas Koefod** — [devinthemtn@gmail.com](mailto:devinthemtn@gmail.com)
