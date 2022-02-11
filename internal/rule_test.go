package internal_test

import (
	"reflect"
	"testing"

	. "github.com/riposo/schema/internal"
	"gopkg.in/yaml.v3"
)

func TestStringSlice(t *testing.T) {
	var typ struct {
		Vals StringSlice
	}
	if err := yaml.Unmarshal([]byte(`{"vals":["a", "b"]}`), &typ); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exp, got := StringSlice([]string{"a", "b"}), typ.Vals; !reflect.DeepEqual(exp, got) {
		t.Fatalf("expected %v, got %v", exp, got)
	}

	if err := yaml.Unmarshal([]byte(`{"vals":"x"}`), &typ); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exp, got := StringSlice([]string{"x"}), typ.Vals; !reflect.DeepEqual(exp, got) {
		t.Fatalf("expected %v, got %v", exp, got)
	}
}
