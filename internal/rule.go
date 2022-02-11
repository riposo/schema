package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/rule"
	"gopkg.in/yaml.v3"

	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
)

var errNoPath = errors.New("configuration must specify at least one path pattern")

// Config contains the configuration for the rule.
type Config struct {
	Schema string
	Path   StringSlice
}

// StringSlice can unmarshal individual strings or array from YAML.
type StringSlice []string

func (s *StringSlice) UnmarshalYAML(n *yaml.Node) error {
	switch n.Kind {
	case yaml.ScalarNode:
		var v string
		if err := n.Decode(&v); err != nil {
			return err
		}
		*s = []string{v}
	default:
		var v []string
		if err := n.Decode(&v); err != nil {
			return err
		}
		*s = v
	}
	return nil
}

// --------------------------------------------------------------------

type loader struct {
	schema *jsonschema.Schema
	globs  []string
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

	for _, path := range c.Path {
		if !doublestar.ValidatePathPattern(path) {
			return nil, fmt.Errorf("configuration contains an invalid path pattern: %q", path)
		}
	}

	return &loader{globs: c.Path, schema: schema}, nil
}

func (l *loader) Init(ctx context.Context, rts *api.Routes, _ riposo.Helpers) error {
	rts.Callbacks(&callbacks{globs: l.globs, schema: l.schema})
	return nil
}

func (*loader) Close() error {
	return nil
}
