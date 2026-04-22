package converter

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// ConvertFunc performs a conversion from input file to output file.
type ConvertFunc func(input, output string) error

type convKey struct {
	from, to string
}

var registry = map[convKey]ConvertFunc{}

// Register adds a converter for the given extension pair (without dots).
func Register(fromExt, toExt string, fn ConvertFunc) {
	registry[convKey{from: strings.ToLower(fromExt), to: strings.ToLower(toExt)}] = fn
}

// Convert looks up the appropriate converter by file extensions and runs it.
func Convert(inputPath, outputPath string) error {
	fromExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(inputPath)), ".")
	toExt := strings.TrimPrefix(strings.ToLower(filepath.Ext(outputPath)), ".")

	if fromExt == "" || toExt == "" {
		return fmt.Errorf("cannot determine file extensions: input=%q output=%q", inputPath, outputPath)
	}

	fn, ok := registry[convKey{from: fromExt, to: toExt}]
	if !ok {
		return fmt.Errorf("unsupported conversion: %s → %s", fromExt, toExt)
	}

	return fn(inputPath, outputPath)
}

// ListSupported returns a sorted list of "ext → ext" strings for all registered conversions.
func ListSupported() []string {
	pairs := make([]string, 0, len(registry))
	for k := range registry {
		pairs = append(pairs, fmt.Sprintf("%s → %s", k.from, k.to))
	}
	sort.Strings(pairs)
	return pairs
}
