package config

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestGetConfigValue_AddressDefaultAndCaseInsensitive(t *testing.T) {
	got := GetConfigValue("address")
	if got != ":8080" {
		t.Fatalf("expected :8080 for key 'address', got %q", got)
	}

	// Case-insensitive key
	got = GetConfigValue("ADDRESS")
	if got != ":8080" {
		t.Fatalf("expected :8080 for key 'ADDRESS', got %q", got)
	}
}

func TestGetConfigValue_ReturnsEnvVarsWhenSet(t *testing.T) {
	// Use t.Setenv to isolate environment for each variable
	t.Setenv("NEO4J_URL", "neo4j://localhost:7687")
	t.Setenv("NEO4J_USERNAME", "neo4j")
	t.Setenv("NEO4J_PASSWORD", "pass")

	if got := GetConfigValue("neo4j_url"); got != "neo4j://localhost:7687" {
		t.Fatalf("unexpected neo4j_url: got %q", got)
	}
	// Mixed case to ensure case-insensitive behavior
	if got := GetConfigValue("Neo4j_Username"); got != "neo4j" {
		t.Fatalf("unexpected neo4j_username: got %q", got)
	}
	if got := GetConfigValue("NEO4J_PASSWORD"); got != "pass" {
		t.Fatalf("unexpected neo4j_password: got %q", got)
	}
}

func TestGetConfigValue_MissingEnvVar_LogsAndReturnsEmpty(t *testing.T) {
	vars := []string{"NEO4J_URL", "NEO4J_USERNAME", "NEO4J_PASSWORD"}
	keys := []string{"neo4j_url", "neo4j_username", "neo4j_password"}

	for i, env := range vars {
		// Save current value and ensure the variable is truly unset (LookupEnv returns found=false)
		prev, hadPrev := os.LookupEnv(env)
		_ = os.Unsetenv(env)

		// Capture logs
		buf := &bytes.Buffer{}
		prevOut := log.Writer()
		prevFlags := log.Flags()
		prevPrefix := log.Prefix()
		log.SetOutput(buf)
		log.SetFlags(0)
		log.SetPrefix("")

		// Restore environment and logger after each iteration
		t.Cleanup(func() {
			if hadPrev {
				_ = os.Setenv(env, prev)
			} else {
				_ = os.Unsetenv(env)
			}
			log.SetOutput(prevOut)
			log.SetFlags(prevFlags)
			log.SetPrefix(prevPrefix)
		})

		// Call the function and assert empty string
		if got := GetConfigValue(keys[i]); got != "" {
			t.Fatalf("expected empty string for %s when %s is missing, got %q", keys[i], env, got)
		}

		// Ensure the log contains the expected message substring
		logged := buf.String()
		expectedSub := "Environment variable " + env + " not found"
		if !strings.Contains(logged, expectedSub) {
			t.Fatalf("expected log to contain %q, got %q", expectedSub, logged)
		}
	}
}

func TestGetConfigValue_UnknownKey_ReturnsEmpty(t *testing.T) {
	if got := GetConfigValue("does_not_exist"); got != "" {
		t.Fatalf("expected empty string for unknown key, got %q", got)
	}
}
