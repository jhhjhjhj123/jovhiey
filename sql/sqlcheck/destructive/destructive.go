// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package destructive

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"ariga.io/atlas/schemahcl"
	"ariga.io/atlas/sql/internal/sqlx"
	"ariga.io/atlas/sql/migrate"
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
	if r, ok := r.Resource(az.Name()); ok {
		if err := r.As(&az.Options); err != nil {
			return nil, fmt.Errorf("sql/sqlcheck: parsing destructive check options: %w", err)
		}
	}
	return az, nil
}

// List of codes.
var (
	codeDropS = sqlcheck.Code("DS101")
	codeDropT = sqlcheck.Code("DS102")
	codeDropC = sqlcheck.Code("DS103")
)

// Name of the analyzer. Implements the sqlcheck.NamedAnalyzer interface.
func (*Analyzer) Name() string {
	return "destructive"
}

// Analyze implements sqlcheck.Analyzer.
func (a *Analyzer) Analyze(_ context.Context, p *sqlcheck.Pass) error {
	var (
		edits []*migrate.Stmt
		diags []sqlcheck.Diagnostic
	)
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
					diags = append(diags, sqlcheck.Diagnostic{
						Code: codeDropS,
						Pos:  sc.Stmt.Pos,
						Text: text,
					})
				}
			case *schema.DropTable:
				if p.File.SchemaSpan(c.T.Schema) != sqlcheck.SpanDropped && p.File.TableSpan(c.T) != sqlcheck.SpanTemporary && !a.hasEmptyTableCheck(p, c.T) {
					diags = append(diags, sqlcheck.Diagnostic{
						Code: codeDropT,
						Pos:  sc.Stmt.Pos,
						Text: fmt.Sprintf("Dropping table %q", c.T.Name),
						SuggestedFixes: []sqlcheck.SuggestedFix{
							{
								Message: fmt.Sprintf("Add a pre-migration check to ensure table %q is empty before dropping it", c.T.Name),
							},
						},
					})
					if stmt, err := a.emptyTableCheckStmt(p, c.T); err == nil {
						edits = append(edits, stmt)
					}
				}
			case *schema.ModifyTable:
				var names []string
				for i := range c.Changes {
					d, ok := c.Changes[i].(*schema.DropColumn)
					if !ok || p.File.ColumnSpan(c.T, d.C) == sqlcheck.SpanTemporary {
						continue
					}
					if g := (schema.GeneratedExpr{}); !sqlx.Has(d.C.Attrs, &g) || strings.ToUpper(g.Type) != "VIRTUAL" {
						names = append(names, strconv.Quote(d.C.Name))
						if stmt, err := a.emptyColumnCheckStmt(p, c.T, d.C.Name); err == nil {
							edits = append(edits, stmt)
						}
					}
				}
				switch n := len(names); {
				case n == 1:
					diags = append(diags, sqlcheck.Diagnostic{
						Code: codeDropC,
						Pos:  sc.Stmt.Pos,
						Text: fmt.Sprintf("Dropping non-virtual column %s", names[0]),
						SuggestedFixes: []sqlcheck.SuggestedFix{
							{
								Message: fmt.Sprintf("Add a pre-migration check to ensure column %q is NULL before dropping it", c.T.Name),
							},
						},
					})
				case n > 1:
					// All changes generated by the same statement (same position).
					diags = append(diags, sqlcheck.Diagnostic{
						Code: codeDropC,
						Pos:  sc.Stmt.Pos,
						Text: fmt.Sprintf("Dropping non-virtual columns %s and %s", strings.Join(names[:n-1], ", "), names[n-1]),
						SuggestedFixes: []sqlcheck.SuggestedFix{
							{
								Message: fmt.Sprintf("Add pre-migration checks to ensure columns %s and %s are NULL before dropping them", strings.Join(names[:n-1], ", "), names[n-1]),
							},
						},
					})
				}
			}
		}
	}
	if len(diags) > 0 {
		const reportText = "destructive changes detected"
		p.Reporter.WriteReport(sqlcheck.Report{Text: reportText, Diagnostics: diags, SuggestedFixes: suggestFix(p, edits)})
		if sqlx.V(a.Error) {
			return errors.New(reportText)
		}
	}
	return nil
}
