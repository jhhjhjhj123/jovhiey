// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package sqlx

import (
	"strconv"
	"testing"

	"ariga.io/atlas/sql/schema"

	"github.com/stretchr/testify/require"
)

func TestModeInspectRealm(t *testing.T) {
	m := ModeInspectRealm(nil)
	require.True(t, m.Is(schema.InspectSchemas))
	require.True(t, m.Is(schema.InspectTables))

	m = ModeInspectRealm(&schema.InspectRealmOption{})
	require.True(t, m.Is(schema.InspectSchemas))
	require.True(t, m.Is(schema.InspectTables))

	m = ModeInspectRealm(&schema.InspectRealmOption{Mode: schema.InspectSchemas})
	require.True(t, m.Is(schema.InspectSchemas))
	require.False(t, m.Is(schema.InspectTables))
}

func TestModeInspectSchema(t *testing.T) {
	m := ModeInspectSchema(nil)
	require.True(t, m.Is(schema.InspectSchemas))
	require.True(t, m.Is(schema.InspectTables))

	m = ModeInspectSchema(&schema.InspectOptions{})
	require.True(t, m.Is(schema.InspectSchemas))
	require.True(t, m.Is(schema.InspectTables))

	m = ModeInspectSchema(&schema.InspectOptions{Mode: schema.InspectSchemas})
	require.True(t, m.Is(schema.InspectSchemas))
	require.False(t, m.Is(schema.InspectTables))
}

func TestBuilder(t *testing.T) {
	var (
		b       = &Builder{QuoteChar: '"'}
		columns = []string{"a", "b", "c"}
	)
	b.P("CREATE TABLE").
		Table(&schema.Table{Name: "users"}).
		Wrap(func(b *Builder) {
			b.MapComma(columns, func(i int, b *Builder) {
				b.Ident(columns[i]).P("int").P("NOT NULL")
			})
			b.Comma().P("PRIMARY KEY").Wrap(func(b *Builder) {
				b.MapComma(columns, func(i int, b *Builder) {
					b.Ident(columns[i])
				})
			})
		})
	require.Equal(t, `CREATE TABLE "users" ("a" int NOT NULL, "b" int NOT NULL, "c" int NOT NULL, PRIMARY KEY ("a", "b", "c"))`, b.String())
}

func TestMayWrap(t *testing.T) {
	tests := []struct {
		input   string
		wrapped bool
	}{
		{"", true},
		{"()", false},
		{"('text')", false},
		{"('(')", false},
		{`('(\\')`, false},
		{`('\')(')`, false},
		{`(a) in (b)`, true},
		{`a in (b)`, true},
		{`("\\\\(((('")`, false},
		{`('(')||(')')`, true},
		// Test examples from SQLite.
		{"b || 'x'", true},
		{"a+1", true},
		{"substr(x, 2)", true},
		{"(json_extract(x, '$.a'), json_extract(x, '$.b'))", false},
		{"(substr(a, 2) COLLATE NOCASE, b)", false},
		{"(a,b+random())", false},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			expect := tt.input
			if tt.wrapped {
				expect = "(" + expect + ")"
			}
			require.Equal(t, expect, MayWrap(tt.input))

		})
	}
}
