package aws

import (
	"bufio"
	"encoding/json"
	"io"

	"get.porter.sh/plugin/aws/pkg/aws/awsconfig"
	"get.porter.sh/porter/pkg/runtime"
	"github.com/pkg/errors"
)

type Plugin struct {
	runtime.RuntimeConfig
	awsconfig.Config
}

// New AWS plugin client, initialized with useful defaults.
func New() *Plugin {
	p := &Plugin{
		RuntimeConfig: runtime.NewConfig(),
	}

	// Tell porter that we are running inside a plugin
	p.IsInternalPlugin = true
	p.InternalPluginKey = "porter.plugins.aws"

	return p
}

func (p *Plugin) LoadConfig() error {
	reader := bufio.NewReader(p.In)
	b, err := io.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "could not read stdin")
	}

	if len(b) == 0 {
		return nil
	}

	err = json.Unmarshal(b, &p.Config)
	if err != nil {
		return errors.Wrapf(err, "error unmarshaling stdin %q as aws.Config", string(b))
	}

	return nil
}
