package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJsonToCSV(t *testing.T) {
	tmpDir := t.TempDir()
	input := filepath.Join(tmpDir, "input.json")
	output := filepath.Join(tmpDir, "output.csv")

	jsonData := `[
		{"name": "Alice", "age": 30, "city": "New York"},
		{"name": "Bob", "age": null, "city": "Los Angeles", "extra": "yes"},
		{"name": "Charlie", "city": "Chicago"}
	]`

	if err := os.WriteFile(input, []byte(jsonData), 0644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	if err := jsonToCSV(input, output); err != nil {
		t.Fatalf("jsonToCSV failed: %v", err)
	}

	csvData, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	// Expecting sorted headers: age, city, extra, name
	// Bob has null age -> empty string
	// Charlie has no age or extra -> empty string
	expectedCSV := `age,city,extra,name
30,New York,,Alice
,Los Angeles,yes,Bob
,Chicago,,Charlie
`
	if string(csvData) != expectedCSV {
		t.Errorf("expected CSV:\n%s\ngot:\n%s", expectedCSV, string(csvData))
	}
}

func TestCsvToJSON(t *testing.T) {
	tmpDir := t.TempDir()
	input := filepath.Join(tmpDir, "input.csv")
	output := filepath.Join(tmpDir, "output.json")

	csvData := `age,city,extra,name
30,New York,,Alice
,Los Angeles,yes,Bob
,Chicago,,Charlie
`

	if err := os.WriteFile(input, []byte(csvData), 0644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	if err := csvToJSON(input, output); err != nil {
		t.Fatalf("csvToJSON failed: %v", err)
	}

	jsonData, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	expectedJSON := `[
  {
    "age": "30",
    "city": "New York",
    "extra": "",
    "name": "Alice"
  },
  {
    "age": "",
    "city": "Los Angeles",
    "extra": "yes",
    "name": "Bob"
  },
  {
    "age": "",
    "city": "Chicago",
    "extra": "",
    "name": "Charlie"
  }
]`
	if string(jsonData) != expectedJSON {
		t.Errorf("expected JSON:\n%s\ngot:\n%s", expectedJSON, string(jsonData))
	}
}

func TestJsonToYAMLAndBack(t *testing.T) {
	tmpDir := t.TempDir()
	inputJSON := filepath.Join(tmpDir, "input.json")
	outputYAML := filepath.Join(tmpDir, "output.yaml")
	outputJSON := filepath.Join(tmpDir, "output2.json")

	jsonData := `[
  {
    "age": 30,
    "city": "New York",
    "name": "Alice"
  }
]`

	if err := os.WriteFile(inputJSON, []byte(jsonData), 0644); err != nil {
		t.Fatalf("failed to write input: %v", err)
	}

	if err := jsonToYAML(inputJSON, outputYAML); err != nil {
		t.Fatalf("jsonToYAML failed: %v", err)
	}

	yamlData, err := os.ReadFile(outputYAML)
	if err != nil {
		t.Fatalf("failed to read YAML output: %v", err)
	}

	expectedYAML := `- age: 30
  city: New York
  name: Alice
`
	if string(yamlData) != expectedYAML {
		t.Errorf("expected YAML:\n%s\ngot:\n%s", expectedYAML, string(yamlData))
	}

	if err := yamlToJSON(outputYAML, outputJSON); err != nil {
		t.Fatalf("yamlToJSON failed: %v", err)
	}

	jsonData2, err := os.ReadFile(outputJSON)
	if err != nil {
		t.Fatalf("failed to read JSON output: %v", err)
	}

	if string(jsonData2) != jsonData {
		t.Errorf("expected JSON:\n%s\ngot:\n%s", jsonData, string(jsonData2))
	}
}
