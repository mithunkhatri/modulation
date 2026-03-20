# Modulation 🎵
![Build Status](https://github.com/mithunkhatri/modulation/actions/workflows/build.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/mithunkhatri/modulation)](https://goreportcard.com/report/github.com/mithunkhatri/modulation)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mithunkhatri/modulation)](https://github.com/mithunkhatri/modulation)
[![Latest Release](https://img.shields.io/github/v/release/mithunkhatri/modulation)](https://github.com/mithunkhatri/modulation/releases)
[![License](https://img.shields.io/github/license/mithunkhatri/modulation)](https://github.com/mithunkhatri/modulation/blob/main/LICENSE)
[![Total Downloads](https://img.shields.io/github/downloads/mithunkhatri/modulation/total)](https://github.com/mithunkhatri/modulation/releases)
![Platforms](https://img.shields.io/badge/platform-macOS%20%7C%20Linux%20%7C%20Windows-blue)

> Resonating through your favorite terminal

Modulation is a high-performance, cross-platform TUI (Terminal User Interface) radio player inspired by `htop` and `pianobar`. It allows you to stream thousands of radio stations from around the world directly from your command line.

## Features ✨

- **HTOP-like Interface**: Modern, responsive TUI with a dedicated header, searchable station list, and footer shortcuts.
- **Global Radio Library**: Dynamically fetches the most popular stations from the Radio Browser API.
- **Robust Audio Engine**: Powered by FFmpeg for high-quality decoding of any audio format (MP3, AAC, FLAC, etc.).
- **Smart Management**:
  - **Search**: Press `/` to search for stations by name.
  - **Filter**: Press `f` to filter by category or genre (e.g., "jazz", "rock").
  - **Favorites**: Mark your favorite stations with `a` and view them with `v`.
- **Seamless Playback**:
  - **Minimized Mode**: Press `b` to keep the music playing in the background while hiding the TUI.
  - **Buffering Animations**: Visual feedback so you know when a stream is connecting.

## Installation 🚀

The easiest way to install Modulation on Mac, Linux, or Windows (via Git Bash/WSL) is with a single command. This script automatically detects your OS, installs any missing dependencies (Go/FFmpeg), and installs the application.

**Using curl:**
```bash
curl -sSL "https://raw.githubusercontent.com/mithunkhatri/modulation/main/install.sh?v=$(date +%s)" | bash
```

**Using wget:**
```bash
wget -qO- "https://raw.githubusercontent.com/mithunkhatri/modulation/main/install.sh?v=$(date +%s)" | bash
```

### Manual Installation (Alternative)
If you prefer to build from source manually:

1. **Install FFmpeg**: Required for audio decoding ([Download here](https://ffmpeg.org/download.html)).
2. **Build**:
```bash
go build -o modulation .
```

## Uninstallation 🗑️

If you need to uninstall Modulation, you can run the following command:

```bash
curl -sSL https://raw.githubusercontent.com/mithunkhatri/modulation/main/uninstall.sh | bash
```

Alternatively, if you have the source code locally:
```bash
./uninstall.sh
```

## Usage ⌨️

Run the app:
```bash
./modulation
```

### Shortcuts
| Key | Action |
| --- | --- |
| `Enter`/`Space` | Play selected station (or toggle Pause if already playing) |
| `p` | Toggle Pause/Resume |
| `s` | Stop Playback |
| `a` | Add/Remove station from Favorites |
| `v` | Toggle Favorites View |
| `/` | Search stations |
| `f` | Filter by category |
| `b` | Toggle Background (Minimized) mode |
| `m` | Refresh Radio Browser Mirrors |
| `c` | Clear error messages |
| `1-9` | Quickly select the first 9 visible stations |
| `q` | Quit |


## Credits ☕

Modulation (v1.0.0) | Made with **Go** by **Mithun Khatri**
