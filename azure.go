package bricked

import (
	"log/slog"
	"os/exec"
	"strings"
)

type azure struct {
	tenant       string // GUID
	subscription string // name or GUID
	token        string
}

func NewAzure(tenant, subscription string) *azure {
	return &azure{
		tenant:       tenant,
		subscription: subscription,
	}
}

func (az *azure) CLI(arg ...string) string {
	slog.Debug("executing", "command", "az "+strings.Join(arg, " "))
	out, err := exec.Command("az", arg...).CombinedOutput()
	sout := strings.TrimSpace(string(out))
	slog.Debug("command", "output", string(out))
	if err != nil {
		panic(err)
	}
	return sout
}

func (az *azure) Login() {
	az.CLI("login", "--tenant", az.tenant)
	az.CLI("account", "set", "--subscription", az.subscription)
}

func (az *azure) getEntraToken() string {
	if az.token != "" {
		return az.token
	}
	az.token = az.CLI(
		"account", "get-access-token",
		"--resource", "2ff814a6-3304-4ab8-85cb-cd0e6f879c1d",
		"--query", "accessToken",
		"--output", "tsv",
	)
	return az.token
}

func (az *azure) findWorkspaceUrl(workspaceName, resourceGroup string) string {
	return "https://" + az.CLI(
		"databricks", "workspace", "show",
		"--name", workspaceName,
		"--resource-group", resourceGroup,
		"--query", "workspaceUrl",
		"--output", "tsv",
	)
}

func (az *azure) NewDatabricks(workspaceName, resourceGroup, apiVersion string) *databricks {
	return &databricks{
		workspaceURL: az.findWorkspaceUrl(workspaceName, resourceGroup),
		apiVersion:   apiVersion,
		token:        az.getEntraToken(),
	}
}
