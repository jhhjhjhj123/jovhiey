// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package main

import (
	"bytes"
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// Job defines an integration job to run.
type Job struct {
	Version string   // version to test (passed to go test as flag which database dialect/version)
	Image   string   // name of service
	Regex   string   // run regex
	Env     []string // env of service
	Ports   []string // port mappings
	Options []string // other options
}

var (
	//go:embed ci_dialect.tmpl ci_go.tmpl
	tplFS embed.FS
	tpl   = template.Must(template.ParseFS(tplFS, "*.tmpl"))

	mysqlOptions = []string{
		`--health-cmd "mysqladmin ping -ppass"`,
		`--health-interval 10s`,
		`--health-start-period 10s`,
		`--health-timeout 5s`,
		`--health-retries 10`,
	}
	mysqlEnv = []string{
		"MYSQL_DATABASE: test",
		"MYSQL_ROOT_PASSWORD: pass",
	}
	pgOptions = []string{
		"--health-cmd pg_isready",
		"--health-interval 10s",
		"--health-timeout 5s",
		"--health-retries 5",
	}
	pgEnv = []string{
		"POSTGRES_DB: test",
		"POSTGRES_PASSWORD: pass",
	}
	jobs = []Job{
		{
			Version: "mysql56",
			Image:   "mysql:5.6.35",
			Regex:   "MySQL",
			Env:     mysqlEnv,
			Ports:   []string{"3306:3306"},
			Options: mysqlOptions,
		},
		{
			Version: "mysql57",
			Image:   "mysql:5.7.26",
			Regex:   "MySQL",
			Env:     mysqlEnv,
			Ports:   []string{"3307:3306"},
			Options: mysqlOptions,
		},
		{
			Version: "mysql8",
			Image:   "mysql:8",
			Regex:   "MySQL",
			Env:     mysqlEnv,
			Ports:   []string{"3308:3306"},
			Options: mysqlOptions,
		},
		{
			Version: "maria107",
			Image:   "mariadb:10.7",
			Regex:   "MySQL",
			Env:     mysqlEnv,
			Ports:   []string{"4306:3306"},
			Options: mysqlOptions,
		},
		{
			Version: "maria102",
			Image:   "mariadb:10.2.32",
			Regex:   "MySQL",
			Env:     mysqlEnv,
			Ports:   []string{"4307:3306"},
			Options: mysqlOptions,
		},
		{
			Version: "maria103",
			Image:   "mariadb:10.3.13",
			Regex:   "MySQL",
			Env:     mysqlEnv,
			Ports:   []string{"4308:3306"},
			Options: mysqlOptions,
		},
		{
			Version: "postgres-ext-postgis",
			Image:   "postgis/postgis:latest",
			Regex:   "Postgres",
			Env:     pgEnv,
			Ports:   []string{"5429:5432"},
			Options: pgOptions,
		},
		{
			Version: "postgres10",
			Image:   "postgres:10",
			Regex:   "Postgres",
			Env:     pgEnv,
			Ports:   []string{"5430:5432"},
			Options: pgOptions,
		},
		{
			Version: "postgres11",
			Image:   "postgres:11",
			Regex:   "Postgres",
			Env:     pgEnv,
			Ports:   []string{"5431:5432"},
			Options: pgOptions,
		},
		{
			Version: "postgres12",
			Image:   "postgres:12.3",
			Regex:   "Postgres",
			Env:     pgEnv,
			Ports:   []string{"5432:5432"},
			Options: pgOptions,
		},
		{
			Version: "postgres13",
			Image:   "postgres:13.1",
			Regex:   "Postgres",
			Env:     pgEnv,
			Ports:   []string{"5433:5432"},
			Options: pgOptions,
		},
		{
			Version: "postgres14",
			Image:   "postgres:14",
			Regex:   "Postgres",
			Env:     pgEnv,
			Ports:   []string{"5434:5432"},
			Options: pgOptions,
		},
		{
			Version: "postgres15",
			Image:   "postgres:15",
			Regex:   "Postgres",
			Env:     pgEnv,
			Ports:   []string{"5435:5432"},
			Options: pgOptions,
		},
		{
			Version: "sqlite",
			Regex:   "SQLite.*",
		},
		{
			Version: "sqlserver-2022",
			Image:   "mcr.microsoft.com/mssql/server:2022-latest",
			Regex:   "SQLServer",
			Ports:   []string{"1433:1433"},
			Env: []string{
				"ACCEPT_EULA: Y",
				"MSSQL_SA_PASSWORD: Passw0rd!995",
			},
		},
		{
			Version: "azure-sql-edge",
			Image:   "mcr.microsoft.com/azure-sql-edge:1.0.7",
			Regex:   "SQLServer",
			Ports:   []string{"1434:1433"},
			Env: []string{
				"ACCEPT_EULA: Y",
				"MSSQL_SA_PASSWORD: Passw0rd!995",
			},
		},
	}
)

func main() {
	var flavor, tags, suffix string
	flag.StringVar(&flavor, "flavor", "", "")
	flag.StringVar(&tags, "tags", "", "")
	flag.StringVar(&suffix, "suffix", "", "")
	flag.Parse()
	for _, n := range []string{"go", "dialect"} {
		var buf bytes.Buffer
		if err := tpl.ExecuteTemplate(&buf, fmt.Sprintf("ci_%s.tmpl", n), struct {
			Jobs         []Job
			Flavor, Tags string
		}{jobs, flavor, tags}); err != nil {
			log.Fatalln(err)
		}
		err := os.WriteFile(filepath.Clean(fmt.Sprintf("../../.github/workflows/ci-%s_%s.yaml", n, suffix)), buf.Bytes(), 0600)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
