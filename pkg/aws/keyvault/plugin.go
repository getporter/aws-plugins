package keyvault

import (
	"os"

	"get.porter.sh/plugin/aws/pkg/aws/awsconfig"
	"get.porter.sh/porter/pkg/portercontext"
	"get.porter.sh/porter/pkg/secrets"
	"get.porter.sh/porter/pkg/secrets/plugins"
	"get.porter.sh/porter/pkg/secrets/pluginstore"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

const PluginInterface = plugins.PluginInterface + ".aws.secretsmanager"

var _ plugins.SecretsProtocol = &Plugin{}

// Plugin is the plugin wrapper for accessing secrets from AWS Secrets Manager.
type Plugin struct {
	secrets.Store
}

func NewPlugin(c *portercontext.Context, cfg awsconfig.Config) plugin.Plugin {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       PluginInterface,
		Output:     os.Stderr,
		Level:      hclog.Debug,
		JSONFormat: true,
	})

	return pluginstore.NewPlugin(c, NewStore(cfg, logger))
}
