package environment_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ronei-kunkel/environment"
)

type TestVars struct {
	ENVIRONMENT string `env:"APP_ENV"`
	DB_NAME     string
	SOME_KEY    string
}

func resetEnv(keys ...string) {
	for _, k := range keys {
		os.Unsetenv(k)
	}
}

func TestLoad_AllVarsSet(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	os.Setenv("DB_NAME", "main_db")
	os.Setenv("SOME_KEY", "value123")
	defer resetEnv("APP_ENV", "DB_NAME", "SOME_KEY")

	vars := environment.Load[TestVars]()

	if vars.ENVIRONMENT != "production" {
		t.Errorf("expected ENVIRONMENT=production, got %s", vars.ENVIRONMENT)
	}
	if vars.DB_NAME != "main_db" {
		t.Errorf("expected DB_NAME=main_db, got %s", vars.DB_NAME)
	}
	if vars.SOME_KEY != "value123" {
		t.Errorf("expected SOME_KEY=value123, got %s", vars.SOME_KEY)
	}
}

func TestLoad_MissingVars_Fatalln(t *testing.T) {
	if os.Getenv("GO_SUBPROCESS") == "1" {
		resetEnv("APP_ENV", "DB_NAME", "SOME_KEY")
		_ = environment.Load[TestVars]()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLoad_MissingVars_Fatalln")
	cmd.Env = append(os.Environ(), "GO_SUBPROCESS=1")
	out, err := cmd.CombinedOutput()

	if e, ok := err.(*exec.ExitError); !ok || e.ExitCode() == 0 {
		t.Fatalf("expected subprocess to exit with code != 0")
	}

	output := string(out)
	expectedMessages := []string{
		"has no `APP_ENV` environment variable defined",
		"has no `DB_NAME` environment variable defined",
		"has no `SOME_KEY` environment variable defined",
		"Aborting due to missing env vars",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("expected output to contain %q, got: %s", msg, output)
		}
	}
}

func TestLoad_WithDotEnvFile(t *testing.T) {
	envContent := `
APP_ENV=staging
DB_NAME=test_db
SOME_KEY=abc123
`
	f, err := os.CreateTemp("", "*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	f.WriteString(envContent)
	f.Close()

	vars := environment.Load[TestVars](f.Name())

	if vars.ENVIRONMENT != "staging" {
		t.Errorf("expected ENVIRONMENT=staging, got %s", vars.ENVIRONMENT)
	}
	if vars.DB_NAME != "test_db" {
		t.Errorf("expected DB_NAME=test_db, got %s", vars.DB_NAME)
	}
	if vars.SOME_KEY != "abc123" {
		t.Errorf("expected SOME_KEY=abc123, got %s", vars.SOME_KEY)
	}
}

func TestLoad_TagMapping(t *testing.T) {
	os.Setenv("APP_ENV", "production")
	os.Setenv("DB_NAME", "main_db")
	os.Setenv("SOME_KEY", "key123")
	defer resetEnv("APP_ENV", "DB_NAME", "SOME_KEY")

	vars := environment.Load[TestVars]()

	if vars.ENVIRONMENT != "production" {
		t.Errorf("tag env mapping failed, expected ENVIRONMENT=production, got %s", vars.ENVIRONMENT)
	}
}
