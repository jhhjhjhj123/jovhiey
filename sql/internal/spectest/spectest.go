// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package spectest

import (
	"reflect"
	"testing"

	"ariga.io/atlas/internal/types"
	"ariga.io/atlas/schema/schemaspec"
	"ariga.io/atlas/sql/internal/specutil"

	"github.com/stretchr/testify/require"
)

// RegistrySanityTest runs a sanity for a Registry, generated a dummy *schemaspec.Type
// then converting it to a schema.Type and back to a *schemaspec.Type.
func RegistrySanityTest(t *testing.T, registry *types.Registry, skip []string) {
	for _, ts := range registry.Specs() {
		if contains(ts.Name, skip) {
			continue
		}
		t.Run(ts.Name, func(t *testing.T) {
			spec := dummyType(t, ts)
			styp, err := registry.Type(spec, nil)
			require.NoError(t, err)
			require.NoErrorf(t, err, "failed formatting: %styp", err)
			convert, err := registry.Convert(styp)
			require.NoError(t, err)
			after, err := registry.Type(convert, nil)
			require.NoError(t, err)
			require.EqualValues(t, styp, after)
		})
	}
}

func contains(s string, l []string) bool {
	for i := range l {
		if s == l[i] {
			return true
		}
	}
	return false
}

func dummyType(t *testing.T, ts *schemaspec.TypeSpec) *schemaspec.Type {
	spec := &schemaspec.Type{T: ts.T}
	for _, attr := range ts.Attributes {
		var a *schemaspec.Attr
		switch attr.Kind {
		case reflect.Int, reflect.Int64:
			a = specutil.LitAttr(attr.Name, "2")
		case reflect.String:
			a = specutil.LitAttr(attr.Name, `"a"`)
		case reflect.Slice:
			a = specutil.ListAttr(attr.Name, `"a"`, `"b"`)
		case reflect.Bool:
			a = specutil.LitAttr(attr.Name, "false")
		default:
			t.Fatalf("unsupported kind: %s", attr.Kind)
		}
		spec.Attrs = append(spec.Attrs, a)
	}
	return spec
}
