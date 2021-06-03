// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"ariga.io/atlas/integration/entinteg/ent/defaultcontainer"
	"ariga.io/atlas/integration/entinteg/ent/predicate"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// DefaultContainerDelete is the builder for deleting a DefaultContainer entity.
type DefaultContainerDelete struct {
	config
	hooks    []Hook
	mutation *DefaultContainerMutation
}

// Where adds a new predicate to the DefaultContainerDelete builder.
func (dcd *DefaultContainerDelete) Where(ps ...predicate.DefaultContainer) *DefaultContainerDelete {
	dcd.mutation.predicates = append(dcd.mutation.predicates, ps...)
	return dcd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (dcd *DefaultContainerDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(dcd.hooks) == 0 {
		affected, err = dcd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*DefaultContainerMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			dcd.mutation = mutation
			affected, err = dcd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(dcd.hooks) - 1; i >= 0; i-- {
			mut = dcd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, dcd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (dcd *DefaultContainerDelete) ExecX(ctx context.Context) int {
	n, err := dcd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (dcd *DefaultContainerDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: defaultcontainer.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: defaultcontainer.FieldID,
			},
		},
	}
	if ps := dcd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, dcd.driver, _spec)
}

// DefaultContainerDeleteOne is the builder for deleting a single DefaultContainer entity.
type DefaultContainerDeleteOne struct {
	dcd *DefaultContainerDelete
}

// Exec executes the deletion query.
func (dcdo *DefaultContainerDeleteOne) Exec(ctx context.Context) error {
	n, err := dcdo.dcd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{defaultcontainer.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (dcdo *DefaultContainerDeleteOne) ExecX(ctx context.Context) {
	dcdo.dcd.ExecX(ctx)
}
