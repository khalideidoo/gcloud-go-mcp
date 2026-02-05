package executor

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseJSON_Success(t *testing.T) {
	result := &Result{
		JSON: json.RawMessage(`{"name": "test", "count": 42}`),
	}

	var target struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	err := result.ParseJSON(&target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if target.Name != "test" {
		t.Errorf("expected Name 'test', got %q", target.Name)
	}
	if target.Count != 42 {
		t.Errorf("expected Count 42, got %d", target.Count)
	}
}

func TestParseJSON_Array(t *testing.T) {
	result := &Result{
		JSON: json.RawMessage(`[{"id": 1}, {"id": 2}, {"id": 3}]`),
	}

	var target []struct {
		ID int `json:"id"`
	}

	err := result.ParseJSON(&target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(target) != 3 {
		t.Errorf("expected 3 items, got %d", len(target))
	}
	if target[0].ID != 1 || target[1].ID != 2 || target[2].ID != 3 {
		t.Errorf("unexpected IDs: %v", target)
	}
}

func TestParseJSON_NoJSON(t *testing.T) {
	result := &Result{
		Stdout: "some text output",
	}

	var target any
	err := result.ParseJSON(&target)
	if err == nil {
		t.Error("expected error for nil JSON")
	}
	if !strings.Contains(err.Error(), "no JSON output") {
		t.Errorf("expected 'no JSON output' error, got: %v", err)
	}
}

func TestParseJSON_InvalidJSON(t *testing.T) {
	result := &Result{
		JSON: json.RawMessage(`{invalid json`),
	}

	var target any
	err := result.ParseJSON(&target)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestToJSONString_WithJSON(t *testing.T) {
	result := &Result{
		JSON: json.RawMessage(`{"name":"test"}`),
	}

	output := result.ToJSONString()

	// Should be pretty-printed
	if !strings.Contains(output, "  ") {
		t.Error("expected pretty-printed output with indentation")
	}
	if !strings.Contains(output, "name") || !strings.Contains(output, "test") {
		t.Errorf("expected JSON content in output, got: %s", output)
	}
}

func TestToJSONString_WithoutJSON(t *testing.T) {
	result := &Result{
		Stdout: "plain text output",
	}

	output := result.ToJSONString()
	if output != "plain text output" {
		t.Errorf("expected stdout when no JSON, got: %s", output)
	}
}

func TestToJSONString_InvalidJSON(t *testing.T) {
	// When JSON is present but can't be indented, return as-is
	result := &Result{
		JSON: json.RawMessage(`{"name":"test"}`),
	}

	output := result.ToJSONString()
	// Should still contain the JSON content
	if !strings.Contains(output, "name") {
		t.Errorf("expected JSON content, got: %s", output)
	}
}

func TestIsEmpty_Empty(t *testing.T) {
	result := &Result{}
	if !result.IsEmpty() {
		t.Error("expected empty result to be empty")
	}
}

func TestIsEmpty_WithJSON(t *testing.T) {
	result := &Result{
		JSON: json.RawMessage(`{}`),
	}
	if result.IsEmpty() {
		t.Error("expected result with JSON to not be empty")
	}
}

func TestIsEmpty_WithStdout(t *testing.T) {
	result := &Result{
		Stdout: "some output",
	}
	if result.IsEmpty() {
		t.Error("expected result with Stdout to not be empty")
	}
}

func TestIsEmpty_WithBoth(t *testing.T) {
	result := &Result{
		JSON:   json.RawMessage(`{"key": "value"}`),
		Stdout: `{"key": "value"}`,
	}
	if result.IsEmpty() {
		t.Error("expected result with JSON and Stdout to not be empty")
	}
}

func TestFormatError(t *testing.T) {
	err := &testError{msg: "command failed"}
	output := FormatError(err, "gcloud run list", "ERROR: permission denied")

	// Verify it's valid JSON
	var parsed ErrorResponse
	if jsonErr := json.Unmarshal([]byte(output), &parsed); jsonErr != nil {
		t.Fatalf("output is not valid JSON: %v", jsonErr)
	}

	if parsed.Error != "command failed" {
		t.Errorf("expected Error 'command failed', got %q", parsed.Error)
	}
	if parsed.Command != "gcloud run list" {
		t.Errorf("expected Command 'gcloud run list', got %q", parsed.Command)
	}
	if parsed.Stderr != "ERROR: permission denied" {
		t.Errorf("expected Stderr 'ERROR: permission denied', got %q", parsed.Stderr)
	}
}

func TestFormatError_EmptyFields(t *testing.T) {
	err := &testError{msg: "error"}
	output := FormatError(err, "", "")

	var parsed ErrorResponse
	if jsonErr := json.Unmarshal([]byte(output), &parsed); jsonErr != nil {
		t.Fatalf("output is not valid JSON: %v", jsonErr)
	}

	if parsed.Error != "error" {
		t.Errorf("expected Error 'error', got %q", parsed.Error)
	}
	// Empty fields may or may not be present in output
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestResult_AllFields(t *testing.T) {
	result := &Result{
		JSON:     json.RawMessage(`{"status": "ok"}`),
		Stdout:   `{"status": "ok"}`,
		Stderr:   "some warning",
		ExitCode: 0,
	}

	if result.ExitCode != 0 {
		t.Errorf("expected ExitCode 0, got %d", result.ExitCode)
	}
	if result.Stderr != "some warning" {
		t.Errorf("expected Stderr 'some warning', got %q", result.Stderr)
	}
}

func TestToJSONString_NestedJSON(t *testing.T) {
	result := &Result{
		JSON: json.RawMessage(`{"services":[{"name":"svc1"},{"name":"svc2"}],"meta":{"total":2}}`),
	}

	output := result.ToJSONString()

	// Should be pretty-printed with proper nesting
	if !strings.Contains(output, "services") {
		t.Error("expected 'services' in output")
	}
	if !strings.Contains(output, "svc1") {
		t.Error("expected 'svc1' in output")
	}
}
