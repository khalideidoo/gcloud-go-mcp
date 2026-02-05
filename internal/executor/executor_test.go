package executor

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"gcloud-go-mcp/internal/config"
)

func newTestConfig() *config.Config {
	return &config.Config{
		Project:        "default-project",
		Region:         "us-central1",
		Zone:           "us-central1-a",
		GCloudPath:     "gcloud",
		CommandTimeout: 5 * time.Minute,
	}
}

func TestNew(t *testing.T) {
	cfg := newTestConfig()
	exec := New(cfg)

	if exec == nil {
		t.Fatal("expected non-nil executor")
	}
	if exec.config != cfg {
		t.Error("executor config mismatch")
	}
}

func TestCommand_Basic(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "list")

	if builder == nil {
		t.Fatal("expected non-nil builder")
	}
	if builder.executor != exec {
		t.Error("builder executor mismatch")
	}
	if !reflect.DeepEqual(builder.components, []string{"run", "services", "list"}) {
		t.Errorf("expected components [run services list], got %v", builder.components)
	}
}

func TestCommand_InheritsDefaults(t *testing.T) {
	cfg := newTestConfig()
	exec := New(cfg)
	builder := exec.Command("run", "services", "list")

	if builder.project != cfg.Project {
		t.Errorf("expected project %q, got %q", cfg.Project, builder.project)
	}
	if builder.region != cfg.Region {
		t.Errorf("expected region %q, got %q", cfg.Region, builder.region)
	}
	if builder.zone != cfg.Zone {
		t.Errorf("expected zone %q, got %q", cfg.Zone, builder.zone)
	}
	if builder.format != "json" {
		t.Errorf("expected format 'json', got %q", builder.format)
	}
}

func TestWithProject(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "list").
		WithProject("custom-project")

	if builder.project != "custom-project" {
		t.Errorf("expected project 'custom-project', got %q", builder.project)
	}
}

func TestWithProject_Empty(t *testing.T) {
	cfg := newTestConfig()
	exec := New(cfg)
	builder := exec.Command("run", "services", "list").
		WithProject("")

	// Should keep default when empty string passed
	if builder.project != cfg.Project {
		t.Errorf("expected project %q (default), got %q", cfg.Project, builder.project)
	}
}

func TestWithRegion(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "list").
		WithRegion("us-west1")

	if builder.region != "us-west1" {
		t.Errorf("expected region 'us-west1', got %q", builder.region)
	}
}

func TestWithRegion_Empty(t *testing.T) {
	cfg := newTestConfig()
	exec := New(cfg)
	builder := exec.Command("run", "services", "list").
		WithRegion("")

	if builder.region != cfg.Region {
		t.Errorf("expected region %q (default), got %q", cfg.Region, builder.region)
	}
}

func TestWithZone(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("compute", "instances", "list").
		WithZone("us-west1-b")

	if builder.zone != "us-west1-b" {
		t.Errorf("expected zone 'us-west1-b', got %q", builder.zone)
	}
}

func TestWithZone_Empty(t *testing.T) {
	cfg := newTestConfig()
	exec := New(cfg)
	builder := exec.Command("compute", "instances", "list").
		WithZone("")

	if builder.zone != cfg.Zone {
		t.Errorf("expected zone %q (default), got %q", cfg.Zone, builder.zone)
	}
}

func TestWithFlag(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "list").
		WithFlag("limit", "100")

	if builder.flags["limit"] != "100" {
		t.Errorf("expected flag 'limit'='100', got %q", builder.flags["limit"])
	}
}

func TestWithFlag_Empty(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "list").
		WithFlag("limit", "")

	if _, ok := builder.flags["limit"]; ok {
		t.Error("expected empty flag to not be added")
	}
}

func TestWithArrayFlag(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "deploy").
		WithArrayFlag("env", "KEY1=value1").
		WithArrayFlag("env", "KEY2=value2")

	expected := []string{"KEY1=value1", "KEY2=value2"}
	if !reflect.DeepEqual(builder.arrayFlags["env"], expected) {
		t.Errorf("expected arrayFlags['env']=%v, got %v", expected, builder.arrayFlags["env"])
	}
}

func TestWithArrayFlag_Empty(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "deploy").
		WithArrayFlag("env", "")

	if len(builder.arrayFlags["env"]) != 0 {
		t.Error("expected empty array flag to not be added")
	}
}

func TestWithBoolFlag(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "list").
		WithBoolFlag("quiet").
		WithBoolFlag("verbose")

	if !reflect.DeepEqual(builder.boolFlags, []string{"quiet", "verbose"}) {
		t.Errorf("expected boolFlags [quiet verbose], got %v", builder.boolFlags)
	}
}

func TestWithFormat(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "list").
		WithFormat("yaml")

	if builder.format != "yaml" {
		t.Errorf("expected format 'yaml', got %q", builder.format)
	}
}

func TestWithTextFormat(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("secrets", "versions", "access").
		WithTextFormat()

	if builder.format != "" {
		t.Errorf("expected empty format, got %q", builder.format)
	}
}

func TestGetProject(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "list").
		WithProject("my-project")

	if builder.GetProject() != "my-project" {
		t.Errorf("expected GetProject() 'my-project', got %q", builder.GetProject())
	}
}

func TestGetRegion(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "list").
		WithRegion("eu-west1")

	if builder.GetRegion() != "eu-west1" {
		t.Errorf("expected GetRegion() 'eu-west1', got %q", builder.GetRegion())
	}
}

func TestGetZone(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("compute", "instances", "list").
		WithZone("asia-east1-a")

	if builder.GetZone() != "asia-east1-a" {
		t.Errorf("expected GetZone() 'asia-east1-a', got %q", builder.GetZone())
	}
}

func TestBuild_Basic(t *testing.T) {
	cfg := &config.Config{
		Project:        "",
		Region:         "",
		Zone:           "",
		GCloudPath:     "gcloud",
		CommandTimeout: time.Minute,
	}
	exec := New(cfg)
	args := exec.Command("run", "services", "list").
		WithTextFormat().
		Build()

	expected := []string{"run", "services", "list"}
	if !reflect.DeepEqual(args, expected) {
		t.Errorf("expected args %v, got %v", expected, args)
	}
}

func TestBuild_WithProject(t *testing.T) {
	exec := New(newTestConfig())
	args := exec.Command("run", "services", "list").
		WithProject("my-project").
		Build()

	// Check that project flag is present
	found := false
	for _, arg := range args {
		if arg == "--project=my-project" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected --project=my-project in args, got %v", args)
	}
}

func TestBuild_WithFlags(t *testing.T) {
	cfg := &config.Config{
		Project:        "",
		GCloudPath:     "gcloud",
		CommandTimeout: time.Minute,
	}
	exec := New(cfg)
	args := exec.Command("run", "services", "list").
		WithFlag("limit", "50").
		WithFlag("filter", "status=READY").
		WithTextFormat().
		Build()

	// Check flags are present (order may vary due to map iteration)
	hasLimit := false
	hasFilter := false
	for _, arg := range args {
		if arg == "--limit=50" {
			hasLimit = true
		}
		if arg == "--filter=status=READY" {
			hasFilter = true
		}
	}
	if !hasLimit {
		t.Errorf("expected --limit=50 in args, got %v", args)
	}
	if !hasFilter {
		t.Errorf("expected --filter=status=READY in args, got %v", args)
	}
}

func TestBuild_WithBoolFlags(t *testing.T) {
	cfg := &config.Config{
		Project:        "",
		GCloudPath:     "gcloud",
		CommandTimeout: time.Minute,
	}
	exec := New(cfg)
	args := exec.Command("run", "services", "delete", "my-service").
		WithBoolFlag("quiet").
		WithTextFormat().
		Build()

	found := false
	for _, arg := range args {
		if arg == "--quiet" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected --quiet in args, got %v", args)
	}
}

func TestBuild_WithArrayFlags(t *testing.T) {
	cfg := &config.Config{
		Project:        "",
		GCloudPath:     "gcloud",
		CommandTimeout: time.Minute,
	}
	exec := New(cfg)
	args := exec.Command("run", "deploy", "my-service").
		WithArrayFlag("set-env-vars", "KEY1=val1").
		WithArrayFlag("set-env-vars", "KEY2=val2").
		WithTextFormat().
		Build()

	count := 0
	for _, arg := range args {
		if arg == "--set-env-vars=KEY1=val1" || arg == "--set-env-vars=KEY2=val2" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected both env vars in args, found %d, args: %v", count, args)
	}
}

func TestBuild_WithFormat(t *testing.T) {
	cfg := &config.Config{
		Project:        "",
		GCloudPath:     "gcloud",
		CommandTimeout: time.Minute,
	}
	exec := New(cfg)
	args := exec.Command("run", "services", "list").Build()

	found := false
	for _, arg := range args {
		if arg == "--format=json" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected --format=json in args, got %v", args)
	}
}

func TestBuild_CompleteCommand(t *testing.T) {
	exec := New(newTestConfig())
	builder := exec.Command("run", "services", "deploy", "my-service").
		WithProject("prod-project").
		WithRegion("us-east1").
		WithFlag("image", "gcr.io/my-project/my-image").
		WithFlag("memory", "512Mi").
		WithBoolFlag("allow-unauthenticated")

	args := builder.Build()

	// Verify components are first
	if args[0] != "run" || args[1] != "services" || args[2] != "deploy" || args[3] != "my-service" {
		t.Errorf("expected components first, got %v", args[:4])
	}

	// Verify all expected flags are present
	argsSet := make(map[string]bool)
	for _, arg := range args {
		argsSet[arg] = true
	}

	expected := []string{
		"--project=prod-project",
		"--image=gcr.io/my-project/my-image",
		"--memory=512Mi",
		"--allow-unauthenticated",
		"--format=json",
	}

	for _, exp := range expected {
		if !argsSet[exp] {
			t.Errorf("expected %q in args, got %v", exp, args)
		}
	}
}

func TestChaining(t *testing.T) {
	exec := New(newTestConfig())

	// Test method chaining returns the same builder
	builder := exec.Command("run", "services", "list")
	builder2 := builder.
		WithProject("p").
		WithRegion("r").
		WithZone("z").
		WithFlag("f", "v").
		WithArrayFlag("a", "v1").
		WithBoolFlag("b").
		WithFormat("yaml")

	if builder != builder2 {
		t.Error("expected chaining to return same builder instance")
	}
}

func TestBuild_Deterministic(t *testing.T) {
	exec := New(newTestConfig())

	// Build the same command multiple times
	for i := 0; i < 10; i++ {
		args := exec.Command("run", "services", "list").
			WithProject("test-project").
			Build()

		// Components should always be first
		if args[0] != "run" || args[1] != "services" || args[2] != "list" {
			t.Errorf("iteration %d: components not first, got %v", i, args)
		}
	}
}

// TestBuild_FlagsOrder tests that flags are consistently ordered
func TestBuild_FlagsOrder(t *testing.T) {
	cfg := &config.Config{
		GCloudPath:     "gcloud",
		CommandTimeout: time.Minute,
	}
	exec := New(cfg)

	// Build command multiple times to verify consistency
	var allResults [][]string
	for i := 0; i < 5; i++ {
		args := exec.Command("test").
			WithTextFormat().
			WithBoolFlag("flag1").
			WithBoolFlag("flag2").
			Build()
		allResults = append(allResults, args)
	}

	// Sort each result and compare
	for i := 1; i < len(allResults); i++ {
		sorted0 := make([]string, len(allResults[0]))
		sortedI := make([]string, len(allResults[i]))
		copy(sorted0, allResults[0])
		copy(sortedI, allResults[i])
		sort.Strings(sorted0)
		sort.Strings(sortedI)
		if !reflect.DeepEqual(sorted0, sortedI) {
			t.Error("flag results should be consistent")
		}
	}
}
