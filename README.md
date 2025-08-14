````markdown
# tpdf - Terminal PDF Viewer

`tpdf` is a terminal-based PDF viewer written in Go. It converts PDF files to plain text, formats them with ASCII layout, and allows navigation through pages and chapters using terminal commands.

## Features

- Converts PDF to text using `pdftotext` (Poppler)
- Displays formatted content with ASCII styling
- Supports:
  - `gotoPage N` – Jump to a specific page
  - `gotoChapter Title` – Jump to a specific chapter
  - `next` – Move to the next page
  - `previous` – Move to the previous page
  - `exit` or `q` – Exit the viewer

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

Once inside the viewer, use the following commands:

* `gotoPage 3` – Go to page 3
* `gotoChapter Introduction` – Jump to the chapter named "Introduction"
* `next` – Next page
* `previous` – Previous page
* `exit` or `q` – Quit the viewer


## Improvements

This version of tpdf has been significantly optimized:

- Removed unnecessary dependencies (Bubbletea TUI framework)
- Consolidated code into a single file for easier maintenance
- Reduced binary size and compilation time
- Improved memory usage
- Simplified code structure while maintaining all functionality
