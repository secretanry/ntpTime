package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/beevik/ntp"
)

// TestMain sets up test environment
func TestMain(m *testing.M) {
	log.SetOutput(&bytes.Buffer{})

	code := m.Run()

	os.Exit(code)
}

// TestGetNTPTime tests the getNTPTime function
func TestGetNTPTime(t *testing.T) {
	timeStr, err := getNTPTime("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		t.Skipf("Skipping test due to NTP server error: %v", err)
	}

	if timeStr == "" {
		t.Error("Expected non-empty time string")
	}

	_, parseErr := time.Parse(time.RFC3339, timeStr)
	if parseErr != nil {
		t.Errorf("Time string %s is not in RFC3339 format: %v", timeStr, parseErr)
	}

	parsedTime, _ := time.Parse(time.RFC3339, timeStr)
	expectedMin := time.Now().Add(-1 * time.Hour)
	expectedMax := time.Now().Add(1 * time.Hour)

	if parsedTime.Before(expectedMin) || parsedTime.After(expectedMax) {
		t.Errorf("NTP time %v is outside reasonable range [%v, %v]",
			parsedTime, expectedMin, expectedMax)
	}
}

// TestGetNTPTimeError tests error handling in getNTPTime function
func TestGetNTPTimeError(t *testing.T) {
	_, err := getNTPTime("invalid-server.example.com")
	if err == nil {
		t.Error("Expected error with invalid server, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") &&
		!strings.Contains(err.Error(), "no such host") &&
		!strings.Contains(err.Error(), "connection refused") {
		t.Logf("Error message: %v", err.Error())
	}
}

// TestNTPTimeRetrieval tests successful NTP time retrieval
func TestNTPTimeRetrieval(t *testing.T) {
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		t.Skipf("Skipping test due to NTP server error: %v", err)
	}

	if response.ClockOffset == 0 && response.RTT == 0 {
		t.Error("Expected NTP response to contain valid clock offset or rtt")
	}

	ntpTime := time.Now().Add(response.ClockOffset)
	formatted := ntpTime.Format(time.RFC3339)

	if formatted == "" {
		t.Error("Formatted time should not be empty")
	}

	expectedMin := time.Now().Add(-1 * time.Hour)
	expectedMax := time.Now().Add(1 * time.Hour)

	if ntpTime.Before(expectedMin) || ntpTime.After(expectedMax) {
		t.Errorf("NTP time %v is outside reasonable range [%v, %v]",
			ntpTime, expectedMin, expectedMax)
	}
}

// TestNTPErrorHandling tests error handling with invalid server
func TestNTPErrorHandling(t *testing.T) {
	_, err := ntp.Query("invalid-server.example.com")
	if err == nil {
		t.Error("Expected error with invalid server, got nil")
	}

	if !strings.Contains(err.Error(), "timeout") &&
		!strings.Contains(err.Error(), "no such host") &&
		!strings.Contains(err.Error(), "connection refused") {
		t.Logf("Error message: %v", err.Error())
	}
}

// TestNTPTimeFormat tests the time formatting
func TestNTPTimeFormat(t *testing.T) {
	// Test RFC3339 format
	testTime := time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC)
	formatted := testTime.Format(time.RFC3339)
	expected := "2024-01-15T12:30:45Z"

	if formatted != expected {
		t.Errorf("Expected format %s, got %s", expected, formatted)
	}
}

// TestNTPResponseFields tests NTP response structure
func TestNTPResponseFields(t *testing.T) {
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		t.Skipf("Skipping test due to NTP server error: %v", err)
	}

	if response.Time.IsZero() {
		t.Error("Time should not be zero")
	}

	if response.ReferenceTime.IsZero() {
		t.Error("ReferenceTime should not be zero")
	}

	if err := response.Validate(); err != nil {
		t.Errorf("Response validation failed: %v", err)
	}

	if response.Stratum < 1 || response.Stratum > 15 {
		t.Errorf("Unexpected stratum level: %d", response.Stratum)
	}

	if response.Version < 3 || response.Version > 4 {
		t.Errorf("Unexpected NTP version: %d", response.Version)
	}
}

// TestNTPMultipleServers tests with different NTP servers
func TestNTPMultipleServers(t *testing.T) {
	servers := []string{
		"0.beevik-ntp.pool.ntp.org",
		"1.beevik-ntp.pool.ntp.org",
		"2.beevik-ntp.pool.ntp.org",
	}

	for _, server := range servers {
		t.Run(server, func(t *testing.T) {
			response, err := ntp.Query(server)
			if err != nil {
				t.Skipf("Server %s unavailable: %v", server, err)
			}

			if response.ClockOffset == 0 && response.RTT == 0 {
				t.Error("Expected valid response data")
			}
		})
	}
}

// TestNTPPerformance tests reasonable response times
func TestNTPPerformance(t *testing.T) {
	start := time.Now()
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	duration := time.Since(start)

	if err != nil {
		t.Skipf("Skipping test due to NTP server error: %v", err)
	}

	if duration > 5*time.Second {
		t.Errorf("NTP query took too long: %v", duration)
	}

	if response.RTT > time.Second {
		t.Errorf("rtt too high: %v", response.RTT)
	}
}

// BenchmarkNTPQuery benchmarks NTP query performance
func BenchmarkNTPQuery(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
		if err != nil {
			b.Skipf("NTP server error: %v", err)
		}
	}
}
