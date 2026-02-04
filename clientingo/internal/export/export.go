package export

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"clientingo/internal/models"
)

func ExportToJSON(data []models.DeviceData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create JSON file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}

func ExportDetailedReport(data []models.DeviceData, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()

	// Header
	fmt.Fprintf(file, "EnergyGrid Solar Inverter Telemetry Report\n")
	fmt.Fprintf(file, "Generated: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "Total Devices: %d\n", len(data))
	fmt.Fprintf(file, "%s\n\n", "=================================================================")

	// Statistics
	var onlineCount, offlineCount int
	var totalPower float64

	for _, device := range data {
		if device.Status == "Online" {
			onlineCount++
		} else {
			offlineCount++
		}

		var power float64
		fmt.Sscanf(device.Power, "%f", &power)
		totalPower += power
	}

	fmt.Fprintf(file, "SUMMARY STATISTICS\n")
	fmt.Fprintf(file, "%s\n", "-----------------------------------------------------------------")
	fmt.Fprintf(file, "Online Devices:       %d (%.1f%%)\n", onlineCount, float64(onlineCount)/float64(len(data))*100)
	fmt.Fprintf(file, "Offline Devices:      %d (%.1f%%)\n", offlineCount, float64(offlineCount)/float64(len(data))*100)
	fmt.Fprintf(file, "Total Power Output:   %.2f kW\n", totalPower)
	fmt.Fprintf(file, "Average Power/Device: %.2f kW\n", totalPower/float64(len(data)))
	fmt.Fprintf(file, "\n%s\n\n", "=================================================================")

	// Detailed device list
	fmt.Fprintf(file, "DETAILED DEVICE INFORMATION\n")
	fmt.Fprintf(file, "%s\n", "-----------------------------------------------------------------")
	fmt.Fprintf(file, "%-10s %-12s %-10s %-25s\n", "Serial #", "Power", "Status", "Last Updated")
	fmt.Fprintf(file, "%s\n", "-----------------------------------------------------------------")

	for _, device := range data {
		fmt.Fprintf(file, "%-10s %-12s %-10s %-25s\n",
			device.SN,
			device.Power,
			device.Status,
			device.LastUpdated,
		)
	}

	fmt.Fprintf(file, "%s\n", "=================================================================")

	return nil
}

// ExportAll exports data in multiple formats
func ExportAll(data []models.DeviceData) error {
	timestamp := time.Now().Format("20060102_150405")

	jsonFilename := fmt.Sprintf("energygrid_devices_%s.json", timestamp)
	if err := ExportToJSON(data, jsonFilename); err != nil { //json export
		return err
	}
	fmt.Printf("--> JSON exported: %s\n", jsonFilename)

	reportFilename := fmt.Sprintf("energygrid_report_%s.txt", timestamp)
	if err := ExportDetailedReport(data, reportFilename); err != nil { //txt report
		return err
	}
	fmt.Printf("--> Detailed report exported: %s\n", reportFilename)

	return nil
}
