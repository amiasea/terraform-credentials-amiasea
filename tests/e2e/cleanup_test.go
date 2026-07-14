package e2e

import "testing"

func runCleanupTests(
	t *testing.T,
) {
	t.Helper()

	runTfcred(
		t,
		"init",
		"--domain",
		"app.terraform.io",
	)

	runTfcred(
		t,
		"add",
		"--context",
		"cleanup",
		"--org",
		"cleanup-org",
		"--token-type",
		"org",
		"--token",
		dummyToken,
		"--switch",
	)

	output := runTfcred(
		t,
		"remove",
		"cleanup",
	)

	assertContains(
		t,
		output,
		"Context 'cleanup' and associated bindings removed successfully.",
	)

	output = runTfcred(
		t,
		"current",
	)

	assertContains(
		t,
		output,
		"Current directory is not bound to any context.",
	)

	output = runTfcred(
		t,
		"purge",
		"--force",
	)

	assertContains(
		t,
		output,
		"Wiping all data...",
	)

	assertContains(
		t,
		output,
		"SUCCESS: System entirely purged.",
	)
}
