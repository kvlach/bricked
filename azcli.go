package bricked

func azCli(arguments ...string) {
	// e.g. azCli("databricks", "workspace", "list")
	panic("TODO: Implement")
}

func azLogin() {
	// Potential helpful error message:
	//   Use `az login` or `az login --tenant <tenant-id>` to login.
	//   Might also need to select the correct subscription using `az account set --subscription <name|id>`
	panic("TODO: Implement")
}

func findWorkspaceUrl(workspaceName string) string {
	panic("TODO: Implement")
}

func getEntraToken() {
	// az account get-access-token \
	// 	--resource 2ff814a6-3304-4ab8-85cb-cd0e6f879c1d \
	// 	--query 'accessToken' \
	// 	--output tsv

	// azCli(
	// 	"account", "get-access-token",
	// 	"--resource", "2ff814a6-3304-4ab8-85cb-cd0e6f879c1d",
	// 	"--query", "accessToken",
	// 	"--output", "tsv",
	// )

	panic("TODO: Implement")
}
