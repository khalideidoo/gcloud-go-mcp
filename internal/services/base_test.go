package services

import (
	"testing"
	"time"

	"gcloud-go-mcp/internal/config"
)

func TestNewBaseService(t *testing.T) {
	cfg := &config.Config{
		Project:        "test-project",
		Region:         "us-central1",
		Zone:           "us-central1-a",
		GCloudPath:     "gcloud",
		CommandTimeout: 5 * time.Minute,
	}

	base := NewBaseService(cfg)

	if base == nil {
		t.Fatal("expected non-nil BaseService")
	}
	if base.Executor == nil {
		t.Error("expected non-nil Executor")
	}
	if base.Config != cfg {
		t.Error("expected Config to match input")
	}
}

func TestToolResult(t *testing.T) {
	result := ToolResult("test message")

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.IsError {
		t.Error("expected IsError to be false")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}
}

func TestToolError(t *testing.T) {
	err := &testError{msg: "test error"}
	result := ToolError(err)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if !result.IsError {
		t.Error("expected IsError to be true")
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}
}

func TestGetRequiredString_Success(t *testing.T) {
	args := map[string]any{
		"name": "test-value",
	}

	val, err := GetRequiredString(args, "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "test-value" {
		t.Errorf("expected 'test-value', got %q", val)
	}
}

func TestGetRequiredString_Missing(t *testing.T) {
	args := map[string]any{}

	_, err := GetRequiredString(args, "name")
	if err == nil {
		t.Error("expected error for missing parameter")
	}
}

func TestGetRequiredString_WrongType(t *testing.T) {
	args := map[string]any{
		"name": 123,
	}

	_, err := GetRequiredString(args, "name")
	if err == nil {
		t.Error("expected error for wrong type")
	}
}

func TestGetRequiredString_Empty(t *testing.T) {
	args := map[string]any{
		"name": "",
	}

	_, err := GetRequiredString(args, "name")
	if err == nil {
		t.Error("expected error for empty string")
	}
}

func TestGetOptionalString_Present(t *testing.T) {
	args := map[string]any{
		"region": "us-west1",
	}

	val := GetOptionalString(args, "region", "us-central1")
	if val != "us-west1" {
		t.Errorf("expected 'us-west1', got %q", val)
	}
}

func TestGetOptionalString_Missing(t *testing.T) {
	args := map[string]any{}

	val := GetOptionalString(args, "region", "us-central1")
	if val != "us-central1" {
		t.Errorf("expected default 'us-central1', got %q", val)
	}
}

func TestGetOptionalString_WrongType(t *testing.T) {
	args := map[string]any{
		"region": 123,
	}

	val := GetOptionalString(args, "region", "default")
	if val != "default" {
		t.Errorf("expected default for wrong type, got %q", val)
	}
}

func TestGetOptionalInt_Present(t *testing.T) {
	args := map[string]any{
		"limit": float64(100), // JSON numbers are float64
	}

	val := GetOptionalInt(args, "limit", 50)
	if val != 100 {
		t.Errorf("expected 100, got %d", val)
	}
}

func TestGetOptionalInt_Missing(t *testing.T) {
	args := map[string]any{}

	val := GetOptionalInt(args, "limit", 50)
	if val != 50 {
		t.Errorf("expected default 50, got %d", val)
	}
}

func TestGetOptionalInt_WrongType(t *testing.T) {
	args := map[string]any{
		"limit": "not a number",
	}

	val := GetOptionalInt(args, "limit", 50)
	if val != 50 {
		t.Errorf("expected default for wrong type, got %d", val)
	}
}

func TestGetOptionalInt_FromIntValue(t *testing.T) {
	// This tests the case where an int is passed (though JSON typically gives float64)
	args := map[string]any{
		"limit": 100, // actual int
	}

	// Should return default since it expects float64
	val := GetOptionalInt(args, "limit", 50)
	if val != 50 {
		t.Errorf("expected default 50 for int type (not float64), got %d", val)
	}
}

func TestGetOptionalBool_True(t *testing.T) {
	args := map[string]any{
		"verbose": true,
	}

	val := GetOptionalBool(args, "verbose", false)
	if !val {
		t.Error("expected true")
	}
}

func TestGetOptionalBool_False(t *testing.T) {
	args := map[string]any{
		"verbose": false,
	}

	val := GetOptionalBool(args, "verbose", true)
	if val {
		t.Error("expected false")
	}
}

func TestGetOptionalBool_Missing(t *testing.T) {
	args := map[string]any{}

	val := GetOptionalBool(args, "verbose", true)
	if !val {
		t.Error("expected default true")
	}
}

func TestGetOptionalBool_WrongType(t *testing.T) {
	args := map[string]any{
		"verbose": "true", // string, not bool
	}

	val := GetOptionalBool(args, "verbose", false)
	if val {
		t.Error("expected default false for wrong type")
	}
}

func TestGetOptionalStringArray_Present(t *testing.T) {
	args := map[string]any{
		"tags": []any{"tag1", "tag2", "tag3"},
	}

	val := GetOptionalStringArray(args, "tags")
	if len(val) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(val))
	}
	if val[0] != "tag1" || val[1] != "tag2" || val[2] != "tag3" {
		t.Errorf("unexpected tags: %v", val)
	}
}

func TestGetOptionalStringArray_Empty(t *testing.T) {
	args := map[string]any{
		"tags": []any{},
	}

	val := GetOptionalStringArray(args, "tags")
	if val == nil {
		t.Error("expected empty slice, got nil")
	}
	if len(val) != 0 {
		t.Errorf("expected 0 tags, got %d", len(val))
	}
}

func TestGetOptionalStringArray_Missing(t *testing.T) {
	args := map[string]any{}

	val := GetOptionalStringArray(args, "tags")
	if val != nil {
		t.Errorf("expected nil for missing array, got %v", val)
	}
}

func TestGetOptionalStringArray_WrongType(t *testing.T) {
	args := map[string]any{
		"tags": "not an array",
	}

	val := GetOptionalStringArray(args, "tags")
	if val != nil {
		t.Errorf("expected nil for wrong type, got %v", val)
	}
}

func TestGetOptionalStringArray_MixedTypes(t *testing.T) {
	args := map[string]any{
		"tags": []any{"tag1", 123, "tag2", true},
	}

	val := GetOptionalStringArray(args, "tags")
	// Should only include string values
	if len(val) != 2 {
		t.Errorf("expected 2 string tags, got %d: %v", len(val), val)
	}
}

func TestGetOptionalStringMap_Present(t *testing.T) {
	args := map[string]any{
		"labels": map[string]any{
			"env":  "production",
			"team": "backend",
		},
	}

	val := GetOptionalStringMap(args, "labels")
	if len(val) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(val))
	}
	if val["env"] != "production" {
		t.Errorf("expected env='production', got %q", val["env"])
	}
	if val["team"] != "backend" {
		t.Errorf("expected team='backend', got %q", val["team"])
	}
}

func TestGetOptionalStringMap_Empty(t *testing.T) {
	args := map[string]any{
		"labels": map[string]any{},
	}

	val := GetOptionalStringMap(args, "labels")
	if val == nil {
		t.Error("expected empty map, got nil")
	}
	if len(val) != 0 {
		t.Errorf("expected 0 labels, got %d", len(val))
	}
}

func TestGetOptionalStringMap_Missing(t *testing.T) {
	args := map[string]any{}

	val := GetOptionalStringMap(args, "labels")
	if val != nil {
		t.Errorf("expected nil for missing map, got %v", val)
	}
}

func TestGetOptionalStringMap_WrongType(t *testing.T) {
	args := map[string]any{
		"labels": "not a map",
	}

	val := GetOptionalStringMap(args, "labels")
	if val != nil {
		t.Errorf("expected nil for wrong type, got %v", val)
	}
}

func TestGetOptionalStringMap_NonStringValues(t *testing.T) {
	args := map[string]any{
		"labels": map[string]any{
			"count":   float64(42),
			"enabled": true,
			"name":    "test",
		},
	}

	val := GetOptionalStringMap(args, "labels")
	if len(val) != 3 {
		t.Fatalf("expected 3 labels, got %d", len(val))
	}
	// Non-string values should be converted to string representation
	if val["count"] != "42" {
		t.Errorf("expected count='42', got %q", val["count"])
	}
	if val["enabled"] != "true" {
		t.Errorf("expected enabled='true', got %q", val["enabled"])
	}
	if val["name"] != "test" {
		t.Errorf("expected name='test', got %q", val["name"])
	}
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

// Benchmark tests for parameter extraction
func BenchmarkGetRequiredString(b *testing.B) {
	args := map[string]any{
		"name": "test-value",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetRequiredString(args, "name")
	}
}

func BenchmarkGetOptionalString(b *testing.B) {
	args := map[string]any{
		"region": "us-west1",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetOptionalString(args, "region", "default")
	}
}

func BenchmarkGetOptionalInt(b *testing.B) {
	args := map[string]any{
		"limit": float64(100),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetOptionalInt(args, "limit", 50)
	}
}

func BenchmarkGetOptionalStringArray(b *testing.B) {
	args := map[string]any{
		"tags": []any{"tag1", "tag2", "tag3"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetOptionalStringArray(args, "tags")
	}
}

func BenchmarkGetOptionalStringMap(b *testing.B) {
	args := map[string]any{
		"labels": map[string]any{
			"env":  "production",
			"team": "backend",
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetOptionalStringMap(args, "labels")
	}
}
