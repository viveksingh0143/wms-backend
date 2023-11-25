package utils

import (
	"encoding/csv"
	"os"
)

func ProcessCsvFile(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var headers []string
	var processedRows []map[string]string

	for i, row := range records {
		if i == 0 {
			// Process headers
			for _, header := range row {
				headers = append(headers, GenerateSlug(header))
			}
		} else {
			// Process each row
			rowMap := make(map[string]string)
			for j, cell := range row {
				rowMap[headers[j]] = cell
			}
			processedRows = append(processedRows, rowMap)
		}
	}

	return processedRows, nil
}
