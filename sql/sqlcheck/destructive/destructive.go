// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package destructive

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"ariga.io/atlas/schemahcl"
	"ariga.io/atlas/sql/internal/sqlx"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlcheck"
)

// Analyzer checks for destructive changes.
type Analyzer struct {
	sqlcheck.Options
}

// New creates a new destructive changes Analyzer with the given options.
func New(r *schemahcl.Resource) (*Analyzer, error) {
	az := &Analyzer{}
	az.Error = sqlx.P(true)
	if r, ok := r.Resource("destructive"); ok {
		if err := r.As(&az.Options); err != nil {
			return nil, fmt.Errorf("sql/sqlcheck: parsing destructive check options: %w", err)
		}
	}
	return az, nil
}

// Analyze implements sqlcheck.Analyzer.
func (a *Analyzer) Analyze(_ context.Context, p *sqlcheck.Pass) error {
	var diags []sqlcheck.Diagnostic
	for _, sc := range p.File.Changes {
		for _, c := range sc.Changes {
			switch c := c.(type) {
			case *schema.DropSchema:
				if p.File.SchemaSpan(c.S) != sqlcheck.SpanTemporary {
					var text string
					switch n := len(c.S.Tables); {
					case n == 0:
						text = fmt.Sprintf("Dropping schema %q", c.S.Name)
					case n == 1:
						text = fmt.Sprintf("Dropping non-empty schema %q with 1 table", c.S.Name)
					case n > 1:
						text = fmt.Sprintf("Dropping non-empty schema %q with %d tables", c.S.Name, n)
					}
					diags = append(diags, sqlcheck.Diagnostic{Pos: sc.Pos, Text: text})
				}
			case *schema.DropTable:
				if p.File.SchemaSpan(c.T.Schema) != sqlcheck.SpanDropped && p.File.TableSpan(c.T) != sqlcheck.SpanTemporary {
					diags = append(diags, sqlcheck.Diagnostic{
						Pos:  sc.Pos,
						Text: fmt.Sprintf("Dropping table %q", c.T.Name),
					})
				}
			case *schema.ModifyTable:
				for i := range c.Changes {
					d, ok := c.Changes[i].(*schema.DropColumn)
					if !ok || p.File.ColumnSpan(c.T, d.C) == sqlcheck.SpanTemporary {
						continue
					}
					if g := (schema.GeneratedExpr{}); !sqlx.Has(d.C.Attrs, &g) || strings.ToUpper(g.Type) != "VIRTUAL" {
						diags = append(diags, sqlcheck.Diagnostic{
							Pos:  sc.Pos,
							Text: fmt.Sprintf("Dropping non-virtual column %q", d.C.Name),
						})
					}
				}
			}
		}
	}
	if len(diags) > 0 {
		text := reportText(len(diags))
		p.Reporter.WriteReport(sqlcheck.Report{Text: text, Diagnostics: diags})
		if sqlx.V(a.Error) {
			return errors.New(text)
		}
	}
	return nil
}

func reportText(n int) string {
	if n > 1 {
		return fmt.Sprintf("%d destructive changes detected", n)
	}
	return "destructive change detected"
}
