package bricked

import (
	"fmt"
	"os/exec"
	"strings"
)

func azCli(arg ...string) string {
	fmt.Printf("[DEBUG] Executing command: %s\n", "az "+strings.Join(arg, " "))
	out, err := exec.Command("az", arg...).CombinedOutput()
	if err != nil {
		panic(err)
	}
	sout := strings.TrimSpace(string(out))
	fmt.Printf("[DEBUG] Command output: %s\n", sout)
	return sout
}

// Potential helpful error message:
//
//	Use `az login` or `az login --tenant <tenant-id>` to login.
//	Might also need to select the correct subscription using `az account set --subscription <name|id>`
func AzLogin(tenant string, subscription string) {
	azCli("login", "--tenant", tenant)
	azCli("account", "set", "--subscription", subscription)
}

func getEntraToken() string {
	return azCli(
		"account", "get-access-token",
		"--resource", "2ff814a6-3304-4ab8-85cb-cd0e6f879c1d",
		"--query", "accessToken",
		"--output", "tsv",
	)
}

func findWorkspaceUrl(workspaceName, resourceGroup string) string {
	return "https://" + azCli(
		"databricks", "workspace", "show",
		"--name", workspaceName,
		"--resource-group", resourceGroup,
		"--query", "workspaceUrl",
		"--output", "tsv",
	)
}
