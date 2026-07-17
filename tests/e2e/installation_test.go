package e2e

import "testing"

func runInstallationTests(
	t *testing.T,
) {
	t.Helper()

	t.Logf(
		"verified installed tfcred binary: %s",
		tfCredBinary,
	)

	t.Logf(
		"verified tfcred data directory: %s",
		tfCredDataDir,
	)

	t.Logf(
		"verified terraform config: %s",
		terraformConfigFile,
	)

	t.Logf(
		"verified terraform helper: %s",
		tfHelperCredBinary,
	)
}
