package aggregator

import (
	"fmt"

	"clientingo/internal/models"
)

func AggregateResults(data []models.DeviceData) {
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

	fmt.Println("==================================================")
	fmt.Println("AGGREGATION REPORT")
	fmt.Println("==================================================")
	fmt.Printf("Total Devices: %d\n", len(data))
	fmt.Printf("Online: %d (%.1f%%)\n", onlineCount, float64(onlineCount)/float64(len(data))*100)
	fmt.Printf("Offline: %d (%.1f%%)\n", offlineCount, float64(offlineCount)/float64(len(data))*100)
	fmt.Printf("Total Power: %.2f kW\n", totalPower)
	fmt.Printf("Average Power: %.2f kW per device\n", totalPower/float64(len(data)))
	fmt.Println("==================================================")
}
