package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	Register("docx", "pdf", libreOfficeConvert)
	Register("md", "pdf", libreOfficeConvert)
	Register("html", "pdf", libreOfficeConvert)
	Register("pdf", "docx", libreOfficeConvert)
	Register("pdf", "txt", pdfToTxt)
	Register("pdf", "md", pdfToMd)
}

// libreOfficeConvert uses LibreOffice headless to convert documents.
func libreOfficeConvert(input, output string) error {
	if _, err := exec.LookPath("soffice"); err != nil {
		return fmt.Errorf("LibreOffice is required for this conversion but 'soffice' was not found on PATH. Install LibreOffice: https://www.libreoffice.org/download")
	}

	outExt := filepath.Ext(output)
	if len(outExt) > 0 {
		outExt = outExt[1:] // strip leading dot
	}

	inExt := filepath.Ext(input)
	if len(inExt) > 0 {
		inExt = inExt[1:]
	}

	absInput, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("failed to resolve input path: %w", err)
	}

	absOutput, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	outDir := filepath.Dir(absOutput)

	// LibreOffice needs explicit filters for certain conversions.
	convertArg := outExt
	filters := map[string]map[string]string{
		"pdf":  {"docx": "MS Word 2007 XML"},
		"docx": {"pdf": "writer_pdf_Export"},
		"md":   {"pdf": "writer_pdf_Export"},
		"html": {"pdf": "writer_pdf_Export"},
	}
	if filterMap, ok := filters[inExt]; ok {
		if filter, ok := filterMap[outExt]; ok {
			convertArg = outExt + ":" + filter
		}
	}

	cmd := exec.Command("soffice", "--headless", "--convert-to", convertArg, "--outdir", outDir, absInput)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("LibreOffice conversion failed: %w", err)
	}

	// LibreOffice names the output based on the input filename. Rename if needed.
	baseName := filepath.Base(absInput)
	ext := filepath.Ext(baseName)
	loOutput := filepath.Join(outDir, baseName[:len(baseName)-len(ext)]+"."+outExt)

	if loOutput != absOutput {
		if err := os.Rename(loOutput, absOutput); err != nil {
			return fmt.Errorf("failed to rename output file: %w", err)
		}
	}

	return nil
}

func pdfToTxt(input, output string) error {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return fmt.Errorf("pdftotext is required but not found on PATH (install poppler-utils)")
	}
	absInput, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("failed to resolve input path: %w", err)
	}
	absOutput, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}
	cmd := exec.Command("pdftotext", "-nopgbrk", "-enc", "UTF-8", absInput, absOutput)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pdftotext conversion failed: %w", err)
	}
	return nil
}

func pdfToMd(input, output string) error {
	if _, err := exec.LookPath("pdftotext"); err != nil {
		return fmt.Errorf("pdftotext is required but not found on PATH (install poppler-utils)")
	}
	tmp, err := os.CreateTemp("", "cv-pdf2md-*.txt")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

	absInput, err := filepath.Abs(input)
	if err != nil {
		return fmt.Errorf("failed to resolve input path: %w", err)
	}
	absOutput, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	cmd := exec.Command("pdftotext", "-nopgbrk", "-enc", "UTF-8", absInput, tmpPath)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pdftotext extraction failed: %w", err)
	}

	raw, err := os.ReadFile(tmpPath)
	if err != nil {
		return fmt.Errorf("failed to read extracted text: %w", err)
	}
	return os.WriteFile(absOutput, []byte(textToMarkdown(string(raw))), 0644)
}

func textToMarkdown(text string) string {
	lines := strings.Split(text, "\n")
	var buf strings.Builder
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			buf.WriteString("\n")
		} else if isHeading(trimmed, i, lines) {
			buf.WriteString("## ")
			buf.WriteString(trimmed)
			buf.WriteString("\n\n")
		} else {
			buf.WriteString(trimmed)
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func isHeading(line string, idx int, lines []string) bool {
	if len(line) > 80 {
		return false
	}
	for _, suffix := range []string{".", ",", ";", ":"} {
		if strings.HasSuffix(line, suffix) {
			return false
		}
	}
	nextIsBlank := idx+1 >= len(lines) || strings.TrimSpace(lines[idx+1]) == ""
	return nextIsBlank
}
