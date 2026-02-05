package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ParseJSON parses the JSON result into a target struct.
func (r *Result) ParseJSON(target any) error {
	if r.JSON == nil {
		return fmt.Errorf("no JSON output available")
	}
	return json.Unmarshal(r.JSON, target)
}

// ToJSONString returns the JSON as a formatted string for MCP response.
func (r *Result) ToJSONString() string {
	if r.JSON != nil {
		// Pretty print for readability
		var pretty bytes.Buffer
		if err := json.Indent(&pretty, r.JSON, "", "  "); err == nil {
			return pretty.String()
		}
		return string(r.JSON)
	}
	return r.Stdout
}

// IsEmpty returns true if the result has no meaningful output.
func (r *Result) IsEmpty() bool {
	return r.JSON == nil && r.Stdout == ""
}

// ErrorResponse creates a standardized error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Command string `json:"command,omitempty"`
	Stderr  string `json:"stderr,omitempty"`
}

// FormatError creates a formatted error response.
func FormatError(err error, command string, stderr string) string {
	resp := ErrorResponse{
		Error:   err.Error(),
		Command: command,
		Stderr:  stderr,
	}
	b, _ := json.MarshalIndent(resp, "", "  ")
	return string(b)
}
