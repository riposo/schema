package schema

import (
	"github.com/riposo/riposo/pkg/plugin"
	"github.com/riposo/riposo/pkg/rule"
	"github.com/riposo/schema/internal"
	"gopkg.in/yaml.v3"
)

func init() {
	rule.Register("schema", func(nodes []*yaml.Node) (rule.Set, error) {
		cc := make([]*internal.Config, 0, len(nodes))
		for _, n := range nodes {
			c := new(internal.Config)
			if err := n.Decode(c); err != nil {
				return nil, err
			}
			cc = append(cc, c)
		}
		return internal.New(cc)
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
