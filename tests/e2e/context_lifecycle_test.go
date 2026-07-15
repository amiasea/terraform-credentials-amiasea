package e2e

import (
	"os"
	"testing"
)

func runContextLifecycleTests(
	t *testing.T,
) {
	t.Helper()

	defer purgeTfcred(t)

	runTfcred(
		t,
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
		"Token securely vaulted in Windows Credential Manager.",
	)

	assertContains(
		t,
		output,
		"Context 'platform' configured successfully.",
	)

	assertContains(
		t,
		output,
		"Current directory bound to context 'platform'",
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
		"status",
	)

	assertContains(
		t,
		output,
		"context=platform",
	)

	assertContains(
		t,
		output,
		"type=org",
	)

	assertContains(
		t,
		output,
		"org=acme-corp",
	)

	output = runTfcred(
		t,
		"list",
	)

	assertContains(
		t,
		output,
		"platform",
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

	assertContains(
		t,
		output,
		"vault_key=tfcred:domain:app_terraform_io:org:acme_corp",
	)

	assertNotContains(
		t,
		output,
		"token_value",
	)

	assertNotContains(
		t,
		output,
		"has_vaulted_token",
	)
}

func runCredentialsHelperBasic(
	t *testing.T,
) {
	t.Helper()

	defer purgeTfcred(t)

	runTfcred(
		t,
		"init",
		"--domain",
		"app.terraform.io",
	)

	token := os.Getenv("TF_TOKEN_ORG")
	if token == "" {
		t.Fatal("TF_TOKEN_ORG is not set")
	}

	output := runTfcred(
		t,
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

	assertContains(
		t,
		output,
		"Token securely vaulted in Windows Credential Manager.",
	)

	assertContains(
		t,
		output,
		"Context 'platform' configured successfully.",
	)

	assertContains(
		t,
		output,
		"Current directory bound to context 'platform'",
	)

	runTerraform(
		t,
		"init",
	)
}
