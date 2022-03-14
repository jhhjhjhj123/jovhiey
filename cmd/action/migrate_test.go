// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package action

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMigrate_Diff(t *testing.T) {
	_, err := runCmd(RootCmd, "migrate")
	require.NoError(t, err)
	p := t.TempDir()

	// Expect no clean dev error.
	s, err := runCmd(
		RootCmd, "migrate", "diff",
		"name",
		"--dir", "file://"+p,
		"--dev-url", openSQLite(t, "create table t (c int);"),
		"--to", hclURL(t),
	)
	require.True(t, strings.HasPrefix(s, "Error: dev database must be clean"))
	require.EqualError(t, err, "dev database must be clean")

	// Works.
	s, err = runCmd(
		RootCmd, "migrate", "diff",
		"name",
		"--dir", "file://"+p,
		"--dev-url", openSQLite(t, ""),
		"--to", hclURL(t),
	)
	require.Zero(t, s)
	require.NoError(t, err)
	require.FileExists(t, filepath.Join(p, fmt.Sprintf("%s_name.sql", time.Now().Format("20060102150405"))))
	require.FileExists(t, filepath.Join(p, "atlas.sum"))
}

const hcl = `
schema "main" {
}

table "table" {
  schema = schema.main
  column "col" {
    type    = int
    comment = "column comment"
  }
  column "age" {
    type = int
  }
  column "price1" {
    type = int
  }
  column "price2" {
    type           = int
  }
  column "account_name" {
    type = varchar(32)
    null = true
  }
  column "created_at" {
    type    = datetime
    default = sql("current_timestamp")
  }
  primary_key {
    columns = [table.table.column.col]
  }
  index "index" {
    unique  = true
    columns = [
      table.table.column.col,
      table.table.column.age,
    ]
    comment = "index comment"
  }
  foreign_key "accounts" {
    columns = [
      table.table.column.account_name,
    ]
    ref_columns = [
      table.accounts.column.name,
    ]
    on_delete = SET_NULL
    on_update = "NO_ACTION"
  }
  check "positive price" {
    expr = "price1 > 0"
  }
  check {
    expr     = "price1 <> price2"
    enforced = true
  }
  check {
    expr     = "price2 <> price1"
    enforced = false
  }
  comment        = "table comment"
}

table "accounts" {
  schema = schema.main
  column "name" {
    type = varchar(32)
  }
  column "unsigned_float" {
    type     = float(10)
    unsigned = true
  }
  column "unsigned_decimal" {
    type     = decimal(10, 2)
    unsigned = true
  }
  primary_key {
    columns = [table.accounts.column.name]
  }
}`

func hclURL(t *testing.T) string {
	p := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(p, "atlas.hcl"), []byte(hcl), 0640))
	return "file://" + filepath.Join(p, "atlas.hcl")
}
