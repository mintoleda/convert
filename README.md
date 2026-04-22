# cv

A fast and versatile CLI tool for converting files between various formats.

## Features

- **Image Conversion**: Supports PNG, JPEG, WEBP, GIF, and TIFF.
- **Data Conversion**: Supports JSON, CSV, and YAML with deterministic header ordering.
- **Document Conversion**: Supports PDF, DOCX, HTML, and Markdown (requires LibreOffice). Supports PDF to plain text and Markdown extraction (requires poppler-utils).

## Installation

Ensure you have [Go](https://golang.org/doc/install) installed.

```bash
go install github.com/mintoleda/convert/cv@latest
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
cv <input_file> <output_file>
```

### Examples

#### Data
```bash
# Convert JSON to CSV
cv data.json data.csv

# Convert CSV to YAML
cv data.csv data.yaml
```

#### Images
```bash
# Convert PNG to JPEG
cv image.png image.jpg

# Convert JPEG to WEBP
cv image.jpg image.webp
```

#### Documents
```bash
# Convert DOCX to PDF
cv document.docx document.pdf

# Convert Markdown to PDF
cv README.md README.pdf

# Extract text from PDF
cv document.pdf document.txt

# Extract PDF content as Markdown
cv document.pdf document.md
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
