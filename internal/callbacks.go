package internal

import (
	"strings"

	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/schema"
	"github.com/tidwall/gjson"
)

type callbacks struct {
	api.NoopCallbacks
	rs ruleSet
}

func (cc *callbacks) OnCreate(_ *api.Txn, path riposo.Path) api.CreateCallback {
	if cc.match(path) {
		return &handler{path: path, rs: cc.rs}
	}
	return nil
}

func (cc *callbacks) OnUpdate(_ *api.Txn, path riposo.Path) api.UpdateCallback {
	if cc.match(path) {
		return &handler{path: path, rs: cc.rs}
	}
	return nil
}

func (cc *callbacks) OnPatch(_ *api.Txn, path riposo.Path) api.PatchCallback {
	if cc.match(path) {
		return &handler{path: path, rs: cc.rs}
	}
	return nil
}
func (cc *callbacks) match(path riposo.Path) bool {
	for _, rc := range cc.rs {
		if path.Match(rc.globs...) {
			return true
		}
	}
	return false
}

// --------------------------------------------------------------------

type handler struct {
	path riposo.Path
	rs   ruleSet
}

func (h *handler) BeforeCreate(payload *schema.Resource) error {
	return h.validate(payload.Data)
}
func (*handler) AfterCreate(_ *schema.Resource) error {
	return nil
}

func (h *handler) BeforeUpdate(_ *schema.Object, payload *schema.Resource) error {
	return h.validate(payload.Data)
}
func (*handler) AfterUpdate(_ *schema.Resource) error {
	return nil
}

func (h *handler) BeforePatch(exst *schema.Object, payload *schema.Resource) error {
	sim := exst.Copy()
	if err := sim.Patch(payload.Data); err != nil {
		return err
	}
	return h.validate(sim)
}
func (*handler) AfterPatch(_ *schema.Resource) error {
	return nil
}

func (h *handler) validate(obj *schema.Object) error {
	extra := gjson.ParseBytes(obj.Extra)
	if !extra.IsObject() {
		return schema.InvalidBody("data", "Is not an object")
	}

	for _, rc := range h.rs {
		if !h.path.Match(rc.globs...) {
			continue
		}

		if err := rc.schema.Validate(extra.Value()); err != nil {
			msg := strings.TrimPrefix(err.Error(), "jsonschema: ")
			return schema.InvalidBody("data", msg)
		}
	}
	return nil
}
