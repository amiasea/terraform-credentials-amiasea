package e2e

import (
	"os"
	"testing"
)

func runContextLifecycleTests(
	t *testing.T,
) {
	t.Helper()

	defer purgeTfcredAfterTest(t)

	runTfcred(
		t,
		"",
		"init",
		"--domain",
		"app.terraform.io",
	)

	output := runTfcred(
		t,
		"add",
		"--context",
		"platform",
		"--org",
		"acme-corp",
		"--token-type",
		"org",
		"--token",
		dummyToken,
		"--switch",
	)

	assertContains(
		t,
		output,
		"Context 'platform' configured successfully.",
	)

	output = runTfcred(
		t,
		"current",
	)

	assertContains(
		t,
		output,
		"active_context=platform",
	)

	output = runTfcred(
		t,
		"context",
	)

	assertContains(
		t,
		output,
		"context=platform",
	)

	assertNotContains(
		t,
		output,
		"token_value",
	)
}

func runCredentialsHelperTests(
	t *testing.T,
) {
	t.Helper()

	defer purgeTfcredAfterTest(t)

	runTfcred(
		t,
		"",
		"init",
		"--domain",
		"app.terraform.io",
	)

	token := os.Getenv("TF_TOKEN_ORG")

	if token == "" {
		t.Fatal("TF_TOKEN_ORG is not set")
	}

	runTfcredAt(
		t,
		workspaceDir,
		"add",
		"--context",
		"platform",
		"--org",
		"amiasea",
		"--token-type",
		"org",
		"--token",
		token,
		"--switch",
	)

	runTerraform(
		t,
		"init",
	)
}
