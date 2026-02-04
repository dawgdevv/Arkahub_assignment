package utils

import (
	"fmt"

	"clientingo/config"
)

// GenerateSerialNumbers creates SN-000 to SN-499
func GenerateSerialNumbers() []string {
	sns := make([]string, config.TotalSNs)
	for i := 0; i < config.TotalSNs; i++ {
		sns[i] = fmt.Sprintf("SN-%03d", i)
	}
	return sns
}

// BatchSerialNumbers splits SNs into batches of 10
func BatchSerialNumbers(sns []string) [][]string {
	var batches [][]string

	for i := 0; i < len(sns); i += config.BatchSize {
		end := i + config.BatchSize
		if end > len(sns) {
			end = len(sns)
		}
		batches = append(batches, sns[i:end])
	}
	return batches
}
