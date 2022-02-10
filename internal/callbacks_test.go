package internal_test

import (
	"testing"

	"github.com/riposo/riposo/pkg/schema"
	"github.com/riposo/schema/internal"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func TestCallbacks(t *testing.T) {
	js, err := jsonschema.CompileString("person.json", `{
		"$id": "mock:///person.json",
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"title": "Person",
		"type": "object",
		"required": [ "firstName", "lastName" ],
		"properties": {
			"firstName": { "type": "string", "minLength": 3 },
			"lastName": { "type": "string" }
		}
	}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	callbacks := internal.SeedCallbacks(js)

	t.Run("OnCreate", func(t *testing.T) {
		err := callbacks.OnCreate(nil, "").BeforeCreate(&schema.Resource{
			Data: &schema.Object{Extra: []byte(`{}`)},
		})
		if exp := `data in body: '' does not validate with mock:///person.json#/required: missing properties: 'firstName', 'lastName'`; exp != err.Error() {
			t.Fatalf("expected %v, got %v", exp, err)
		}

		err = callbacks.OnCreate(nil, "").BeforeCreate(&schema.Resource{
			Data: &schema.Object{Extra: []byte(`{ "firstName": "J", "lastName": "Doe" }`)},
		})
		if exp := `data in body: '/firstName' does not validate with mock:///person.json#/properties/firstName/minLength: length must be >= 3, but got 1`; exp != err.Error() {
			t.Fatalf("expected %v, got %v", exp, err)
		}

		err = callbacks.OnCreate(nil, "").BeforeCreate(&schema.Resource{
			Data: &schema.Object{Extra: []byte(`{ "firstName": "Jane", "lastName": "Doe" }`)},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("OnUpdate", func(t *testing.T) {
		err := callbacks.OnUpdate(nil, "").BeforeUpdate(nil, &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{ "firstName": "J", "lastName": "Doe" }`)},
		})
		if exp := `data in body: '/firstName' does not validate with mock:///person.json#/properties/firstName/minLength: length must be >= 3, but got 1`; exp != err.Error() {
			t.Fatalf("expected %v, got %v", exp, err)
		}
	})

	t.Run("OnPatch", func(t *testing.T) {
		exst := &schema.Object{Extra: []byte(`{ "firstName": "Jane", "lastName": "Doe" }`)}
		err := callbacks.OnPatch(nil, "").BeforePatch(exst, &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{ "firstName": "Alice" }`)},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		err = callbacks.OnPatch(nil, "").BeforePatch(exst, &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{ "firstName": "J" }`)},
		})
		if exp := `data in body: '/firstName' does not validate with mock:///person.json#/properties/firstName/minLength: length must be >= 3, but got 1`; exp != err.Error() {
			t.Fatalf("expected %v, got %v", exp, err)
		}
	})
}
