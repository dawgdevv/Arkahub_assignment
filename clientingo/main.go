package main

import (
	"fmt"

	"clientingo/internal/aggregator"
	"clientingo/internal/client"
	"clientingo/internal/export"
)

func main() {
	c := client.NewClient()

	data, err := c.FetchAllDevices()

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	aggregator.AggregateResults(data)

	// Export data to multiple formats
	fmt.Println("\n📊 Exporting data...")
	if err := export.ExportAll(data); err != nil {
		fmt.Printf("⚠️  Warning: Failed to export data: %v\n", err)
	} else {
		fmt.Println("✅ All exports completed successfully!")
	}
}
