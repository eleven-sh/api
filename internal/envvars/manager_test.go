package envvars

import (
	"os"
	"testing"
)

func TestGetUndefinedEnvVar(t *testing.T) {
	defer func() { _ = recover() }()

	Get("undefined_env_var")

	t.Fatalf("expected panic, got nothing")
}

func TestGetEmptyEnvVar(t *testing.T) {
	defer func() { _ = recover() }()

	envVarName := "eleven_empty_env_var"
	err := os.Setenv(envVarName, "")

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	Get(EnvVarName(envVarName))

	t.Fatalf("expected panic, got nothing")
}

func TestGetDefinedEnvVar(t *testing.T) {
	envVarName := "eleven_env_var"
	envVarVal := "eleven_env_var_value"

	err := os.Setenv(envVarName, envVarVal)

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	returnedVal := Get(EnvVarName(envVarName))

	if returnedVal != envVarVal {
		t.Fatalf(
			"expected env var to equal '%s', got '%s'",
			envVarVal,
			returnedVal,
		)
	}
}

func TestEnsureWithUndefinedEnvVars(t *testing.T) {
	defer func() { _ = recover() }()

	Ensure("./testdata/.env.dist")

	t.Fatalf("expected panic, got nothing")
}

func TestEnsureWithEmptyEnvVars(t *testing.T) {
	defer func() { _ = recover() }()

	err := os.Setenv("ENV_VAR", "")

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	err = os.Setenv("ENV_VAR2", "")

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	Ensure("./testdata/.env.dist")

	t.Fatalf("expected panic, got nothing")
}

func TestEnsureWithDefinedEnvVars(t *testing.T) {
	err := os.Setenv("ENV_VAR", "env_var_val")

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	err = os.Setenv("ENV_VAR2", "env_var2_val")

	if err != nil {
		t.Fatalf("expected no error, got '%+v'", err)
	}

	Ensure("./testdata/.env.dist")
}
