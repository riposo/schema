package internal

import (
	"context"
	"errors"

	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/rule"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
)

var errNoPath = errors.New("configuration must specify a path")

// Config contains the configuration for the rule.
type Config struct {
	Schema string
	Path   []string
}

// New initiates a rule from the config.
func New(c *Config) (rule.Rule, error) {
	schema, err := jsonschema.Compile(c.Schema)
	if err != nil {
		return nil, err
	}

	if len(c.Path) == 0 {
		return nil, errNoPath
	}

	return &loader{paths: c.Path, schema: schema}, nil
}

type loader struct {
	schema *jsonschema.Schema
	paths  []string
}

func (l *loader) Init(ctx context.Context, rts *api.Routes, _ riposo.Helpers) error {
	rts.Hook(l.paths, &hook{s: l.schema})
	return nil
}

func (*loader) Close() error {
	return nil
}
