package test

import (
	"io"
	"net/http"
	"testing"
	"time"
)

// PerformanceResult represents the performance test result
type PerformanceResult struct {
	RequestsPerSecond float64
	AverageTime       time.Duration
	MaxTime           time.Duration
}

// testPerformance tests the performance of the server
func testPerformance() (PerformanceResult, error) {
	// Test parameters
	baseURL := "http://localhost:8081"
	requestCount := 1000
	concurrency := 10

	// Create a client
	client := &http.Client{}

	// Start time
	start := time.Now()

	// Channels for results
	done := make(chan bool)
	error := make(chan error)

	// Run concurrent requests
	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < requestCount/concurrency; j++ {
				resp, err := client.Get(baseURL + "/api/v1/users")
				if err != nil {
					error <- err
					return
				}
				defer resp.Body.Close()
				_, err = io.ReadAll(resp.Body)
				if err != nil {
					error <- err
					return
				}
			}
			done <- true
		}()
	}

	// Wait for all requests to complete
	completed := 0
	for completed < concurrency {
		select {
		case <-done:
			completed++
		case err := <-error:
			return PerformanceResult{}, err
		}
	}

	// Calculate results
	duration := time.Since(start)
	requestsPerSecond := float64(requestCount) / duration.Seconds()
	averageTime := duration / time.Duration(requestCount)

	return PerformanceResult{
		RequestsPerSecond: requestsPerSecond,
		AverageTime:       averageTime,
		MaxTime:           duration,
	}, nil
}

// TestServerPerformance tests the performance of the server
func TestServerPerformance(t *testing.T) {
	t.Logf("\n=== Testing Server Performance ===")
	result, err := testPerformance()
	if err != nil {
		t.Fatalf("Error testing server: %v", err)
	}

	t.Logf("Requests per second: %.2f", result.RequestsPerSecond)
	t.Logf("Average response time: %v", result.AverageTime)
	t.Logf("Total time: %v", result.MaxTime)
}
