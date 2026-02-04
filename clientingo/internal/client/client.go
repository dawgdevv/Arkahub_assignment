package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"clientingo/config"
	"clientingo/internal/auth"
	"clientingo/internal/models"
	"clientingo/internal/utils"
)

type Client struct {
	httpClient  *http.Client
	lastReqTime time.Time
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		lastReqTime: time.Time{},
	}
}

func (c *Client) fetchBatch(batch []string, attempt int) ([]models.DeviceData, error) {
	if !c.lastReqTime.IsZero() {
		elapsed := time.Since(c.lastReqTime)
		if elapsed < time.Second {
			time.Sleep(time.Second - elapsed)
		}
	}

	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	signature := auth.GenerateSignature("/device/real/query", config.Token, timestamp)

	reqBody := models.APIRequest{SNList: batch}
	jsonData, err := json.Marshal(reqBody)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", config.APIURL, bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("content-Type", "application/json")
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("signature", signature)

	c.lastReqTime = time.Now()
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		var apiResp models.APIResponse

		if err := json.Unmarshal(body, &apiResp); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		return apiResp.Data, nil

	case http.StatusTooManyRequests:
		if attempt < 3 {
			fmt.Printf("Rate limited (429), retrying... (attempt %d/3)\n", attempt+1)
			time.Sleep(2 * time.Second)
			return c.fetchBatch(batch, attempt+1)
		}

		return nil, fmt.Errorf("rate limit exceeded after 3 attempts")

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("authentication failed (401): %s", string(body))

	default:
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}
}

func (c *Client) FetchAllDevices() ([]models.DeviceData, error) {
	fmt.Println(" Starting EnergyGrid Data Aggregator")
	fmt.Printf(" Total devices: %d | Batch size: %d | Rate limit: 1 req/sec\n\n", config.TotalSNs, config.BatchSize)

	sns := utils.GenerateSerialNumbers()
	batches := utils.BatchSerialNumbers(sns)

	fmt.Printf("Processing %d batches...\n\n", len(batches))

	var allData []models.DeviceData
	startTime := time.Now()

	for i, batch := range batches {
		batchStart := time.Now()

		data, err := c.fetchBatch(batch, 1)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch batch %d: %w", i+1, err)
		}

		allData = append(allData, data...)

		elapsed := time.Since(batchStart).Seconds()
		progress := float64(i+1) / float64(len(batches)) * 100

		fmt.Printf("--> Batch %2d/%d | Devices: %3d-%3d | Time: %.2fs | Progress: %.1f%%\n",
			i+1, len(batches), i*config.BatchSize, (i+1)*config.BatchSize-1, elapsed, progress)
	}

	totalTime := time.Since(startTime).Seconds()
	fmt.Printf("\n⚡ Complete! Fetched %d devices in %.2f seconds\n\n", len(allData), totalTime)

	return allData, nil
}
