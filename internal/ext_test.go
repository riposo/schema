package internal

import (
	"github.com/riposo/riposo/pkg/api"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func SeedHook(s *jsonschema.Schema) api.Hook {
	return &hook{s: s}
}
