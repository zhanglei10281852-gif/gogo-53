package rail

import (
	"testing"
	"time"
)

func TestApplySimulationDoesNotDuplicateMigrationsOnRetry(t *testing.T) {
	manifest := testManifest()
	manifest.Components[0].Health = []HealthCriterion{{Name: "errors", Metric: "error_rate", Operator: "<", Threshold: 0.1, Required: true}}
	now := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	hash, err := ManifestHash(manifest)
	if err != nil {
		t.Fatal(err)
	}
	state := InitializeState(manifest, hash, now)
	state.Health = []HealthSample{{Environment: "staging", Component: "database", Metric: "error_rate", Value: 0.9, ObservedAt: now}}
	artifacts := []ArtifactResult{{Component: "database", Valid: true}, {Component: "api", Valid: true}}
	plan, err := BuildPlan(manifest, state, artifacts, now)
	if err != nil {
		t.Fatal(err)
	}
	failed, err := ApplySimulation(plan, state, "staging", now.Add(time.Minute))
	if err == nil {
		t.Fatal("expected the unhealthy first attempt to fail")
	}
	failed.Health = []HealthSample{{Environment: "staging", Component: "database", Metric: "error_rate", Value: 0.01, ObservedAt: now.Add(2 * time.Minute)}}
	retried, err := ApplySimulation(plan, failed, "staging", now.Add(3*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	migrations := retried.Environments["staging"].Components["database"].Migrations
	seen := map[string]bool{}
	for _, id := range migrations {
		if seen[id] {
			t.Fatalf("retried apply duplicated migration records: %v", migrations)
		}
		seen[id] = true
	}
	if len(migrations) != 1 || migrations[0] != "schema" {
		t.Fatalf("migrations = %v, want exactly [schema]", migrations)
	}
}
