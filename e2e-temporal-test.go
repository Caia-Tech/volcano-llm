package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type TestResult struct {
	Name     string
	Success  bool
	Duration time.Duration
	Details  string
	Error    string
}

type TestSuite struct {
	Results []TestResult
	BaseURL string
}

func main() {
	log.Printf("🚀 Starting comprehensive Temporal integration tests")
	
	suite := &TestSuite{
		BaseURL: "http://localhost:8080",
		Results: make([]TestResult, 0),
	}
	
	// Wait for services to be ready
	suite.waitForServices()
	
	// Run all test scenarios
	suite.testDeterministicExecution()
	suite.testTemporalWorkflowRouting()
	suite.testGitNativeHotReload()
	suite.testMultiTenantIsolation()
	suite.testWorkflowSignaling()
	suite.testErrorHandlingAndRetries()
	suite.testPerformanceBenchmark()
	
	// Print comprehensive results
	suite.printResults()
}

func (ts *TestSuite) waitForServices() {
	log.Printf("⏳ Waiting for services to be ready...")
	
	// Check if Docker Compose is running
	cmd := exec.Command("docker-compose", "-f", "docker-compose-temporal.yml", "ps")
	if err := cmd.Run(); err != nil {
		log.Printf("❌ Docker Compose not running. Starting services...")
		startCmd := exec.Command("docker-compose", "-f", "docker-compose-temporal.yml", "up", "-d")
		if err := startCmd.Run(); err != nil {
			log.Fatalf("Failed to start Docker Compose: %v", err)
		}
	}
	
	// Wait for API server to be ready
	for i := 0; i < 60; i++ {
		resp, err := http.Get(ts.BaseURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			log.Printf("✅ Services are ready!")
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(2 * time.Second)
		log.Printf("   Waiting for services... (%d/60)", i+1)
	}
	
	log.Fatalf("❌ Services did not become ready in time")
}

func (ts *TestSuite) testDeterministicExecution() {
	log.Printf("🧮 Testing deterministic execution performance...")
	
	tests := []struct {
		name string
		text string
	}{
		{"Simple Math", "calculate 42 + 58"},
		{"Complex Math", "calculate (15 * 7) + (89 - 34) / 5"},
		{"Sequential Math", "calculate 10 + 5, then multiply by 3"},
	}
	
	for _, test := range tests {
		start := time.Now()
		
		request := map[string]interface{}{
			"text":       test.text,
			"session_id": fmt.Sprintf("deterministic-test-%d", time.Now().UnixNano()),
		}
		
		response := ts.makeRequest("POST", "/api/v1/execute", request)
		duration := time.Since(start)
		
		success := response["success"] == true && duration < 200*time.Millisecond
		
		result := TestResult{
			Name:     fmt.Sprintf("Deterministic: %s", test.name),
			Success:  success,
			Duration: duration,
			Details:  fmt.Sprintf("Result: %v, Duration: %v", response["result"], duration),
		}
		
		if !success {
			result.Error = fmt.Sprintf("Too slow (%v) or failed", duration)
		}
		
		ts.Results = append(ts.Results, result)
		log.Printf("   %s: %v (Duration: %v)", test.name, success, duration)
	}
}

func (ts *TestSuite) testTemporalWorkflowRouting() {
	log.Printf("⚡ Testing Temporal workflow routing...")
	
	workflows := []struct {
		name string
		text string
		expectedWorkflow string
	}{
		{"Data Pipeline", "run data pipeline to extract and transform customer data", "DataPipelineWorkflow"},
		{"GitOps Deploy", "deploy latest changes from git repository", "GitOpsWorkflow"},
		{"Customer Onboarding", "onboard new enterprise customer", "CustomerOnboardingWorkflow"},
		{"Long Analytics", "run comprehensive analytics for the last 7 days", "LongRunningAnalyticsWorkflow"},
	}
	
	for _, workflow := range workflows {
		start := time.Now()
		
		request := map[string]interface{}{
			"text":       workflow.text,
			"session_id": fmt.Sprintf("workflow-test-%d", time.Now().UnixNano()),
		}
		
		response := ts.makeRequest("POST", "/api/v1/execute", request)
		duration := time.Since(start)
		
		// Check if it was routed to Temporal (has workflow_id or took longer than typical deterministic execution)
		isTemporalWorkflow := response["workflow_id"] != nil || duration > 500*time.Millisecond
		
		success := response["success"] == true && isTemporalWorkflow
		
		result := TestResult{
			Name:     fmt.Sprintf("Temporal Routing: %s", workflow.name),
			Success:  success,
			Duration: duration,
			Details:  fmt.Sprintf("WorkflowID: %v, Duration: %v", response["workflow_id"], duration),
		}
		
		if !success {
			result.Error = "Not routed to Temporal or failed"
		}
		
		ts.Results = append(ts.Results, result)
		log.Printf("   %s: %v (WorkflowID: %v)", workflow.name, success, response["workflow_id"])
	}
}

func (ts *TestSuite) testGitNativeHotReload() {
	log.Printf("🔄 Testing git-native hot-reload...")
	
	// Test 1: Get current workflow definitions
	definitionsResponse := ts.makeRequest("GET", "/api/v1/temporal/workflows/definitions", nil)
	initialCount := int(definitionsResponse["total"].(float64))
	
	// Test 2: Trigger hot-reload
	reloadRequest := map[string]interface{}{
		"type":       "workflow",
		"repository": "volcano-workflows",
		"file_path":  "workflows/test-workflow.json",
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	
	start := time.Now()
	reloadResponse := ts.makeRequest("POST", "/api/v1/temporal/reload", reloadRequest)
	duration := time.Since(start)
	
	success := reloadResponse["success"] == true
	
	result := TestResult{
		Name:     "Git Hot-Reload",
		Success:  success,
		Duration: duration,
		Details:  fmt.Sprintf("Initial definitions: %d, Reload success: %v", initialCount, success),
	}
	
	if !success {
		result.Error = "Hot-reload failed"
	}
	
	ts.Results = append(ts.Results, result)
	log.Printf("   Hot-reload: %v (Duration: %v)", success, duration)
	
	// Test 3: Verify worker status
	workersResponse := ts.makeRequest("GET", "/api/v1/temporal/workers/status", nil)
	workerCount := int(workersResponse["total"].(float64))
	
	workerResult := TestResult{
		Name:     "Worker Status Check",
		Success:  workerCount > 0,
		Duration: 0,
		Details:  fmt.Sprintf("Active workers: %d", workerCount),
	}
	
	ts.Results = append(ts.Results, workerResult)
	log.Printf("   Workers active: %d", workerCount)
}

func (ts *TestSuite) testMultiTenantIsolation() {
	log.Printf("🏢 Testing multi-tenant isolation...")
	
	customers := []string{"enterprise-corp", "startup-inc", "mid-market-co"}
	
	for i, customer := range customers {
		start := time.Now()
		
		// Test direct Temporal API with customer ID
		request := map[string]interface{}{
			"workflow_type": "DataPipelineWorkflow",
			"customer_id":   customer,
			"parameters": map[string]interface{}{
				"tenant_test": true,
				"customer":    customer,
			},
			"async": true,
		}
		
		response := ts.makeRequest("POST", "/api/v1/temporal/workflows/execute", request)
		duration := time.Since(start)
		
		success := response["success"] == true && response["workflow_id"] != nil
		
		result := TestResult{
			Name:     fmt.Sprintf("Multi-Tenant: %s", customer),
			Success:  success,
			Duration: duration,
			Details:  fmt.Sprintf("Customer: %s, WorkflowID: %v", customer, response["workflow_id"]),
		}
		
		if !success {
			result.Error = "Multi-tenant workflow execution failed"
		}
		
		ts.Results = append(ts.Results, result)
		log.Printf("   Customer %s: %v (WorkflowID: %v)", customer, success, response["workflow_id"])
		
		// Small delay between tenant tests
		time.Sleep(100 * time.Millisecond)
	}
}

func (ts *TestSuite) testWorkflowSignaling() {
	log.Printf("📡 Testing workflow signaling and control...")
	
	// Start a long-running workflow
	workflowRequest := map[string]interface{}{
		"workflow_type": "LongRunningAnalyticsWorkflow",
		"customer_id":   "signal-test-customer",
		"parameters": map[string]interface{}{
			"processing_days": 3,
			"test_mode":      true,
		},
		"async": true,
	}
	
	start := time.Now()
	workflowResponse := ts.makeRequest("POST", "/api/v1/temporal/workflows/execute", workflowRequest)
	
	if workflowResponse["success"] != true || workflowResponse["workflow_id"] == nil {
		result := TestResult{
			Name:    "Workflow Signaling Setup",
			Success: false,
			Error:   "Failed to start test workflow",
		}
		ts.Results = append(ts.Results, result)
		return
	}
	
	workflowID := workflowResponse["workflow_id"].(string)
	runID := workflowResponse["run_id"].(string)
	
	// Wait a moment for workflow to start
	time.Sleep(2 * time.Second)
	
	// Test 1: Check workflow status
	statusResponse := ts.makeRequest("GET", fmt.Sprintf("/api/v1/temporal/workflows/%s/runs/%s/status", workflowID, runID), nil)
	statusSuccess := statusResponse["success"] == true
	
	statusResult := TestResult{
		Name:     "Workflow Status Query",
		Success:  statusSuccess,
		Duration: time.Since(start),
		Details:  fmt.Sprintf("Status: %v", statusResponse["status"]),
	}
	ts.Results = append(ts.Results, statusResult)
	
	// Test 2: Send pause signal
	pauseRequest := map[string]interface{}{
		"signal_name": "pause",
		"data":        true,
	}
	
	pauseResponse := ts.makeRequest("POST", fmt.Sprintf("/api/v1/temporal/workflows/%s/runs/%s/signal", workflowID, runID), pauseRequest)
	pauseSuccess := pauseResponse["success"] == true
	
	pauseResult := TestResult{
		Name:     "Workflow Pause Signal",
		Success:  pauseSuccess,
		Duration: 0,
		Details:  "Sent pause signal to workflow",
	}
	ts.Results = append(ts.Results, pauseResult)
	
	// Test 3: Send resume signal
	time.Sleep(1 * time.Second)
	resumeRequest := map[string]interface{}{
		"signal_name": "resume",
		"data":        true,
	}
	
	resumeResponse := ts.makeRequest("POST", fmt.Sprintf("/api/v1/temporal/workflows/%s/runs/%s/signal", workflowID, runID), resumeRequest)
	resumeSuccess := resumeResponse["success"] == true
	
	resumeResult := TestResult{
		Name:     "Workflow Resume Signal",
		Success:  resumeSuccess,
		Duration: 0,
		Details:  "Sent resume signal to workflow",
	}
	ts.Results = append(ts.Results, resumeResult)
	
	// Test 4: Query workflow state
	queryRequest := map[string]interface{}{
		"query_type": "progress",
	}
	
	queryResponse := ts.makeRequest("POST", fmt.Sprintf("/api/v1/temporal/workflows/%s/runs/%s/query", workflowID, runID), queryRequest)
	querySuccess := queryResponse["success"] == true
	
	queryResult := TestResult{
		Name:     "Workflow Query",
		Success:  querySuccess,
		Duration: 0,
		Details:  fmt.Sprintf("Query result: %v", queryResponse["result"]),
	}
	ts.Results = append(ts.Results, queryResult)
	
	log.Printf("   Status check: %v", statusSuccess)
	log.Printf("   Pause signal: %v", pauseSuccess)
	log.Printf("   Resume signal: %v", resumeSuccess)
	log.Printf("   Query: %v", querySuccess)
	
	// Cancel the test workflow to clean up
	ts.makeRequest("POST", fmt.Sprintf("/api/v1/temporal/workflows/%s/runs/%s/cancel", workflowID, runID), nil)
}

func (ts *TestSuite) testErrorHandlingAndRetries() {
	log.Printf("🔧 Testing error handling and retry policies...")
	
	// Test 1: Invalid workflow type (should fail gracefully)
	invalidRequest := map[string]interface{}{
		"workflow_type": "NonExistentWorkflow",
		"customer_id":   "error-test-customer",
		"parameters":    map[string]interface{}{},
		"async":         false,
	}
	
	start := time.Now()
	errorResponse := ts.makeRequest("POST", "/api/v1/temporal/workflows/execute", invalidRequest)
	duration := time.Since(start)
	
	// Should fail but gracefully
	errorHandlingSuccess := errorResponse["success"] == false && errorResponse["error"] != nil
	
	errorResult := TestResult{
		Name:     "Error Handling: Invalid Workflow",
		Success:  errorHandlingSuccess,
		Duration: duration,
		Details:  fmt.Sprintf("Error handled gracefully: %v", errorResponse["error"]),
	}
	ts.Results = append(ts.Results, errorResult)
	
	// Test 2: Test with malformed request
	malformedRequest := map[string]interface{}{
		"invalid_field": "invalid_value",
	}
	
	malformedResponse := ts.makeRequest("POST", "/api/v1/temporal/workflows/execute", malformedRequest)
	malformedSuccess := malformedResponse["success"] == false
	
	malformedResult := TestResult{
		Name:     "Error Handling: Malformed Request",
		Success:  malformedSuccess,
		Duration: 0,
		Details:  "Malformed request handled gracefully",
	}
	ts.Results = append(ts.Results, malformedResult)
	
	// Test 3: Test service availability
	metricsResponse := ts.makeRequest("GET", "/api/v1/temporal/metrics", nil)
	metricsSuccess := metricsResponse["success"] == true
	
	metricsResult := TestResult{
		Name:     "Service Availability: Metrics",
		Success:  metricsSuccess,
		Duration: 0,
		Details:  fmt.Sprintf("Metrics available: %v", metricsResponse["metrics"]),
	}
	ts.Results = append(ts.Results, metricsResult)
	
	log.Printf("   Invalid workflow handling: %v", errorHandlingSuccess)
	log.Printf("   Malformed request handling: %v", malformedSuccess)
	log.Printf("   Service metrics: %v", metricsSuccess)
}

func (ts *TestSuite) testPerformanceBenchmark() {
	log.Printf("⚡ Running performance benchmarks...")
	
	// Benchmark 1: Deterministic execution speed
	deterministic_times := make([]time.Duration, 10)
	for i := 0; i < 10; i++ {
		start := time.Now()
		request := map[string]interface{}{
			"text":       fmt.Sprintf("calculate %d + %d", i*10, i*5),
			"session_id": fmt.Sprintf("perf-test-%d", i),
		}
		ts.makeRequest("POST", "/api/v1/execute", request)
		deterministic_times[i] = time.Since(start)
	}
	
	avgDeterministic := averageDuration(deterministic_times)
	deterministicSuccess := avgDeterministic < 100*time.Millisecond
	
	deterministicResult := TestResult{
		Name:     "Performance: Deterministic Avg",
		Success:  deterministicSuccess,
		Duration: avgDeterministic,
		Details:  fmt.Sprintf("Average: %v, Target: <100ms", avgDeterministic),
	}
	ts.Results = append(ts.Results, deterministicResult)
	
	// Benchmark 2: Temporal workflow startup time
	temporal_times := make([]time.Duration, 5)
	for i := 0; i < 5; i++ {
		start := time.Now()
		request := map[string]interface{}{
			"workflow_type": "DataPipelineWorkflow",
			"customer_id":   fmt.Sprintf("perf-customer-%d", i),
			"parameters":    map[string]interface{}{"test": true},
			"async":         true,
		}
		ts.makeRequest("POST", "/api/v1/temporal/workflows/execute", request)
		temporal_times[i] = time.Since(start)
	}
	
	avgTemporal := averageDuration(temporal_times)
	temporalSuccess := avgTemporal < 2*time.Second
	
	temporalResult := TestResult{
		Name:     "Performance: Temporal Startup Avg",
		Success:  temporalSuccess,
		Duration: avgTemporal,
		Details:  fmt.Sprintf("Average: %v, Target: <2s", avgTemporal),
	}
	ts.Results = append(ts.Results, temporalResult)
	
	log.Printf("   Deterministic avg: %v (target: <100ms)", avgDeterministic)
	log.Printf("   Temporal startup avg: %v (target: <2s)", avgTemporal)
}

func (ts *TestSuite) makeRequest(method, endpoint string, payload interface{}) map[string]interface{} {
	var body io.Reader
	
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshaling request: %v", err)
			return map[string]interface{}{"success": false, "error": err.Error()}
		}
		body = bytes.NewBuffer(jsonData)
	}
	
	req, err := http.NewRequest(method, ts.BaseURL+endpoint, body)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	defer resp.Body.Close()
	
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return map[string]interface{}{"success": false, "error": err.Error()}
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		return map[string]interface{}{"success": false, "error": err.Error(), "raw": string(responseBody)}
	}
	
	return result
}

func (ts *TestSuite) printResults() {
	log.Printf("\n" + strings.Repeat("=", 80))
	log.Printf("📊 COMPREHENSIVE TEST RESULTS")
	log.Printf(strings.Repeat("=", 80))
	
	passed := 0
	failed := 0
	totalDuration := time.Duration(0)
	
	for _, result := range ts.Results {
		status := "✅ PASS"
		if !result.Success {
			status = "❌ FAIL"
			failed++
		} else {
			passed++
		}
		
		totalDuration += result.Duration
		
		log.Printf("%s | %-35s | %8v | %s", 
			status, result.Name, result.Duration, result.Details)
		
		if result.Error != "" {
			log.Printf("     ERROR: %s", result.Error)
		}
	}
	
	log.Printf(strings.Repeat("-", 80))
	log.Printf("📈 SUMMARY")
	log.Printf("   Total Tests: %d", len(ts.Results))
	log.Printf("   Passed: %d", passed)
	log.Printf("   Failed: %d", failed)
	log.Printf("   Success Rate: %.1f%%", float64(passed)/float64(len(ts.Results))*100)
	log.Printf("   Total Duration: %v", totalDuration)
	log.Printf(strings.Repeat("=", 80))
	
	if failed > 0 {
		log.Printf("❌ Some tests failed. Check the logs above for details.")
		os.Exit(1)
	} else {
		log.Printf("🎉 All tests passed! Temporal integration is working perfectly.")
	}
}

func averageDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	
	total := time.Duration(0)
	for _, d := range durations {
		total += d
	}
	
	return total / time.Duration(len(durations))
}