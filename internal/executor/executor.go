// Package executor provides a fluent API for executing gcloud CLI commands.
package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"gcloud-go-mcp/internal/config"
)

// Result represents the result of a gcloud command execution.
type Result struct {
	// JSON contains the parsed JSON output (if JSON format was used).
	JSON json.RawMessage

	// Stdout contains the raw standard output.
	Stdout string

	// Stderr contains the raw standard error.
	Stderr string

	// ExitCode contains the command exit code.
	ExitCode int
}

// Executor handles gcloud command execution.
type Executor struct {
	config *config.Config
}

// New creates a new gcloud executor.
func New(cfg *config.Config) *Executor {
	return &Executor{config: cfg}
}

// CommandBuilder provides a fluent interface for building gcloud commands.
type CommandBuilder struct {
	executor   *Executor
	components []string
	flags      map[string]string
	arrayFlags map[string][]string
	boolFlags  []string
	project    string
	region     string
	zone       string
	format     string
}

// Command starts building a new gcloud command.
func (e *Executor) Command(components ...string) *CommandBuilder {
	return &CommandBuilder{
		executor:   e,
		components: components,
		flags:      make(map[string]string),
		arrayFlags: make(map[string][]string),
		project:    e.config.Project,
		region:     e.config.Region,
		zone:       e.config.Zone,
		format:     "json",
	}
}

// WithProject sets the project for this command.
func (b *CommandBuilder) WithProject(project string) *CommandBuilder {
	if project != "" {
		b.project = project
	}
	return b
}

// WithRegion sets the region for this command.
func (b *CommandBuilder) WithRegion(region string) *CommandBuilder {
	if region != "" {
		b.region = region
	}
	return b
}

// WithZone sets the zone for this command.
func (b *CommandBuilder) WithZone(zone string) *CommandBuilder {
	if zone != "" {
		b.zone = zone
	}
	return b
}

// WithFlag adds a flag with a value.
func (b *CommandBuilder) WithFlag(name, value string) *CommandBuilder {
	if value != "" {
		b.flags[name] = value
	}
	return b
}

// WithArrayFlag adds a flag that can be specified multiple times.
func (b *CommandBuilder) WithArrayFlag(name, value string) *CommandBuilder {
	if value != "" {
		b.arrayFlags[name] = append(b.arrayFlags[name], value)
	}
	return b
}

// WithBoolFlag adds a boolean flag (no value).
func (b *CommandBuilder) WithBoolFlag(name string) *CommandBuilder {
	b.boolFlags = append(b.boolFlags, name)
	return b
}

// WithFormat sets the output format.
func (b *CommandBuilder) WithFormat(format string) *CommandBuilder {
	b.format = format
	return b
}

// WithTextFormat sets text output format (disables JSON parsing).
func (b *CommandBuilder) WithTextFormat() *CommandBuilder {
	b.format = ""
	return b
}

// Build constructs the full command arguments.
func (b *CommandBuilder) Build() []string {
	args := make([]string, 0, len(b.components)+len(b.flags)*2+len(b.boolFlags)+4)
	args = append(args, b.components...)

	// Add flags
	for name, value := range b.flags {
		args = append(args, fmt.Sprintf("--%s=%s", name, value))
	}

	// Add array flags
	for name, values := range b.arrayFlags {
		for _, value := range values {
			args = append(args, fmt.Sprintf("--%s=%s", name, value))
		}
	}

	// Add boolean flags
	for _, flag := range b.boolFlags {
		args = append(args, fmt.Sprintf("--%s", flag))
	}

	// Add project if set
	if b.project != "" {
		args = append(args, fmt.Sprintf("--project=%s", b.project))
	}

	// Add format if set
	if b.format != "" {
		args = append(args, fmt.Sprintf("--format=%s", b.format))
	}

	return args
}

// Execute runs the command and returns the result.
func (b *CommandBuilder) Execute(ctx context.Context) (*Result, error) {
	args := b.Build()

	ctx, cancel := context.WithTimeout(ctx, b.executor.config.CommandTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, b.executor.config.GCloudPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	result := &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
		return result, fmt.Errorf("gcloud command failed: %w\nstderr: %s", err, stderr.String())
	}

	// Parse JSON if format was JSON and output is not empty
	if b.format == "json" && stdout.Len() > 0 {
		trimmed := strings.TrimSpace(stdout.String())
		if trimmed != "" {
			result.JSON = json.RawMessage(trimmed)
		}
	}

	return result, nil
}

// ExecuteWithRegion runs the command with a region flag (for regional resources).
func (b *CommandBuilder) ExecuteWithRegion(ctx context.Context) (*Result, error) {
	if b.region != "" {
		b.WithFlag("region", b.region)
	}
	return b.Execute(ctx)
}

// ExecuteWithZone runs the command with a zone flag (for zonal resources).
func (b *CommandBuilder) ExecuteWithZone(ctx context.Context) (*Result, error) {
	if b.zone != "" {
		b.WithFlag("zone", b.zone)
	}
	return b.Execute(ctx)
}

// GetProject returns the current project setting.
func (b *CommandBuilder) GetProject() string {
	return b.project
}

// GetRegion returns the current region setting.
func (b *CommandBuilder) GetRegion() string {
	return b.region
}

// GetZone returns the current zone setting.
func (b *CommandBuilder) GetZone() string {
	return b.zone
}
