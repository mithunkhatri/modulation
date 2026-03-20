# Modulation v1.0.2 🚀

- **Fixed Installation Script**: Resolved a syntax error in `install.sh` that affected macOS users.
- **Robust CI/CD**: Fixed GitHub Actions to support `linux/arm64` cross-compilation with proper CGO and multi-arch dependencies.
- **Improved Uninstall**: Added an `uninstall.sh` script for easy removal.

# Modulation v1.0.1 🎙️

- **Robust Radio Support**: Automatically rotates through multiple Radio Browser mirrors on failure.
- **Mirror Refresh**: New `m` shortcut to manually refresh the list of available API mirrors.
- **Improved Reliability**: Better handling of flaky streams and mirror-specific issues.

# Modulation v1.0.0 🎙️

Modulation is a high-performance, cross-platform TUI radio player that brings any station in the world directly to your terminal.

## Key Features ✨

- **Modern TUI**: Responsive, `htop`-style interface with real-time updates and smooth animations.
- **Global Radio Access**: Fetches the most popular stations from the Radio Browser API.
- **Powerhouse Audio Engine**: Uses FFmpeg for robust decoding of any stream format (MP3, AAC, FLAC, etc.).
- **Smart Management**: Full support for searching, category filtering, and personal favorites.
- **Cross-Platform Support**: Built for macOS, Linux, and Windows, with idiomatic configuration handling.
- **Single-Command Installation**: A universal installer script that handles all dependencies and OS-specific setup.

## Installation 📦

Install across all platforms with one command:
```bash
curl -sSL "https://raw.githubusercontent.com/mithunkhatri/modulation/main/install.sh?v=$(date +%s)" | bash
```

Enjoy your music! 🎵
