package internal

import (
	"strings"

	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/schema"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"github.com/tidwall/gjson"
)

type callbacks struct {
	api.NoopCallbacks

	globs  []string
	schema *jsonschema.Schema
}

func (cc *callbacks) Match(p riposo.Path) bool {
	return p.Match(cc.globs...)
}

func (cc *callbacks) OnCreate(_ *api.Txn, _ riposo.Path) api.CreateCallback { return cc }
func (cc *callbacks) OnUpdate(_ *api.Txn, _ riposo.Path) api.UpdateCallback { return cc }
func (cc *callbacks) OnPatch(_ *api.Txn, _ riposo.Path) api.PatchCallback   { return cc }

func (cc *callbacks) BeforeCreate(payload *schema.Resource) error {
	return cc.validate(payload.Data)
}
func (*callbacks) AfterCreate(_ *schema.Resource) error {
	return nil
}

func (cc *callbacks) BeforeUpdate(_ *schema.Object, payload *schema.Resource) error {
	return cc.validate(payload.Data)
}
func (*callbacks) AfterUpdate(_ *schema.Resource) error {
	return nil
}

func (cc *callbacks) BeforePatch(exst *schema.Object, payload *schema.Resource) error {
	sim := exst.Copy()
	if err := sim.Patch(payload.Data); err != nil {
		return err
	}
	return cc.validate(sim)
}
func (*callbacks) AfterPatch(_ *schema.Resource) error {
	return nil
}

func (cc *callbacks) validate(obj *schema.Object) error {
	extra := gjson.ParseBytes(obj.Extra)
	if !extra.IsObject() {
		return schema.InvalidBody("data", "Is not an object")
	}

	if err := cc.schema.Validate(extra.Value()); err != nil {
		msg := strings.TrimPrefix(err.Error(), "jsonschema: ")
		return schema.InvalidBody("data", msg)
	}

	return nil
}
