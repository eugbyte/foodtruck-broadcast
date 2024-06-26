// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/vendor"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/vendorproductdetail"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// VendorCreate is the builder for creating a Vendor entity.
type VendorCreate struct {
	config
	mutation *VendorMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (vc *VendorCreate) SetName(s string) *VendorCreate {
	vc.mutation.SetName(s)
	return vc
}

// SetDescription sets the "description" field.
func (vc *VendorCreate) SetDescription(s string) *VendorCreate {
	vc.mutation.SetDescription(s)
	return vc
}

// AddVendorProductDetailIDs adds the "vendor_product_details" edge to the VendorProductDetail entity by IDs.
func (vc *VendorCreate) AddVendorProductDetailIDs(ids ...int) *VendorCreate {
	vc.mutation.AddVendorProductDetailIDs(ids...)
	return vc
}

// AddVendorProductDetails adds the "vendor_product_details" edges to the VendorProductDetail entity.
func (vc *VendorCreate) AddVendorProductDetails(v ...*VendorProductDetail) *VendorCreate {
	ids := make([]int, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return vc.AddVendorProductDetailIDs(ids...)
}

// Mutation returns the VendorMutation object of the builder.
func (vc *VendorCreate) Mutation() *VendorMutation {
	return vc.mutation
}

// Save creates the Vendor in the database.
func (vc *VendorCreate) Save(ctx context.Context) (*Vendor, error) {
	return withHooks(ctx, vc.sqlSave, vc.mutation, vc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (vc *VendorCreate) SaveX(ctx context.Context) *Vendor {
	v, err := vc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (vc *VendorCreate) Exec(ctx context.Context) error {
	_, err := vc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vc *VendorCreate) ExecX(ctx context.Context) {
	if err := vc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (vc *VendorCreate) check() error {
	if _, ok := vc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Vendor.name"`)}
	}
	if _, ok := vc.mutation.Description(); !ok {
		return &ValidationError{Name: "description", err: errors.New(`ent: missing required field "Vendor.description"`)}
	}
	return nil
}

func (vc *VendorCreate) sqlSave(ctx context.Context) (*Vendor, error) {
	if err := vc.check(); err != nil {
		return nil, err
	}
	_node, _spec := vc.createSpec()
	if err := sqlgraph.CreateNode(ctx, vc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	vc.mutation.id = &_node.ID
	vc.mutation.done = true
	return _node, nil
}

func (vc *VendorCreate) createSpec() (*Vendor, *sqlgraph.CreateSpec) {
	var (
		_node = &Vendor{config: vc.config}
		_spec = sqlgraph.NewCreateSpec(vendor.Table, sqlgraph.NewFieldSpec(vendor.FieldID, field.TypeInt))
	)
	if value, ok := vc.mutation.Name(); ok {
		_spec.SetField(vendor.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := vc.mutation.Description(); ok {
		_spec.SetField(vendor.FieldDescription, field.TypeString, value)
		_node.Description = value
	}
	if nodes := vc.mutation.VendorProductDetailsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   vendor.VendorProductDetailsTable,
			Columns: []string{vendor.VendorProductDetailsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(vendorproductdetail.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// VendorCreateBulk is the builder for creating many Vendor entities in bulk.
type VendorCreateBulk struct {
	config
	builders []*VendorCreate
}

// Save creates the Vendor entities in the database.
func (vcb *VendorCreateBulk) Save(ctx context.Context) ([]*Vendor, error) {
	specs := make([]*sqlgraph.CreateSpec, len(vcb.builders))
	nodes := make([]*Vendor, len(vcb.builders))
	mutators := make([]Mutator, len(vcb.builders))
	for i := range vcb.builders {
		func(i int, root context.Context) {
			builder := vcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*VendorMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, vcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, vcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, vcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (vcb *VendorCreateBulk) SaveX(ctx context.Context) []*Vendor {
	v, err := vcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (vcb *VendorCreateBulk) Exec(ctx context.Context) error {
	_, err := vcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vcb *VendorCreateBulk) ExecX(ctx context.Context) {
	if err := vcb.Exec(ctx); err != nil {
		panic(err)
	}
}
