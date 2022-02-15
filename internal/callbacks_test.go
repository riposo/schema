package internal_test

import (
	"testing"

	"github.com/riposo/riposo/pkg/api"
	"github.com/riposo/riposo/pkg/mock"
	"github.com/riposo/riposo/pkg/riposo"
	"github.com/riposo/riposo/pkg/schema"
	"github.com/riposo/schema/internal"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

func TestCallbacks(t *testing.T) {
	js, err := jsonschema.CompileString("person.json", `{
		"$id":"mock:///person.json",
		"$schema":"https://json-schema.org/draft/2020-12/schema",
		"title":"Person",
		"type":"object",
		"required":["firstName","lastName"],
		"properties":{
			"firstName":{"type":"string","minLength":3},
			"lastName":{"type":"string"}
		}
	}`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	act := api.NewActions(api.DefaultModel{}, []api.Callbacks{internal.SeedCallbacks(js)})
	txn := mock.Txn()
	defer txn.Abort()

	t.Run("Create", func(t *testing.T) {
		err = act.Create(txn, "/buckets/foo/people/*", &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{}`)},
		})
		if exp := `data in body: '' does not validate with mock:///person.json#/required: missing properties: 'firstName', 'lastName'`; exp != err.Error() {
			t.Fatalf("expected %v, got %v", exp, err)
		}

		err = act.Create(txn, "/buckets/foo/people/*", &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{"firstName":"J","lastName":"Doe"}`)},
		})
		if exp := `data in body: '/firstName' does not validate with mock:///person.json#/properties/firstName/minLength: length must be >= 3, but got 1`; exp != err.Error() {
			t.Fatalf("expected %v, got %v", exp, err)
		}

		err = act.Create(txn, "/buckets/foo/people/*", &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{"firstName":"Jane","lastName":"Doe"}`)},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Update", func(t *testing.T) {
		res := &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{"firstName":"Jane","lastName":"Doe"}`)},
		}
		if err := act.Create(txn, "/buckets/foo/people/*", res); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		path := riposo.Path("/buckets/foo/people/" + res.Data.ID)

		handle, err := txn.Store.GetForUpdate(path)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		res, err = act.Update(txn, path, handle, &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{"firstName":"Alice","lastName":"Glass"}`)},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		} else if exp, got := `{"id":"ITR.ID","last_modified":1515151515679,"firstName":"Alice","lastName":"Glass"}`, res.Data.String(); exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}

		_, err = act.Update(txn, path, handle, &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{"firstName":"J","lastName":"Doe"}`)},
		})
		if exp := `data in body: '/firstName' does not validate with mock:///person.json#/properties/firstName/minLength: length must be >= 3, but got 1`; exp != err.Error() {
			t.Fatalf("expected %v, got %v", exp, err)
		}
	})

	t.Run("Patch", func(t *testing.T) {
		res := &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{"firstName":"Jane","lastName":"Doe"}`)},
		}
		if err := act.Create(txn, "/buckets/foo/people/*", res); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		path := riposo.Path("/buckets/foo/people/" + res.Data.ID)

		handle, err := txn.Store.GetForUpdate(path)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		res, err = act.Patch(txn, path, handle, &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{"firstName":"Alice"}`)},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		} else if exp, got := `{"id":"MXR.ID","last_modified":1515151515681,"firstName":"Alice","lastName":"Doe"}`, res.Data.String(); exp != got {
			t.Fatalf("expected %v, got %v", exp, got)
		}

		_, err = act.Patch(txn, path, handle, &schema.Resource{
			Data: &schema.Object{Extra: []byte(`{"firstName":"J"}`)},
		})
		if exp := `data in body: '/firstName' does not validate with mock:///person.json#/properties/firstName/minLength: length must be >= 3, but got 1`; exp != err.Error() {
			t.Fatalf("expected %v, got %v", exp, err)
		}
	})
}
