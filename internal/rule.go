package internal

import (
	"context"

	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/rule"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader" // use http loader
)

type ruleSet []constraint

type constraint struct {
	schema *jsonschema.Schema
	globs  []string
}

// New initiates a rule from the config.
func New(cfg []*Config) (rule.Set, error) {
	rs := make(ruleSet, 0, len(cfg))
	for _, c := range cfg {
		if err := c.validate(); err != nil {
			return nil, err
		}

		schema, err := jsonschema.Compile(c.Schema)
		if err != nil {
			return nil, err
		}
		rs = append(rs, constraint{schema: schema, globs: c.Path})
	}
	return rs, nil
}

func (rs ruleSet) Init(ctx context.Context, rts *api.Routes, _ riposo.Helpers) error {
	rts.Callbacks(&callbacks{rs: rs})
	return nil
}

func (ruleSet) Close() error {
	return nil
}
