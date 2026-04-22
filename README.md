# convert

A fast and versatile CLI tool for converting files between various formats.

## Features

- **Image Conversion**: Supports PNG, JPEG, WEBP, GIF, and TIFF.
- **Data Conversion**: Supports JSON, CSV, and YAML with deterministic header ordering.
- **Document Conversion**: Supports PDF, DOCX, HTML, and Markdown (requires LibreOffice). Supports PDF to plain text and Markdown extraction (requires poppler-utils).

## Installation

Ensure you have [Go](https://golang.org/doc/install) installed.

```bash
go install github.com/mintoleda/convert@latest
```

### Dependencies

For document conversions (e.g., `.docx` to `.pdf`), this tool requires **LibreOffice** (`soffice`) to be installed and available in your system's PATH.

For PDF text extraction (`.pdf` to `.txt` or `.md`), this tool requires **poppler-utils** (`pdftotext`) to be installed and available in your system's PATH.

```bash
# Ubuntu/Debian
sudo apt install poppler-utils

# macOS
brew install poppler
```

## Usage

### Basic Command

```bash
convert <input_file> <output_file>
```

### Examples

#### Data
```bash
# Convert JSON to CSV
convert data.json data.csv

# Convert CSV to YAML
convert data.csv data.yaml
```

#### Images
```bash
# Convert PNG to JPEG
convert image.png image.jpg

# Convert JPEG to WEBP
convert image.jpg image.webp
```

#### Documents
```bash
# Convert DOCX to PDF
convert document.docx document.pdf

# Convert Markdown to PDF
convert README.md README.pdf

# Extract text from PDF
convert document.pdf document.txt

# Extract PDF content as Markdown
convert document.pdf document.md
```

### Options

- `--list`: List all supported format conversions.
- `--force`: Overwrite the output file if it already exists.
- `-h, --help`: Show help message.

## Development

### Running Tests

```bash
go test -v ./...
```

### Adding New Converters

New converters can be added by implementing a conversion function and registering it in the `converter` package using `Register(from, to, func)`.
