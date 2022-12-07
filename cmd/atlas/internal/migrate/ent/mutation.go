// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"ariga.io/atlas/cmd/atlas/internal/migrate/ent/predicate"
	"ariga.io/atlas/cmd/atlas/internal/migrate/ent/revision"
	"ariga.io/atlas/sql/migrate"

	"entgo.io/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeRevision = "Revision"
)

// RevisionMutation represents an operation that mutates the Revision nodes in the graph.
type RevisionMutation struct {
	config
	op                   Op
	typ                  string
	id                   *string
	description          *string
	_type                *migrate.RevisionType
	add_type             *migrate.RevisionType
	applied              *int
	addapplied           *int
	total                *int
	addtotal             *int
	executed_at          *time.Time
	execution_time       *time.Duration
	addexecution_time    *time.Duration
	error                *string
	error_stmt           *string
	hash                 *string
	partial_hashes       *[]string
	appendpartial_hashes []string
	operator_version     *string
	clearedFields        map[string]struct{}
	done                 bool
	oldValue             func(context.Context) (*Revision, error)
	predicates           []predicate.Revision
}

var _ ent.Mutation = (*RevisionMutation)(nil)

// revisionOption allows management of the mutation configuration using functional options.
type revisionOption func(*RevisionMutation)

// newRevisionMutation creates new mutation for the Revision entity.
func newRevisionMutation(c config, op Op, opts ...revisionOption) *RevisionMutation {
	m := &RevisionMutation{
		config:        c,
		op:            op,
		typ:           TypeRevision,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withRevisionID sets the ID field of the mutation.
func withRevisionID(id string) revisionOption {
	return func(m *RevisionMutation) {
		var (
			err   error
			once  sync.Once
			value *Revision
		)
		m.oldValue = func(ctx context.Context) (*Revision, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Revision.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withRevision sets the old Revision of the mutation.
func withRevision(node *Revision) revisionOption {
	return func(m *RevisionMutation) {
		m.oldValue = func(context.Context) (*Revision, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m RevisionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m RevisionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Revision entities.
func (m *RevisionMutation) SetID(id string) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *RevisionMutation) ID() (id string, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *RevisionMutation) IDs(ctx context.Context) ([]string, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []string{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Revision.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetDescription sets the "description" field.
func (m *RevisionMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the value of the "description" field in the mutation.
func (m *RevisionMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// OldDescription returns the old "description" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldDescription(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldDescription is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldDescription requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldDescription: %w", err)
	}
	return oldValue.Description, nil
}

// ResetDescription resets all changes to the "description" field.
func (m *RevisionMutation) ResetDescription() {
	m.description = nil
}

// SetType sets the "type" field.
func (m *RevisionMutation) SetType(mt migrate.RevisionType) {
	m._type = &mt
	m.add_type = nil
}

// GetType returns the value of the "type" field in the mutation.
func (m *RevisionMutation) GetType() (r migrate.RevisionType, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// OldType returns the old "type" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldType(ctx context.Context) (v migrate.RevisionType, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldType is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldType requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldType: %w", err)
	}
	return oldValue.Type, nil
}

// AddType adds mt to the "type" field.
func (m *RevisionMutation) AddType(mt migrate.RevisionType) {
	if m.add_type != nil {
		*m.add_type += mt
	} else {
		m.add_type = &mt
	}
}

// AddedType returns the value that was added to the "type" field in this mutation.
func (m *RevisionMutation) AddedType() (r migrate.RevisionType, exists bool) {
	v := m.add_type
	if v == nil {
		return
	}
	return *v, true
}

// ResetType resets all changes to the "type" field.
func (m *RevisionMutation) ResetType() {
	m._type = nil
	m.add_type = nil
}

// SetApplied sets the "applied" field.
func (m *RevisionMutation) SetApplied(i int) {
	m.applied = &i
	m.addapplied = nil
}

// Applied returns the value of the "applied" field in the mutation.
func (m *RevisionMutation) Applied() (r int, exists bool) {
	v := m.applied
	if v == nil {
		return
	}
	return *v, true
}

// OldApplied returns the old "applied" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldApplied(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldApplied is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldApplied requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldApplied: %w", err)
	}
	return oldValue.Applied, nil
}

// AddApplied adds i to the "applied" field.
func (m *RevisionMutation) AddApplied(i int) {
	if m.addapplied != nil {
		*m.addapplied += i
	} else {
		m.addapplied = &i
	}
}

// AddedApplied returns the value that was added to the "applied" field in this mutation.
func (m *RevisionMutation) AddedApplied() (r int, exists bool) {
	v := m.addapplied
	if v == nil {
		return
	}
	return *v, true
}

// ResetApplied resets all changes to the "applied" field.
func (m *RevisionMutation) ResetApplied() {
	m.applied = nil
	m.addapplied = nil
}

// SetTotal sets the "total" field.
func (m *RevisionMutation) SetTotal(i int) {
	m.total = &i
	m.addtotal = nil
}

// Total returns the value of the "total" field in the mutation.
func (m *RevisionMutation) Total() (r int, exists bool) {
	v := m.total
	if v == nil {
		return
	}
	return *v, true
}

// OldTotal returns the old "total" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldTotal(ctx context.Context) (v int, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTotal is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTotal requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTotal: %w", err)
	}
	return oldValue.Total, nil
}

// AddTotal adds i to the "total" field.
func (m *RevisionMutation) AddTotal(i int) {
	if m.addtotal != nil {
		*m.addtotal += i
	} else {
		m.addtotal = &i
	}
}

// AddedTotal returns the value that was added to the "total" field in this mutation.
func (m *RevisionMutation) AddedTotal() (r int, exists bool) {
	v := m.addtotal
	if v == nil {
		return
	}
	return *v, true
}

// ResetTotal resets all changes to the "total" field.
func (m *RevisionMutation) ResetTotal() {
	m.total = nil
	m.addtotal = nil
}

// SetExecutedAt sets the "executed_at" field.
func (m *RevisionMutation) SetExecutedAt(t time.Time) {
	m.executed_at = &t
}

// ExecutedAt returns the value of the "executed_at" field in the mutation.
func (m *RevisionMutation) ExecutedAt() (r time.Time, exists bool) {
	v := m.executed_at
	if v == nil {
		return
	}
	return *v, true
}

// OldExecutedAt returns the old "executed_at" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldExecutedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldExecutedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldExecutedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldExecutedAt: %w", err)
	}
	return oldValue.ExecutedAt, nil
}

// ResetExecutedAt resets all changes to the "executed_at" field.
func (m *RevisionMutation) ResetExecutedAt() {
	m.executed_at = nil
}

// SetExecutionTime sets the "execution_time" field.
func (m *RevisionMutation) SetExecutionTime(t time.Duration) {
	m.execution_time = &t
	m.addexecution_time = nil
}

// ExecutionTime returns the value of the "execution_time" field in the mutation.
func (m *RevisionMutation) ExecutionTime() (r time.Duration, exists bool) {
	v := m.execution_time
	if v == nil {
		return
	}
	return *v, true
}

// OldExecutionTime returns the old "execution_time" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldExecutionTime(ctx context.Context) (v time.Duration, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldExecutionTime is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldExecutionTime requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldExecutionTime: %w", err)
	}
	return oldValue.ExecutionTime, nil
}

// AddExecutionTime adds t to the "execution_time" field.
func (m *RevisionMutation) AddExecutionTime(t time.Duration) {
	if m.addexecution_time != nil {
		*m.addexecution_time += t
	} else {
		m.addexecution_time = &t
	}
}

// AddedExecutionTime returns the value that was added to the "execution_time" field in this mutation.
func (m *RevisionMutation) AddedExecutionTime() (r time.Duration, exists bool) {
	v := m.addexecution_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetExecutionTime resets all changes to the "execution_time" field.
func (m *RevisionMutation) ResetExecutionTime() {
	m.execution_time = nil
	m.addexecution_time = nil
}

// SetError sets the "error" field.
func (m *RevisionMutation) SetError(s string) {
	m.error = &s
}

// Error returns the value of the "error" field in the mutation.
func (m *RevisionMutation) Error() (r string, exists bool) {
	v := m.error
	if v == nil {
		return
	}
	return *v, true
}

// OldError returns the old "error" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldError(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldError is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldError requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldError: %w", err)
	}
	return oldValue.Error, nil
}

// ClearError clears the value of the "error" field.
func (m *RevisionMutation) ClearError() {
	m.error = nil
	m.clearedFields[revision.FieldError] = struct{}{}
}

// ErrorCleared returns if the "error" field was cleared in this mutation.
func (m *RevisionMutation) ErrorCleared() bool {
	_, ok := m.clearedFields[revision.FieldError]
	return ok
}

// ResetError resets all changes to the "error" field.
func (m *RevisionMutation) ResetError() {
	m.error = nil
	delete(m.clearedFields, revision.FieldError)
}

// SetErrorStmt sets the "error_stmt" field.
func (m *RevisionMutation) SetErrorStmt(s string) {
	m.error_stmt = &s
}

// ErrorStmt returns the value of the "error_stmt" field in the mutation.
func (m *RevisionMutation) ErrorStmt() (r string, exists bool) {
	v := m.error_stmt
	if v == nil {
		return
	}
	return *v, true
}

// OldErrorStmt returns the old "error_stmt" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldErrorStmt(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldErrorStmt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldErrorStmt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldErrorStmt: %w", err)
	}
	return oldValue.ErrorStmt, nil
}

// ClearErrorStmt clears the value of the "error_stmt" field.
func (m *RevisionMutation) ClearErrorStmt() {
	m.error_stmt = nil
	m.clearedFields[revision.FieldErrorStmt] = struct{}{}
}

// ErrorStmtCleared returns if the "error_stmt" field was cleared in this mutation.
func (m *RevisionMutation) ErrorStmtCleared() bool {
	_, ok := m.clearedFields[revision.FieldErrorStmt]
	return ok
}

// ResetErrorStmt resets all changes to the "error_stmt" field.
func (m *RevisionMutation) ResetErrorStmt() {
	m.error_stmt = nil
	delete(m.clearedFields, revision.FieldErrorStmt)
}

// SetHash sets the "hash" field.
func (m *RevisionMutation) SetHash(s string) {
	m.hash = &s
}

// Hash returns the value of the "hash" field in the mutation.
func (m *RevisionMutation) Hash() (r string, exists bool) {
	v := m.hash
	if v == nil {
		return
	}
	return *v, true
}

// OldHash returns the old "hash" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldHash(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldHash is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldHash requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldHash: %w", err)
	}
	return oldValue.Hash, nil
}

// ResetHash resets all changes to the "hash" field.
func (m *RevisionMutation) ResetHash() {
	m.hash = nil
}

// SetPartialHashes sets the "partial_hashes" field.
func (m *RevisionMutation) SetPartialHashes(s []string) {
	m.partial_hashes = &s
	m.appendpartial_hashes = nil
}

// PartialHashes returns the value of the "partial_hashes" field in the mutation.
func (m *RevisionMutation) PartialHashes() (r []string, exists bool) {
	v := m.partial_hashes
	if v == nil {
		return
	}
	return *v, true
}

// OldPartialHashes returns the old "partial_hashes" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldPartialHashes(ctx context.Context) (v []string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPartialHashes is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPartialHashes requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPartialHashes: %w", err)
	}
	return oldValue.PartialHashes, nil
}

// AppendPartialHashes adds s to the "partial_hashes" field.
func (m *RevisionMutation) AppendPartialHashes(s []string) {
	m.appendpartial_hashes = append(m.appendpartial_hashes, s...)
}

// AppendedPartialHashes returns the list of values that were appended to the "partial_hashes" field in this mutation.
func (m *RevisionMutation) AppendedPartialHashes() ([]string, bool) {
	if len(m.appendpartial_hashes) == 0 {
		return nil, false
	}
	return m.appendpartial_hashes, true
}

// ClearPartialHashes clears the value of the "partial_hashes" field.
func (m *RevisionMutation) ClearPartialHashes() {
	m.partial_hashes = nil
	m.appendpartial_hashes = nil
	m.clearedFields[revision.FieldPartialHashes] = struct{}{}
}

// PartialHashesCleared returns if the "partial_hashes" field was cleared in this mutation.
func (m *RevisionMutation) PartialHashesCleared() bool {
	_, ok := m.clearedFields[revision.FieldPartialHashes]
	return ok
}

// ResetPartialHashes resets all changes to the "partial_hashes" field.
func (m *RevisionMutation) ResetPartialHashes() {
	m.partial_hashes = nil
	m.appendpartial_hashes = nil
	delete(m.clearedFields, revision.FieldPartialHashes)
}

// SetOperatorVersion sets the "operator_version" field.
func (m *RevisionMutation) SetOperatorVersion(s string) {
	m.operator_version = &s
}

// OperatorVersion returns the value of the "operator_version" field in the mutation.
func (m *RevisionMutation) OperatorVersion() (r string, exists bool) {
	v := m.operator_version
	if v == nil {
		return
	}
	return *v, true
}

// OldOperatorVersion returns the old "operator_version" field's value of the Revision entity.
// If the Revision object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *RevisionMutation) OldOperatorVersion(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldOperatorVersion is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldOperatorVersion requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldOperatorVersion: %w", err)
	}
	return oldValue.OperatorVersion, nil
}

// ResetOperatorVersion resets all changes to the "operator_version" field.
func (m *RevisionMutation) ResetOperatorVersion() {
	m.operator_version = nil
}

// Where appends a list predicates to the RevisionMutation builder.
func (m *RevisionMutation) Where(ps ...predicate.Revision) {
	m.predicates = append(m.predicates, ps...)
}

// Op returns the operation name.
func (m *RevisionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Revision).
func (m *RevisionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *RevisionMutation) Fields() []string {
	fields := make([]string, 0, 11)
	if m.description != nil {
		fields = append(fields, revision.FieldDescription)
	}
	if m._type != nil {
		fields = append(fields, revision.FieldType)
	}
	if m.applied != nil {
		fields = append(fields, revision.FieldApplied)
	}
	if m.total != nil {
		fields = append(fields, revision.FieldTotal)
	}
	if m.executed_at != nil {
		fields = append(fields, revision.FieldExecutedAt)
	}
	if m.execution_time != nil {
		fields = append(fields, revision.FieldExecutionTime)
	}
	if m.error != nil {
		fields = append(fields, revision.FieldError)
	}
	if m.error_stmt != nil {
		fields = append(fields, revision.FieldErrorStmt)
	}
	if m.hash != nil {
		fields = append(fields, revision.FieldHash)
	}
	if m.partial_hashes != nil {
		fields = append(fields, revision.FieldPartialHashes)
	}
	if m.operator_version != nil {
		fields = append(fields, revision.FieldOperatorVersion)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *RevisionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case revision.FieldDescription:
		return m.Description()
	case revision.FieldType:
		return m.GetType()
	case revision.FieldApplied:
		return m.Applied()
	case revision.FieldTotal:
		return m.Total()
	case revision.FieldExecutedAt:
		return m.ExecutedAt()
	case revision.FieldExecutionTime:
		return m.ExecutionTime()
	case revision.FieldError:
		return m.Error()
	case revision.FieldErrorStmt:
		return m.ErrorStmt()
	case revision.FieldHash:
		return m.Hash()
	case revision.FieldPartialHashes:
		return m.PartialHashes()
	case revision.FieldOperatorVersion:
		return m.OperatorVersion()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *RevisionMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case revision.FieldDescription:
		return m.OldDescription(ctx)
	case revision.FieldType:
		return m.OldType(ctx)
	case revision.FieldApplied:
		return m.OldApplied(ctx)
	case revision.FieldTotal:
		return m.OldTotal(ctx)
	case revision.FieldExecutedAt:
		return m.OldExecutedAt(ctx)
	case revision.FieldExecutionTime:
		return m.OldExecutionTime(ctx)
	case revision.FieldError:
		return m.OldError(ctx)
	case revision.FieldErrorStmt:
		return m.OldErrorStmt(ctx)
	case revision.FieldHash:
		return m.OldHash(ctx)
	case revision.FieldPartialHashes:
		return m.OldPartialHashes(ctx)
	case revision.FieldOperatorVersion:
		return m.OldOperatorVersion(ctx)
	}
	return nil, fmt.Errorf("unknown Revision field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RevisionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case revision.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case revision.FieldType:
		v, ok := value.(migrate.RevisionType)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	case revision.FieldApplied:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetApplied(v)
		return nil
	case revision.FieldTotal:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTotal(v)
		return nil
	case revision.FieldExecutedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetExecutedAt(v)
		return nil
	case revision.FieldExecutionTime:
		v, ok := value.(time.Duration)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetExecutionTime(v)
		return nil
	case revision.FieldError:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetError(v)
		return nil
	case revision.FieldErrorStmt:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetErrorStmt(v)
		return nil
	case revision.FieldHash:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetHash(v)
		return nil
	case revision.FieldPartialHashes:
		v, ok := value.([]string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPartialHashes(v)
		return nil
	case revision.FieldOperatorVersion:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetOperatorVersion(v)
		return nil
	}
	return fmt.Errorf("unknown Revision field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *RevisionMutation) AddedFields() []string {
	var fields []string
	if m.add_type != nil {
		fields = append(fields, revision.FieldType)
	}
	if m.addapplied != nil {
		fields = append(fields, revision.FieldApplied)
	}
	if m.addtotal != nil {
		fields = append(fields, revision.FieldTotal)
	}
	if m.addexecution_time != nil {
		fields = append(fields, revision.FieldExecutionTime)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *RevisionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case revision.FieldType:
		return m.AddedType()
	case revision.FieldApplied:
		return m.AddedApplied()
	case revision.FieldTotal:
		return m.AddedTotal()
	case revision.FieldExecutionTime:
		return m.AddedExecutionTime()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *RevisionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case revision.FieldType:
		v, ok := value.(migrate.RevisionType)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddType(v)
		return nil
	case revision.FieldApplied:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddApplied(v)
		return nil
	case revision.FieldTotal:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddTotal(v)
		return nil
	case revision.FieldExecutionTime:
		v, ok := value.(time.Duration)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddExecutionTime(v)
		return nil
	}
	return fmt.Errorf("unknown Revision numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *RevisionMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(revision.FieldError) {
		fields = append(fields, revision.FieldError)
	}
	if m.FieldCleared(revision.FieldErrorStmt) {
		fields = append(fields, revision.FieldErrorStmt)
	}
	if m.FieldCleared(revision.FieldPartialHashes) {
		fields = append(fields, revision.FieldPartialHashes)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *RevisionMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *RevisionMutation) ClearField(name string) error {
	switch name {
	case revision.FieldError:
		m.ClearError()
		return nil
	case revision.FieldErrorStmt:
		m.ClearErrorStmt()
		return nil
	case revision.FieldPartialHashes:
		m.ClearPartialHashes()
		return nil
	}
	return fmt.Errorf("unknown Revision nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *RevisionMutation) ResetField(name string) error {
	switch name {
	case revision.FieldDescription:
		m.ResetDescription()
		return nil
	case revision.FieldType:
		m.ResetType()
		return nil
	case revision.FieldApplied:
		m.ResetApplied()
		return nil
	case revision.FieldTotal:
		m.ResetTotal()
		return nil
	case revision.FieldExecutedAt:
		m.ResetExecutedAt()
		return nil
	case revision.FieldExecutionTime:
		m.ResetExecutionTime()
		return nil
	case revision.FieldError:
		m.ResetError()
		return nil
	case revision.FieldErrorStmt:
		m.ResetErrorStmt()
		return nil
	case revision.FieldHash:
		m.ResetHash()
		return nil
	case revision.FieldPartialHashes:
		m.ResetPartialHashes()
		return nil
	case revision.FieldOperatorVersion:
		m.ResetOperatorVersion()
		return nil
	}
	return fmt.Errorf("unknown Revision field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *RevisionMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *RevisionMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *RevisionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *RevisionMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *RevisionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *RevisionMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *RevisionMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Revision unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *RevisionMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Revision edge %s", name)
}
