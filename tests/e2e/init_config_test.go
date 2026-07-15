package e2e

import "testing"

func runConfigurationTests(
	t *testing.T,
) {
	t.Helper()

	defer purgeTfcred(t)

	output := runTfcred(
		t,
		"init",
		"--domain",
		"app.terraform.io",
	)

	assertContains(
		t,
		output,
		"[tfcred] ✅ Initialization completed successfully.",
	)

	output = runTfcred(
		t,
		"config",
	)

	assertContains(
		t,
		output,
		"default_domain=app.terraform.io",
	)

	output = runTfcred(
		t,
		"config",
		"--default-domain",
		"app.terraform.io",
	)

	assertContains(
		t,
		output,
		"[tfcred] default domain set to app.terraform.io",
	)

	output = runTfcred(
		t,
		"config",
	)

	assertContains(
		t,
		output,
		"default_domain=app.terraform.io",
	)
}
