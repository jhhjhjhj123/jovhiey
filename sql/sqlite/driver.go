// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package sqlite

import (
	"fmt"

	"ariga.io/atlas/sql/internal/sqlx"
	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/schema"
)

type (
	// Driver represents a SQLite driver for introspecting database schemas,
	// generating diff between schema elements and apply migrations changes.
	Driver struct {
		conn
		schema.Differ
		schema.Inspector
		migrate.PlanApplier
	}

	// database connection and its information.
	conn struct {
		schema.ExecQuerier
		// System variables that are set on `Open`.
		fkEnabled  bool
		version    string
		collations []string
	}
)

// Open opens a new SQLite driver.
func Open(db schema.ExecQuerier) (*Driver, error) {
	c := conn{ExecQuerier: db}
	if err := db.QueryRow("SELECT sqlite_version()").Scan(&c.version); err != nil {
		return nil, fmt.Errorf("sqlite: scanning database version: %w", err)
	}
	if err := db.QueryRow("PRAGMA foreign_keys").Scan(&c.fkEnabled); err != nil {
		return nil, fmt.Errorf("sqlite: check foreign_keys pragma: %w", err)
	}
	rows, err := db.Query("SELECT name FROM pragma_collation_list()")
	if err != nil {
		return nil, fmt.Errorf("sqlite: check collation_list pragma: %w", err)
	}
	if c.collations, err = sqlx.ScanStrings(rows); err != nil {
		return nil, fmt.Errorf("sqlite: scanning database collations: %w", err)
	}
	return &Driver{
		conn:        c,
		Differ:      &sqlx.Diff{DiffDriver: &diff{c}},
		Inspector:   &inspect{c},
		PlanApplier: &planApply{c},
	}, nil
}

// SQLite standard data types as defined in its codebase and documentation.
// https://www.sqlite.org/datatype3.html
// https://github.com/sqlite/sqlite/blob/master/src/global.c
const (
	TypeInteger = "integer" // SQLITE_TYPE_INTEGER
	TypeReal    = "real"    // SQLITE_TYPE_REAL
	TypeText    = "text"    // SQLITE_TYPE_TEXT
	TypeBlob    = "blob"    // SQLITE_TYPE_BLOB
)
