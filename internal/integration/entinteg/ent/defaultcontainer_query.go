// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"ariga.io/atlas/integration/entinteg/ent/defaultcontainer"
	"ariga.io/atlas/integration/entinteg/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// DefaultContainerQuery is the builder for querying DefaultContainer entities.
type DefaultContainerQuery struct {
	config
	limit      *int
	offset     *int
	unique     *bool
	order      []OrderFunc
	fields     []string
	predicates []predicate.DefaultContainer
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the DefaultContainerQuery builder.
func (dcq *DefaultContainerQuery) Where(ps ...predicate.DefaultContainer) *DefaultContainerQuery {
	dcq.predicates = append(dcq.predicates, ps...)
	return dcq
}

// Limit adds a limit step to the query.
func (dcq *DefaultContainerQuery) Limit(limit int) *DefaultContainerQuery {
	dcq.limit = &limit
	return dcq
}

// Offset adds an offset step to the query.
func (dcq *DefaultContainerQuery) Offset(offset int) *DefaultContainerQuery {
	dcq.offset = &offset
	return dcq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (dcq *DefaultContainerQuery) Unique(unique bool) *DefaultContainerQuery {
	dcq.unique = &unique
	return dcq
}

// Order adds an order step to the query.
func (dcq *DefaultContainerQuery) Order(o ...OrderFunc) *DefaultContainerQuery {
	dcq.order = append(dcq.order, o...)
	return dcq
}

// First returns the first DefaultContainer entity from the query.
// Returns a *NotFoundError when no DefaultContainer was found.
func (dcq *DefaultContainerQuery) First(ctx context.Context) (*DefaultContainer, error) {
	nodes, err := dcq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{defaultcontainer.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (dcq *DefaultContainerQuery) FirstX(ctx context.Context) *DefaultContainer {
	node, err := dcq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first DefaultContainer ID from the query.
// Returns a *NotFoundError when no DefaultContainer ID was found.
func (dcq *DefaultContainerQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = dcq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{defaultcontainer.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (dcq *DefaultContainerQuery) FirstIDX(ctx context.Context) int {
	id, err := dcq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single DefaultContainer entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when exactly one DefaultContainer entity is not found.
// Returns a *NotFoundError when no DefaultContainer entities are found.
func (dcq *DefaultContainerQuery) Only(ctx context.Context) (*DefaultContainer, error) {
	nodes, err := dcq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{defaultcontainer.Label}
	default:
		return nil, &NotSingularError{defaultcontainer.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (dcq *DefaultContainerQuery) OnlyX(ctx context.Context) *DefaultContainer {
	node, err := dcq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only DefaultContainer ID in the query.
// Returns a *NotSingularError when exactly one DefaultContainer ID is not found.
// Returns a *NotFoundError when no entities are found.
func (dcq *DefaultContainerQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = dcq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{defaultcontainer.Label}
	default:
		err = &NotSingularError{defaultcontainer.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (dcq *DefaultContainerQuery) OnlyIDX(ctx context.Context) int {
	id, err := dcq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of DefaultContainers.
func (dcq *DefaultContainerQuery) All(ctx context.Context) ([]*DefaultContainer, error) {
	if err := dcq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return dcq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (dcq *DefaultContainerQuery) AllX(ctx context.Context) []*DefaultContainer {
	nodes, err := dcq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of DefaultContainer IDs.
func (dcq *DefaultContainerQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := dcq.Select(defaultcontainer.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (dcq *DefaultContainerQuery) IDsX(ctx context.Context) []int {
	ids, err := dcq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (dcq *DefaultContainerQuery) Count(ctx context.Context) (int, error) {
	if err := dcq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return dcq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (dcq *DefaultContainerQuery) CountX(ctx context.Context) int {
	count, err := dcq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (dcq *DefaultContainerQuery) Exist(ctx context.Context) (bool, error) {
	if err := dcq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return dcq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (dcq *DefaultContainerQuery) ExistX(ctx context.Context) bool {
	exist, err := dcq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the DefaultContainerQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (dcq *DefaultContainerQuery) Clone() *DefaultContainerQuery {
	if dcq == nil {
		return nil
	}
	return &DefaultContainerQuery{
		config:     dcq.config,
		limit:      dcq.limit,
		offset:     dcq.offset,
		order:      append([]OrderFunc{}, dcq.order...),
		predicates: append([]predicate.DefaultContainer{}, dcq.predicates...),
		// clone intermediate query.
		sql:  dcq.sql.Clone(),
		path: dcq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		String string `json:"string,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.DefaultContainer.Query().
//		GroupBy(defaultcontainer.FieldString).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (dcq *DefaultContainerQuery) GroupBy(field string, fields ...string) *DefaultContainerGroupBy {
	group := &DefaultContainerGroupBy{config: dcq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := dcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return dcq.sqlQuery(ctx), nil
	}
	return group
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		String string `json:"string,omitempty"`
//	}
//
//	client.DefaultContainer.Query().
//		Select(defaultcontainer.FieldString).
//		Scan(ctx, &v)
//
func (dcq *DefaultContainerQuery) Select(field string, fields ...string) *DefaultContainerSelect {
	dcq.fields = append([]string{field}, fields...)
	return &DefaultContainerSelect{DefaultContainerQuery: dcq}
}

func (dcq *DefaultContainerQuery) prepareQuery(ctx context.Context) error {
	for _, f := range dcq.fields {
		if !defaultcontainer.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if dcq.path != nil {
		prev, err := dcq.path(ctx)
		if err != nil {
			return err
		}
		dcq.sql = prev
	}
	return nil
}

func (dcq *DefaultContainerQuery) sqlAll(ctx context.Context) ([]*DefaultContainer, error) {
	var (
		nodes = []*DefaultContainer{}
		_spec = dcq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]interface{}, error) {
		node := &DefaultContainer{config: dcq.config}
		nodes = append(nodes, node)
		return node.scanValues(columns)
	}
	_spec.Assign = func(columns []string, values []interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(columns, values)
	}
	if err := sqlgraph.QueryNodes(ctx, dcq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (dcq *DefaultContainerQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := dcq.querySpec()
	return sqlgraph.CountNodes(ctx, dcq.driver, _spec)
}

func (dcq *DefaultContainerQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := dcq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %w", err)
	}
	return n > 0, nil
}

func (dcq *DefaultContainerQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   defaultcontainer.Table,
			Columns: defaultcontainer.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: defaultcontainer.FieldID,
			},
		},
		From:   dcq.sql,
		Unique: true,
	}
	if unique := dcq.unique; unique != nil {
		_spec.Unique = *unique
	}
	if fields := dcq.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, defaultcontainer.FieldID)
		for i := range fields {
			if fields[i] != defaultcontainer.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := dcq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := dcq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := dcq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := dcq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (dcq *DefaultContainerQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(dcq.driver.Dialect())
	t1 := builder.Table(defaultcontainer.Table)
	selector := builder.Select(t1.Columns(defaultcontainer.Columns...)...).From(t1)
	if dcq.sql != nil {
		selector = dcq.sql
		selector.Select(selector.Columns(defaultcontainer.Columns...)...)
	}
	for _, p := range dcq.predicates {
		p(selector)
	}
	for _, p := range dcq.order {
		p(selector)
	}
	if offset := dcq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := dcq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// DefaultContainerGroupBy is the group-by builder for DefaultContainer entities.
type DefaultContainerGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (dcgb *DefaultContainerGroupBy) Aggregate(fns ...AggregateFunc) *DefaultContainerGroupBy {
	dcgb.fns = append(dcgb.fns, fns...)
	return dcgb
}

// Scan applies the group-by query and scans the result into the given value.
func (dcgb *DefaultContainerGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := dcgb.path(ctx)
	if err != nil {
		return err
	}
	dcgb.sql = query
	return dcgb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (dcgb *DefaultContainerGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := dcgb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by.
// It is only allowed when executing a group-by query with one field.
func (dcgb *DefaultContainerGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(dcgb.fields) > 1 {
		return nil, errors.New("ent: DefaultContainerGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := dcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (dcgb *DefaultContainerGroupBy) StringsX(ctx context.Context) []string {
	v, err := dcgb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (dcgb *DefaultContainerGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = dcgb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{defaultcontainer.Label}
	default:
		err = fmt.Errorf("ent: DefaultContainerGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (dcgb *DefaultContainerGroupBy) StringX(ctx context.Context) string {
	v, err := dcgb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by.
// It is only allowed when executing a group-by query with one field.
func (dcgb *DefaultContainerGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(dcgb.fields) > 1 {
		return nil, errors.New("ent: DefaultContainerGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := dcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (dcgb *DefaultContainerGroupBy) IntsX(ctx context.Context) []int {
	v, err := dcgb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (dcgb *DefaultContainerGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = dcgb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{defaultcontainer.Label}
	default:
		err = fmt.Errorf("ent: DefaultContainerGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (dcgb *DefaultContainerGroupBy) IntX(ctx context.Context) int {
	v, err := dcgb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by.
// It is only allowed when executing a group-by query with one field.
func (dcgb *DefaultContainerGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(dcgb.fields) > 1 {
		return nil, errors.New("ent: DefaultContainerGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := dcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (dcgb *DefaultContainerGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := dcgb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (dcgb *DefaultContainerGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = dcgb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{defaultcontainer.Label}
	default:
		err = fmt.Errorf("ent: DefaultContainerGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (dcgb *DefaultContainerGroupBy) Float64X(ctx context.Context) float64 {
	v, err := dcgb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by.
// It is only allowed when executing a group-by query with one field.
func (dcgb *DefaultContainerGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(dcgb.fields) > 1 {
		return nil, errors.New("ent: DefaultContainerGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := dcgb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (dcgb *DefaultContainerGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := dcgb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a group-by query.
// It is only allowed when executing a group-by query with one field.
func (dcgb *DefaultContainerGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = dcgb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{defaultcontainer.Label}
	default:
		err = fmt.Errorf("ent: DefaultContainerGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (dcgb *DefaultContainerGroupBy) BoolX(ctx context.Context) bool {
	v, err := dcgb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (dcgb *DefaultContainerGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range dcgb.fields {
		if !defaultcontainer.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := dcgb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := dcgb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (dcgb *DefaultContainerGroupBy) sqlQuery() *sql.Selector {
	selector := dcgb.sql
	columns := make([]string, 0, len(dcgb.fields)+len(dcgb.fns))
	columns = append(columns, dcgb.fields...)
	for _, fn := range dcgb.fns {
		columns = append(columns, fn(selector))
	}
	return selector.Select(columns...).GroupBy(dcgb.fields...)
}

// DefaultContainerSelect is the builder for selecting fields of DefaultContainer entities.
type DefaultContainerSelect struct {
	*DefaultContainerQuery
	// intermediate query (i.e. traversal path).
	sql *sql.Selector
}

// Scan applies the selector query and scans the result into the given value.
func (dcs *DefaultContainerSelect) Scan(ctx context.Context, v interface{}) error {
	if err := dcs.prepareQuery(ctx); err != nil {
		return err
	}
	dcs.sql = dcs.DefaultContainerQuery.sqlQuery(ctx)
	return dcs.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (dcs *DefaultContainerSelect) ScanX(ctx context.Context, v interface{}) {
	if err := dcs.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from a selector. It is only allowed when selecting one field.
func (dcs *DefaultContainerSelect) Strings(ctx context.Context) ([]string, error) {
	if len(dcs.fields) > 1 {
		return nil, errors.New("ent: DefaultContainerSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := dcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (dcs *DefaultContainerSelect) StringsX(ctx context.Context) []string {
	v, err := dcs.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from a selector. It is only allowed when selecting one field.
func (dcs *DefaultContainerSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = dcs.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{defaultcontainer.Label}
	default:
		err = fmt.Errorf("ent: DefaultContainerSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (dcs *DefaultContainerSelect) StringX(ctx context.Context) string {
	v, err := dcs.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from a selector. It is only allowed when selecting one field.
func (dcs *DefaultContainerSelect) Ints(ctx context.Context) ([]int, error) {
	if len(dcs.fields) > 1 {
		return nil, errors.New("ent: DefaultContainerSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := dcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (dcs *DefaultContainerSelect) IntsX(ctx context.Context) []int {
	v, err := dcs.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from a selector. It is only allowed when selecting one field.
func (dcs *DefaultContainerSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = dcs.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{defaultcontainer.Label}
	default:
		err = fmt.Errorf("ent: DefaultContainerSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (dcs *DefaultContainerSelect) IntX(ctx context.Context) int {
	v, err := dcs.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from a selector. It is only allowed when selecting one field.
func (dcs *DefaultContainerSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(dcs.fields) > 1 {
		return nil, errors.New("ent: DefaultContainerSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := dcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (dcs *DefaultContainerSelect) Float64sX(ctx context.Context) []float64 {
	v, err := dcs.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from a selector. It is only allowed when selecting one field.
func (dcs *DefaultContainerSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = dcs.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{defaultcontainer.Label}
	default:
		err = fmt.Errorf("ent: DefaultContainerSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (dcs *DefaultContainerSelect) Float64X(ctx context.Context) float64 {
	v, err := dcs.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from a selector. It is only allowed when selecting one field.
func (dcs *DefaultContainerSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(dcs.fields) > 1 {
		return nil, errors.New("ent: DefaultContainerSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := dcs.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (dcs *DefaultContainerSelect) BoolsX(ctx context.Context) []bool {
	v, err := dcs.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from a selector. It is only allowed when selecting one field.
func (dcs *DefaultContainerSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = dcs.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{defaultcontainer.Label}
	default:
		err = fmt.Errorf("ent: DefaultContainerSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (dcs *DefaultContainerSelect) BoolX(ctx context.Context) bool {
	v, err := dcs.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (dcs *DefaultContainerSelect) sqlScan(ctx context.Context, v interface{}) error {
	rows := &sql.Rows{}
	query, args := dcs.sqlQuery().Query()
	if err := dcs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (dcs *DefaultContainerSelect) sqlQuery() sql.Querier {
	selector := dcs.sql
	selector.Select(selector.Columns(dcs.fields...)...)
	return selector
}
