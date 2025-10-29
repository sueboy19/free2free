package testutils

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// TestResult represents the result of a single test
type TestResult struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"` // "PASS", "FAIL", "SKIP"
	Duration  string    `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

// TestSuiteResult represents the result of a test suite
type TestSuiteResult struct {
	Name      string       `json:"name"`
	Timestamp time.Time    `json:"timestamp"`
	Total     int          `json:"total"`
	Passed    int          `json:"passed"`
	Failed    int          `json:"failed"`
	Skipped   int          `json:"skipped"`
	Duration  string       `json:"duration"`
	Results   []TestResult `json:"results"`
}

// ResultReporter handles reporting of test results
type ResultReporter struct {
	suiteName   string
	startTime   time.Time
	testResults []TestResult
	outputFile  string
}

// NewResultReporter creates a new result reporter
func NewResultReporter(suiteName, outputFile string) *ResultReporter {
	return &ResultReporter{
		suiteName:   suiteName,
		startTime:   time.Now(),
		testResults: []TestResult{},
		outputFile:  outputFile,
	}
}

// AddTestResult adds a test result to the reporter
func (rr *ResultReporter) AddTestResult(name string, status string, duration time.Duration, err error) {
	result := TestResult{
		Name:      name,
		Status:    status,
		Duration:  duration.String(),
		Timestamp: time.Now(),
	}

	if err != nil {
		result.Error = err.Error()
	}

	rr.testResults = append(rr.testResults, result)
}

// GetSuiteResult returns the test suite result
func (rr *ResultReporter) GetSuiteResult() TestSuiteResult {
	// Calculate counts
	var passed, failed, skipped int
	for _, result := range rr.testResults {
		switch result.Status {
		case "PASS":
			passed++
		case "FAIL":
			failed++
		case "SKIP":
			skipped++
		}
	}

	// Calculate total duration
	totalDuration := time.Since(rr.startTime)

	return TestSuiteResult{
		Name:      rr.suiteName,
		Timestamp: rr.startTime,
		Total:     len(rr.testResults),
		Passed:    passed,
		Failed:    failed,
		Skipped:   skipped,
		Duration:  totalDuration.String(),
		Results:   rr.testResults,
	}
}

// WriteJSONReport writes the test results as JSON to the specified file
func (rr *ResultReporter) WriteJSONReport() error {
	suiteResult := rr.GetSuiteResult()

	file, err := os.Create(rr.outputFile)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(suiteResult); err != nil {
		return fmt.Errorf("failed to encode report to JSON: %w", err)
	}

	return nil
}

// WriteSimpleReport writes a simple text report to the specified file
func (rr *ResultReporter) WriteSimpleReport() error {
	suiteResult := rr.GetSuiteResult()

	file, err := os.Create(rr.outputFile)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer file.Close()

	// Write summary
	summary := fmt.Sprintf("Test Suite: %s\n", suiteResult.Name)
	summary += fmt.Sprintf("Executed: %d tests\n", suiteResult.Total)
	summary += fmt.Sprintf("Passed: %d tests\n", suiteResult.Passed)
	summary += fmt.Sprintf("Failed: %d tests\n", suiteResult.Failed)
	summary += fmt.Sprintf("Skipped: %d tests\n", suiteResult.Skipped)
	summary += fmt.Sprintf("Total Duration: %s\n", suiteResult.Duration)
	summary += fmt.Sprintf("Completed at: %s\n\n", suiteResult.Timestamp.Format(time.RFC3339))

	if _, err := file.WriteString(summary); err != nil {
		return fmt.Errorf("failed to write summary to report: %w", err)
	}

	// Write individual test results
	for _, result := range suiteResult.Results {
		testLine := fmt.Sprintf("%s [%s] (%s)", result.Name, result.Status, result.Duration)
		if result.Status == "FAIL" && result.Error != "" {
			testLine += fmt.Sprintf(" - Error: %s", result.Error)
		}
		testLine += "\n"

		if _, err := file.WriteString(testLine); err != nil {
			return fmt.Errorf("failed to write test result to report: %w", err)
		}
	}

	return nil
}

// PrintSummary prints a summary to stdout
func (rr *ResultReporter) PrintSummary() {
	suiteResult := rr.GetSuiteResult()

	fmt.Printf("\n=== Test Suite Summary ===\n")
	fmt.Printf("Suite: %s\n", suiteResult.Name)
	fmt.Printf("Total: %d | Passed: %d | Failed: %d | Skipped: %d\n",
		suiteResult.Total, suiteResult.Passed, suiteResult.Failed, suiteResult.Skipped)
	fmt.Printf("Duration: %s\n", suiteResult.Duration)
	fmt.Printf("========================\n\n")

	// Print failed tests
	if suiteResult.Failed > 0 {
		fmt.Printf("Failed Tests:\n")
		for _, result := range suiteResult.Results {
			if result.Status == "FAIL" {
				fmt.Printf("  ❌ %s - %s\n", result.Name, result.Error)
			}
		}
		fmt.Printf("\n")
	}
}
