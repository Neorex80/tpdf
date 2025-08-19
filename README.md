
# tpdf - Terminal PDF Viewer

`tpdf` is a terminal-based PDF viewer written in Go. It converts PDF files to plain text, formats them with ASCII layout, and allows navigation through pages and chapters using terminal commands.

## Features

- Cross-platform PDF to text conversion using `pdftotext` (Poppler)
- Automatic detection of pdftotext on Windows, macOS, and Linux
- Displays formatted content with ASCII styling
- Short command aliases for faster navigation:
  - `n` or `next` – Next page
  - `p` or `prev` – Previous page  
  - `g N` or `gotoPage N` – Jump to page N
  - `gc Title` or `gotoChapter Title` – Jump to chapter
  - `q` or `exit` – Quit viewer
  - `h` or `help` – Show help

## Requirements

- Go 1.24.5 or higher
- `pdftotext` from the `poppler-utils` package

### Install pdftotext

On Debian/Ubuntu:

```bash
sudo apt install poppler-utils
```

On Arch Linux:

```bash
sudo pacman -S poppler
```

On Windows:

1. Download Poppler for Windows from: https://github.com/oschwartz10612/poppler-windows/releases/
2. Extract the downloaded archive to a folder (e.g., `C:\poppler`)
## Installation

Clone the repository:

```bash
git clone https://github.com/prem0x01/tpdf.git
cd tpdf
```

Build the project:

```bash
go build -o tpdf main.go
```

## Usage

```bash
./tpdf your_file.pdf
```

Once inside the viewer, use these commands:

* `n` or `next` – Next page
* `p` or `prev` – Previous page
* `g 3` or `gotoPage 3` – Go to page 3
* `gc Introduction` – Jump to chapter "Introduction"
* `q` or `exit` – Quit viewer
* `h` or `help` – Show all commands


## Improvements

This optimized version includes:

- **Cross-platform compatibility** - Auto-detects pdftotext on Windows/macOS/Linux
- **Better error handling** - Clear error messages with installation instructions
- **Shorter commands** - Single-letter aliases (n, p, g, q, h) for faster navigation
- **Robust file handling** - Validates PDF existence and conversion success
- **Minimal codebase** - Single file, no external dependencies beyond Go stdlib
- **Automatic cleanup** - Removes temporary files on exit
