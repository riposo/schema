package internal

import (
	"github.com/riposo/riposo/pkg/api"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func SeedCallbacks(schema *jsonschema.Schema) api.Callbacks {
	return &callbacks{
		rs: ruleSet{
			{globs: []string{"/buckets/*/people/*"}, schema: schema},
		},
	}
}
