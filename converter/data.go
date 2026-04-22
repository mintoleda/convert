package converter

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"gopkg.in/yaml.v3"
)

func init() {
	Register("csv", "json", csvToJSON)
	Register("csv", "yaml", csvToYAML)
	Register("json", "csv", jsonToCSV)
	Register("json", "yaml", jsonToYAML)
	Register("yaml", "json", yamlToJSON)
}

// csvToRecords reads a CSV file and returns a slice of maps (header → value).
func csvToRecords(input string) ([]map[string]interface{}, error) {
	f, err := os.Open(input)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read csv: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("csv file must have a header row and at least one data row")
	}

	headers := rows[0]
	records := make([]map[string]interface{}, 0, len(rows)-1)
	for _, row := range rows[1:] {
		record := make(map[string]interface{}, len(headers))
		for i, header := range headers {
			if i < len(row) {
				record[header] = row[i]
			}
		}
		records = append(records, record)
	}
	return records, nil
}

func csvToJSON(input, output string) error {
	records, err := csvToRecords(input)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	return os.WriteFile(output, data, 0644)
}

func csvToYAML(input, output string) error {
	records, err := csvToRecords(input)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(records)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	return os.WriteFile(output, data, 0644)
}

func jsonToCSV(input, output string) error {
	raw, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read json: %w", err)
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(raw, &records); err != nil {
		return fmt.Errorf("failed to parse json (expected array of objects): %w", err)
	}

	if len(records) == 0 {
		return fmt.Errorf("json array is empty")
	}

	// Collect all unique keys for headers.
	headerSet := map[string]bool{}
	for _, r := range records {
		for k := range r {
			headerSet[k] = true
		}
	}
	headers := make([]string, 0, len(headerSet))
	for k := range headerSet {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create csv: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write(headers); err != nil {
		return err
	}

	for _, r := range records {
		row := make([]string, len(headers))
		for i, h := range headers {
			if v, ok := r[h]; ok && v != nil {
				row[i] = fmt.Sprintf("%v", v)
			} else {
				row[i] = ""
			}
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()
	return w.Error()
}

func jsonToYAML(input, output string) error {
	raw, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read json: %w", err)
	}

	var data interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return fmt.Errorf("failed to parse json: %w", err)
	}

	out, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %w", err)
	}

	return os.WriteFile(output, out, 0644)
}

func yamlToJSON(input, output string) error {
	raw, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read yaml: %w", err)
	}

	var data interface{}
	if err := yaml.Unmarshal(raw, &data); err != nil {
		return fmt.Errorf("failed to parse yaml: %w", err)
	}

	// yaml.v3 unmarshals map keys as string, but we normalize to ensure JSON compat.
	data = normalizeYAML(data)

	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	return os.WriteFile(output, out, 0644)
}

// normalizeYAML converts map[string]interface{} (from yaml.v3) recursively for JSON compat.
func normalizeYAML(v interface{}) interface{} {
	switch v := v.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{}, len(v))
		for k, val := range v {
			m[k] = normalizeYAML(val)
		}
		return m
	case []interface{}:
		for i, val := range v {
			v[i] = normalizeYAML(val)
		}
		return v
	default:
		return v
	}
}
