package internal

import (
	"strings"

	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/schema"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/tidwall/gjson"
)

type hook struct {
	api.NoopHook

	s *jsonschema.Schema
}

func (h *hook) BeforeCreate(_ *api.Txn, _ riposo.Path, payload *schema.Resource) error {
	return h.validate(payload.Data)
}

func (h *hook) BeforeUpdate(_ *api.Txn, _ riposo.Path, _ *schema.Object, payload *schema.Resource) error {
	return h.validate(payload.Data)
}

func (h *hook) BeforePatch(_ *api.Txn, _ riposo.Path, exst *schema.Object, payload *schema.Resource) error {
	sim := exst.Copy()
	if err := sim.Patch(payload.Data); err != nil {
		return err
	}
	return h.validate(sim)
}

func (h *hook) validate(obj *schema.Object) error {
	extra := gjson.ParseBytes(obj.Extra)
	if !extra.IsObject() {
		return schema.InvalidBody("data", "Is not an object")
	}

	if err := h.s.Validate(extra.Value()); err != nil {
		msg := strings.TrimPrefix(err.Error(), "jsonschema: ")
		return schema.InvalidBody("data", msg)
	}

	return nil
}
