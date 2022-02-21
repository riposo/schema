package internal

import (
	"errors"
	"fmt"

	"github.com/bmatcuk/doublestar/v4"
	"gopkg.in/yaml.v3"
)

var errNoPath = errors.New("configuration must specify at least one path pattern")

// Config contains the configuration for the rule.
type Config struct {
	Schema string
	Path   StringSlice
}

func (c *Config) validate() error {
	if len(c.Path) == 0 {
		return errNoPath
	}

	for _, path := range c.Path {
		if !doublestar.ValidatePathPattern(path) {
			return fmt.Errorf("configuration contains an invalid path pattern: %q", path)
		}
	}

	return nil
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
