// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/predicate"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/vendor"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// VendorDelete is the builder for deleting a Vendor entity.
type VendorDelete struct {
	config
	hooks    []Hook
	mutation *VendorMutation
}

// Where appends a list predicates to the VendorDelete builder.
func (vd *VendorDelete) Where(ps ...predicate.Vendor) *VendorDelete {
	vd.mutation.Where(ps...)
	return vd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (vd *VendorDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, vd.sqlExec, vd.mutation, vd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (vd *VendorDelete) ExecX(ctx context.Context) int {
	n, err := vd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (vd *VendorDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(vendor.Table, sqlgraph.NewFieldSpec(vendor.FieldID, field.TypeInt))
	if ps := vd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, vd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	vd.mutation.done = true
	return affected, err
}

// VendorDeleteOne is the builder for deleting a single Vendor entity.
type VendorDeleteOne struct {
	vd *VendorDelete
}

// Where appends a list predicates to the VendorDelete builder.
func (vdo *VendorDeleteOne) Where(ps ...predicate.Vendor) *VendorDeleteOne {
	vdo.vd.mutation.Where(ps...)
	return vdo
}

// Exec executes the deletion query.
func (vdo *VendorDeleteOne) Exec(ctx context.Context) error {
	n, err := vdo.vd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{vendor.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (vdo *VendorDeleteOne) ExecX(ctx context.Context) {
	if err := vdo.Exec(ctx); err != nil {
		panic(err)
	}
}