package schema

import (
	"github.com/riposo/riposo/pkg/plugin"
	"github.com/riposo/riposo/pkg/rule"
	"github.com/riposo/schema/internal"
	"gopkg.in/yaml.v3"
)

func init() {
	rule.Register("schema", func(n *yaml.Node) (rule.Rule, error) {
		var cfg internal.Config
		if err := n.Decode(&cfg); err != nil {
			return nil, err
		}
		return internal.New(&cfg)
	})

	plugin.Register(plugin.New(
		"schema",
		map[string]interface{}{
			"description": "Uses a JSON schema to validate records.",
			"url":         "https://github.com/riposo/schema",
		},
		nil,
		nil,
	))
}
