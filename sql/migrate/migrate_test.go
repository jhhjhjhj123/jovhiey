// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package migrate_test

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"io/fs"
	"path/filepath"
	"testing"
	"text/template"
	"time"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/schema"

	"github.com/stretchr/testify/require"
)

func TestPlanner_WritePlan(t *testing.T) {
	p := t.TempDir()
	d, err := migrate.NewLocalDir(p)
	require.NoError(t, err)
	plan := &migrate.Plan{
		Name: "add_t1_and_t2",
		Changes: []*migrate.Change{
			{Cmd: "CREATE TABLE t1(c int)", Reverse: "DROP TABLE t1 IF EXISTS"},
			{Cmd: "CREATE TABLE t2(c int)", Reverse: "DROP TABLE t2"},
		},
	}

	// DefaultFormatter
	pl := migrate.NewPlanner(nil, d, migrate.DisableChecksum())
	require.NotNil(t, pl)
	require.NoError(t, pl.WritePlan(plan))
	v := time.Now().UTC().Format("20060102150405")
	require.Equal(t, countFiles(t, d), 1)
	requireFileEqual(t, d, v+"_add_t1_and_t2.sql", "CREATE TABLE t1(c int);\nCREATE TABLE t2(c int);\n")

	// Custom formatter (creates "up" and "down" migration files).
	fmt, err := migrate.NewTemplateFormatter(
		template.Must(template.New("").Parse("{{ .Name }}.up.sql")),
		template.Must(template.New("").Parse("{{ range .Changes }}{{ println .Cmd }}{{ end }}")),
		template.Must(template.New("").Parse("{{ .Name }}.down.sql")),
		template.Must(template.New("").Parse("{{ range .Changes }}{{ println .Reverse }}{{ end }}")),
	)
	require.NoError(t, err)
	pl = migrate.NewPlanner(nil, d, migrate.WithFormatter(fmt), migrate.DisableChecksum())
	require.NotNil(t, pl)
	require.NoError(t, pl.WritePlan(plan))
	require.Equal(t, countFiles(t, d), 3)
	requireFileEqual(t, d, "add_t1_and_t2.up.sql", "CREATE TABLE t1(c int)\nCREATE TABLE t2(c int)\n")
	requireFileEqual(t, d, "add_t1_and_t2.down.sql", "DROP TABLE t1 IF EXISTS\nDROP TABLE t2\n")
}

func TestPlanner_Plan(t *testing.T) {
	var (
		drv = &lockMockDriver{&mockDriver{}}
		ctx = context.Background()
	)
	d, err := migrate.NewLocalDir(t.TempDir())
	require.NoError(t, err)

	// nothing to do
	pl := migrate.NewPlanner(drv, d)
	plan, err := pl.Plan(ctx, "empty", migrate.Realm(nil))
	require.ErrorIs(t, err, migrate.ErrNoPlan)
	require.Nil(t, plan)

	// there are changes
	drv.changes = []schema.Change{
		&schema.AddTable{T: schema.NewTable("t1").AddColumns(schema.NewIntColumn("c", "int"))},
		&schema.AddTable{T: schema.NewTable("t2").AddColumns(schema.NewIntColumn("c", "int"))},
	}
	drv.plan = &migrate.Plan{
		Changes: []*migrate.Change{
			{Cmd: "CREATE TABLE t1(c int);"},
			{Cmd: "CREATE TABLE t2(c int);"},
		},
	}
	plan, err = pl.Plan(ctx, "", migrate.Realm(nil))
	require.NoError(t, err)
	require.Equal(t, drv.plan, plan)
}

func TestPlanner_PlanSchemaMode(t *testing.T) {
	var (
		drv = &lockMockDriver{&mockDriver{}}
		ctx = context.Background()
	)
	d, err := migrate.NewLocalDir(t.TempDir())
	require.NoError(t, err)

	// Schema is missing in dev connection.
	pl := migrate.NewPlanner(drv, d, migrate.PlanSchema("test"))
	plan, err := pl.Plan(ctx, "empty", migrate.Realm(nil))
	require.EqualError(t, err, `missing schema "test" in realm after replaying migration directory`)
	require.Nil(t, plan)

	drv.realm = *schema.NewRealm(schema.New("test"))
	pl = migrate.NewPlanner(drv, d, migrate.PlanSchema("test"))
	plan, err = pl.Plan(ctx, "empty", migrate.Realm(schema.NewRealm()))
	require.EqualError(t, err, `missing schema definition in desired state`)
	require.Nil(t, plan)

	drv.realm = *schema.NewRealm(schema.New("test"))
	pl = migrate.NewPlanner(drv, d, migrate.PlanSchema("test"))
	plan, err = pl.Plan(ctx, "multi", migrate.Realm(schema.NewRealm(schema.New("test"), schema.New("dev"))))
	require.EqualError(t, err, `found 2 schema definitions in desired state on schema mode`)
	require.Nil(t, plan)

	drv.realm = *schema.NewRealm(schema.New("test"))
	pl = migrate.NewPlanner(drv, d, migrate.PlanSchema("test"))
	plan, err = pl.Plan(ctx, "multi", migrate.Realm(schema.NewRealm(schema.New("test"))))
	require.ErrorIs(t, err, migrate.ErrNoPlan)
	require.Nil(t, plan)
}

func TestExecutor_ReadState(t *testing.T) {
	ctx := context.Background()
	d, err := migrate.NewLocalDir("testdata/migrate")
	require.NoError(t, err)

	// Locking not supported.
	_, err = migrate.NewExecutor(&mockDriver{}, d, migrate.NopRevisionReadWriter{})
	require.ErrorIs(t, err, migrate.ErrLockUnsupported)

	drv := &lockMockDriver{&mockDriver{}}
	ex, err := migrate.NewExecutor(drv, d, migrate.NopRevisionReadWriter{})
	require.NoError(t, err)

	_, err = ex.ReadState(ctx)
	require.NoError(t, err)
	require.Equal(t, []string{"DROP TABLE IF EXISTS t;", "CREATE TABLE t(c int);"}, drv.executed)
	require.Equal(t, 2, drv.lockCounter)
	require.Equal(t, 2, drv.unlockCounter)
	require.True(t, drv.released())

	// Does not work if locked.
	drv.locks = map[string]struct{}{"atlas_migration_directory_state": {}}
	_, err = ex.ReadState(ctx)
	require.EqualError(t, err, "sql/migrate: acquiring database lock: lockErr")
	require.Equal(t, 2, drv.lockCounter)
	require.Equal(t, 2, drv.unlockCounter)
	require.False(t, drv.released())
	drv.locks = make(map[string]struct{})

	// Does not work if database is not clean.
	drv.dirty = true
	drv.realm = schema.Realm{Schemas: []*schema.Schema{{Name: "schema"}}}
	_, err = ex.ReadState(ctx)
	require.ErrorAs(t, err, &migrate.NotCleanError{})
	require.Equal(t, 3, drv.lockCounter)
	require.Equal(t, 3, drv.unlockCounter)
	require.True(t, drv.released())
}

func TestExecutor_Pending(t *testing.T) {
	var (
		drv  = &lockMockDriver{&mockDriver{}}
		rrw  = &mockRevisionReadWriter{}
		log  = &mockLogger{}
		rev1 = &migrate.Revision{
			Version:     "1.a",
			Description: "sub.up",
			Applied:     2,
			Total:       2,
			Hash:        "nXyZR020M/mH7LxkoTkJr7BcQkipVg90imQ9I4595dw=",
		}
		rev2 = &migrate.Revision{
			Version:     "2.10.x-20",
			Description: "description",
			Applied:     1,
			Total:       1,
			Hash:        "wQB3Vh3PHVXQg9OD3Gn7TBxbZN3r1Qb7TtAE1g3q9mQ=",
		}
		rev3 = &migrate.Revision{
			Version:     "3",
			Description: "partly",
			Applied:     1,
			Total:       2,
			Error:       "this is an migration error",
			Hash:        "+O40cAXHgvMClnynHd5wggPAeZAk7zSEaNgzXCZOfmY=",
		}
	)
	dir, err := migrate.NewLocalDir(filepath.Join("testdata/migrate", "sub"))
	require.NoError(t, err)
	ex, err := migrate.NewExecutor(drv, dir, rrw, migrate.WithLogger(log))
	require.NoError(t, err)

	// All are pending
	p, err := ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, p, 3)

	// 2 are pending.
	*rrw = mockRevisionReadWriter(migrate.Revisions{rev1})
	p, err = ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, p, 2)

	// Only the last one is pending (in full).
	*rrw = mockRevisionReadWriter(migrate.Revisions{rev1, rev2})
	p, err = ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, p, 1)

	// First statement of last one is marked as applied, second isn't. Third file is still pending.
	*rrw = mockRevisionReadWriter(migrate.Revisions{rev1, rev2, rev3})
	p, err = ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, p, 1)
}

func TestExecutor(t *testing.T) {
	// Passing nil raises error.
	ex, err := migrate.NewExecutor(nil, nil, nil)
	require.EqualError(t, err, "sql/migrate: execute: no driver given")
	require.Nil(t, ex)

	ex, err = migrate.NewExecutor(&mockDriver{}, nil, nil)
	require.EqualError(t, err, "sql/migrate: execute: no dir given")
	require.Nil(t, ex)

	dir, err := migrate.NewLocalDir(t.TempDir())
	require.NoError(t, err)
	ex, err = migrate.NewExecutor(&mockDriver{}, dir, nil)
	require.EqualError(t, err, "sql/migrate: execute: no revision storage given")
	require.Nil(t, ex)

	// Does not work if no locking mechanism is provided.
	ex, err = migrate.NewExecutor(&mockDriver{}, dir, &mockRevisionReadWriter{})
	require.ErrorIs(t, err, migrate.ErrLockUnsupported)
	require.Nil(t, ex)

	// Does not operate on invalid migration dir.
	dir, err = migrate.NewLocalDir(t.TempDir())
	require.NoError(t, err)
	require.NoError(t, dir.WriteFile("atlas.sum", hash))
	ex, err = migrate.NewExecutor(&lockMockDriver{&mockDriver{}}, dir, &mockRevisionReadWriter{})
	require.NoError(t, err)
	require.NotNil(t, ex)
	require.ErrorIs(t, ex.ExecuteN(context.Background(), 0), migrate.ErrChecksumMismatch)

	// Prerequisites.
	var (
		drv  = &lockMockDriver{&mockDriver{}}
		rrw  = &mockRevisionReadWriter{}
		log  = &mockLogger{}
		rev1 = &migrate.Revision{
			Version:     "1.a",
			Description: "sub.up",
			Type:        migrate.RevisionTypeExecute,
			Applied:     2,
			Total:       2,
			Hash:        "nXyZR020M/mH7LxkoTkJr7BcQkipVg90imQ9I4595dw=",
		}
		rev2 = &migrate.Revision{
			Version:     "2.10.x-20",
			Description: "description",
			Type:        migrate.RevisionTypeExecute,
			Applied:     1,
			Total:       1,
			Hash:        "wQB3Vh3PHVXQg9OD3Gn7TBxbZN3r1Qb7TtAE1g3q9mQ=",
		}
	)
	dir, err = migrate.NewLocalDir(filepath.Join("testdata/migrate", "sub"))
	require.NoError(t, err)
	ex, err = migrate.NewExecutor(drv, dir, rrw, migrate.WithLogger(log))
	require.NoError(t, err)

	// Applies two of them.
	require.NoError(t, ex.ExecuteN(context.Background(), 2))
	require.Equal(t, drv.executed, []string{
		"CREATE TABLE t_sub(c int);", "ALTER TABLE t_sub ADD c1 int;", "ALTER TABLE t_sub ADD c2 int;",
	})
	requireEqualRevisions(t, migrate.Revisions{rev1, rev2}, migrate.Revisions(*rrw))
	require.Equal(t, []migrate.LogEntry{
		migrate.LogExecution{To: "2.10.x-20", Files: []string{"1.a_sub.up.sql", "2.10.x-20_description.sql"}},
		migrate.LogFile{Version: "1.a", Desc: "sub.up"},
		migrate.LogStmt{SQL: "CREATE TABLE t_sub(c int);"},
		migrate.LogStmt{SQL: "ALTER TABLE t_sub ADD c1 int;"},
		migrate.LogFile{Version: "2.10.x-20", Desc: "description"},
		migrate.LogStmt{SQL: "ALTER TABLE t_sub ADD c2 int;"},
		migrate.LogDone{},
	}, []migrate.LogEntry(*log))
	require.Equal(t, drv.lockCounter, 1)
	require.Equal(t, drv.unlockCounter, 1)
	require.True(t, drv.released())

	// Partly is pending.
	p, err := ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, p, 1)
	require.Equal(t, "3_partly.sql", p[0].Name())

	// Apply one by one.
	*rrw = mockRevisionReadWriter{}
	*drv = lockMockDriver{&mockDriver{}}

	require.NoError(t, ex.ExecuteN(context.Background(), 1))
	require.Equal(t, drv.executed, []string{"CREATE TABLE t_sub(c int);", "ALTER TABLE t_sub ADD c1 int;"})
	requireEqualRevisions(t, migrate.Revisions{rev1}, migrate.Revisions(*rrw))

	require.NoError(t, ex.ExecuteN(context.Background(), 1))
	require.Equal(t, drv.executed, []string{
		"CREATE TABLE t_sub(c int);", "ALTER TABLE t_sub ADD c1 int;", "ALTER TABLE t_sub ADD c2 int;",
	})
	requireEqualRevisions(t, migrate.Revisions{rev1, rev2}, migrate.Revisions(*rrw))
	require.Equal(t, 2, drv.lockCounter)
	require.Equal(t, 2, drv.unlockCounter)
	require.True(t, drv.released())

	// Partly is pending.
	p, err = ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, p, 1)
	require.Equal(t, "3_partly.sql", p[0].Name())

	// Suppose first revision is already executed, only execute second migration file.
	*rrw = mockRevisionReadWriter(migrate.Revisions{rev1})
	*drv = lockMockDriver{&mockDriver{}}

	require.NoError(t, ex.ExecuteN(context.Background(), 1))
	require.Equal(t, []string{"ALTER TABLE t_sub ADD c2 int;"}, drv.executed)
	requireEqualRevisions(t, migrate.Revisions{rev1, rev2}, migrate.Revisions(*rrw))

	// Partly is pending.
	p, err = ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, p, 1)
	require.Equal(t, "3_partly.sql", p[0].Name())

	require.Equal(t, 1, drv.lockCounter)
	require.Equal(t, 1, drv.unlockCounter)
	require.True(t, drv.released())

	// Failing, counter will be correct.
	*rrw = mockRevisionReadWriter(migrate.Revisions{rev1, rev2})
	*drv = lockMockDriver{&mockDriver{}}
	drv.failOn(2, errors.New("this is an error"))
	require.ErrorContains(t, ex.ExecuteN(context.Background(), 1), "this is an error")
	revs, err := rrw.ReadRevisions(context.Background())
	require.NoError(t, err)
	requireEqualRevision(t, &migrate.Revision{
		Version:     "3",
		Description: "partly",
		Type:        migrate.RevisionTypeExecute,
		Applied:     1,
		Total:       2,
		Error:       "Statement:\nALTER TABLE t_sub ADD c4 int;\n\nError:\nthis is an error",
	}, revs[len(revs)-1])

	// Will fail if applied contents hash has changed (like when editing a partially applied file to fix an error).
	h := revs[len(revs)-1].PartialHashes[0]
	revs[len(revs)-1].PartialHashes[0] += h
	require.ErrorAs(t, ex.ExecuteN(context.Background(), 1), &migrate.HistoryChangedError{})

	// Re-attempting to migrate will pick up where the execution was left off.
	revs[len(revs)-1].PartialHashes[0] = h
	*drv = lockMockDriver{&mockDriver{}}
	require.NoError(t, ex.ExecuteN(context.Background(), 1))
	require.Equal(t, []string{"ALTER TABLE t_sub ADD c4 int;"}, drv.executed)

	// Everything is applied.
	require.ErrorIs(t, ex.ExecuteN(context.Background(), 0), migrate.ErrNoPendingFiles)
}

func TestExecutor_Baseline(t *testing.T) {
	var (
		rrw mockRevisionReadWriter
		drv = &lockMockDriver{&mockDriver{dirty: true}}
		log = &mockLogger{}
	)
	dir, err := migrate.NewLocalDir(filepath.Join("testdata/migrate", "sub"))
	require.NoError(t, err)
	ex, err := migrate.NewExecutor(drv, dir, &rrw, migrate.WithLogger(log))
	require.NoError(t, err)

	// Require baseline-version or explicit flag to work on a dirty workspace.
	files, err := ex.Pending(context.Background())
	require.EqualError(t, err, "sql/migrate: connected database is not clean: found table. baseline version or allow-dirty are required")
	require.Nil(t, files)

	rrw = mockRevisionReadWriter{}
	ex, err = migrate.NewExecutor(drv, dir, &rrw, migrate.WithLogger(log), migrate.WithAllowDirty(true))
	require.NoError(t, err)
	files, err = ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, files, 3)

	rrw = mockRevisionReadWriter{}
	ex, err = migrate.NewExecutor(drv, dir, &rrw, migrate.WithLogger(log), migrate.WithBaselineVersion("2.10.x-20"))
	require.NoError(t, err)
	files, err = ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, files, 1)
	require.Len(t, rrw, 1)
	require.Equal(t, "2.10.x-20", rrw[0].Version)
	require.Equal(t, "description", rrw[0].Description)
	require.Equal(t, migrate.RevisionTypeBaseline, rrw[0].Type)

	rrw = mockRevisionReadWriter{}
	ex, err = migrate.NewExecutor(drv, dir, &rrw, migrate.WithLogger(log), migrate.WithBaselineVersion("3"))
	require.NoError(t, err)
	files, err = ex.Pending(context.Background())
	require.ErrorIs(t, err, migrate.ErrNoPendingFiles)
	require.Len(t, rrw, 1)
	require.Equal(t, "3", rrw[0].Version)
	require.Equal(t, "partly", rrw[0].Description)
	require.Equal(t, migrate.RevisionTypeBaseline, rrw[0].Type)
}

func TestExecutor_FromVersion(t *testing.T) {
	var (
		drv = &lockMockDriver{&mockDriver{}}
		log = &mockLogger{}
		rrw = &mockRevisionReadWriter{
			{
				Version:     "1.a",
				Description: "sub.up",
				Applied:     2,
				Total:       2,
				Hash:        "nXyZR020M/mH7LxkoTkJr7BcQkipVg90imQ9I4595dw=",
			},
		}
	)
	dir, err := migrate.NewLocalDir(filepath.Join("testdata/migrate", "sub"))
	require.NoError(t, err)
	ex, err := migrate.NewExecutor(drv, dir, rrw, migrate.WithLogger(log))
	require.NoError(t, err)
	files, err := ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, files, 2)

	// Control the starting point.
	ex, err = migrate.NewExecutor(drv, dir, rrw, migrate.WithLogger(log), migrate.WithFromVersion("3"))
	require.NoError(t, err)
	files, err = ex.Pending(context.Background())
	require.NoError(t, err)
	require.Len(t, files, 1)

	// Starting point was not found.
	ex, err = migrate.NewExecutor(drv, dir, rrw, migrate.WithLogger(log), migrate.WithFromVersion("4"))
	require.NoError(t, err)
	files, err = ex.Pending(context.Background())
	require.EqualError(t, err, `starting point version "4" was not found in the migration directory`)
	require.Nil(t, files)
}

type (
	mockDriver struct {
		migrate.Driver
		plan          *migrate.Plan
		changes       []schema.Change
		applied       []schema.Change
		realm         schema.Realm
		executed      []string
		locks         map[string]struct{}
		lockCounter   int
		unlockCounter int
		failCounter   int
		failWith      error
		dirty         bool
	}
	lockMockDriver struct{ *mockDriver }
)

// the nth call to ExecContext will fail with the given error.
func (m *mockDriver) failOn(n int, err error) {
	m.failCounter = n
	m.failWith = err
}

func (m *mockDriver) ExecContext(_ context.Context, query string, _ ...interface{}) (sql.Result, error) {
	if m.failCounter > 0 {
		m.failCounter--
		if m.failCounter == 0 {
			return nil, m.failWith
		}
	}
	m.executed = append(m.executed, query)
	return nil, nil
}

func (m *mockDriver) InspectRealm(context.Context, *schema.InspectRealmOption) (*schema.Realm, error) {
	return &m.realm, nil
}

func (m *mockDriver) SchemaDiff(_, _ *schema.Schema) ([]schema.Change, error) {
	return m.changes, nil
}

func (m *mockDriver) RealmDiff(_, _ *schema.Realm) ([]schema.Change, error) {
	return m.changes, nil
}

func (m *mockDriver) PlanChanges(context.Context, string, []schema.Change) (*migrate.Plan, error) {
	return m.plan, nil
}

func (m *mockDriver) ApplyChanges(_ context.Context, changes []schema.Change) error {
	m.applied = changes
	return nil
}

func (m *mockDriver) Snapshot(context.Context) (migrate.RestoreFunc, error) {
	if m.dirty {
		return nil, migrate.NotCleanError{}
	}
	realm := m.realm
	return func(context.Context) error {
		m.realm = realm
		return nil
	}, nil
}

func (m *mockDriver) CheckClean(context.Context, *migrate.TableIdent) error {
	if m.dirty {
		return &migrate.NotCleanError{Reason: "found table"}
	}
	return nil
}

func (m *lockMockDriver) Lock(_ context.Context, name string, _ time.Duration) (schema.UnlockFunc, error) {
	if _, ok := m.locks[name]; ok {
		return nil, errors.New("lockErr")
	}
	m.lockCounter++
	if m.locks == nil {
		m.locks = make(map[string]struct{})
	}
	m.locks[name] = struct{}{}
	return func() error {
		m.unlockCounter++
		delete(m.locks, name)
		return nil
	}, nil
}

func (m *lockMockDriver) released() bool {
	return len(m.locks) == 0
}

type mockRevisionReadWriter migrate.Revisions

func (mockRevisionReadWriter) Ident() *migrate.TableIdent {
	return nil
}

func (mockRevisionReadWriter) Exists(_ context.Context) (bool, error) {
	return true, nil
}

func (mockRevisionReadWriter) Init(_ context.Context) error {
	return nil
}

func (rrw *mockRevisionReadWriter) WriteRevision(_ context.Context, r *migrate.Revision) error {
	for i, rev := range *rrw {
		if rev.Version == r.Version {
			(*rrw)[i] = r
			return nil
		}
	}
	*rrw = append(*rrw, r)
	return nil
}

func (rrw *mockRevisionReadWriter) ReadRevision(_ context.Context, v string) (*migrate.Revision, error) {
	for _, r := range migrate.Revisions(*rrw) {
		if r.Version == v {
			return r, nil
		}
	}
	return nil, migrate.ErrRevisionNotExist
}

func (rrw *mockRevisionReadWriter) ReadRevisions(context.Context) (migrate.Revisions, error) {
	return migrate.Revisions(*rrw), nil
}

func (rrw *mockRevisionReadWriter) clean() {
	*rrw = mockRevisionReadWriter(migrate.Revisions{})
}

type mockLogger []migrate.LogEntry

func (m *mockLogger) Log(e migrate.LogEntry) { *m = append(*m, e) }

func requireEqualRevisions(t *testing.T, expected, actual migrate.Revisions) {
	require.Equal(t, len(expected), len(actual))
	for i := range expected {
		requireEqualRevision(t, expected[i], actual[i])
	}
}

func requireEqualRevision(t *testing.T, expected, actual *migrate.Revision) {
	require.Equal(t, expected.Version, actual.Version)
	require.Equal(t, expected.Description, actual.Description)
	require.Equal(t, expected.Type, actual.Type)
	require.Equal(t, expected.Applied, actual.Applied)
	require.Equal(t, expected.Total, actual.Total)
	require.Equal(t, expected.Error, actual.Error)
	if expected.Hash != "" {
		require.Equal(t, expected.Hash, actual.Hash)
	}
	require.Equal(t, expected.OperatorVersion, actual.OperatorVersion)
}

func countFiles(t *testing.T, d migrate.Dir) int {
	files, err := fs.ReadDir(d, "")
	require.NoError(t, err)
	return len(files)
}

func requireFileEqual(t *testing.T, d migrate.Dir, name, contents string) {
	c, err := fs.ReadFile(d, name)
	require.NoError(t, err)
	require.Equal(t, contents, string(c))
}
