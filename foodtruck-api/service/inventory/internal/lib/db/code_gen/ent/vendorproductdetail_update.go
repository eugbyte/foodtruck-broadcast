// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/predicate"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/product"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/vendor"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/vendorproductdetail"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// VendorProductDetailUpdate is the builder for updating VendorProductDetail entities.
type VendorProductDetailUpdate struct {
	config
	hooks    []Hook
	mutation *VendorProductDetailMutation
}

// Where appends a list predicates to the VendorProductDetailUpdate builder.
func (vpdu *VendorProductDetailUpdate) Where(ps ...predicate.VendorProductDetail) *VendorProductDetailUpdate {
	vpdu.mutation.Where(ps...)
	return vpdu
}

// SetVendorID sets the "vendor_id" field.
func (vpdu *VendorProductDetailUpdate) SetVendorID(i int) *VendorProductDetailUpdate {
	vpdu.mutation.SetVendorID(i)
	return vpdu
}

// SetNillableVendorID sets the "vendor_id" field if the given value is not nil.
func (vpdu *VendorProductDetailUpdate) SetNillableVendorID(i *int) *VendorProductDetailUpdate {
	if i != nil {
		vpdu.SetVendorID(*i)
	}
	return vpdu
}

// ClearVendorID clears the value of the "vendor_id" field.
func (vpdu *VendorProductDetailUpdate) ClearVendorID() *VendorProductDetailUpdate {
	vpdu.mutation.ClearVendorID()
	return vpdu
}

// SetProductID sets the "product_id" field.
func (vpdu *VendorProductDetailUpdate) SetProductID(i int) *VendorProductDetailUpdate {
	vpdu.mutation.SetProductID(i)
	return vpdu
}

// SetNillableProductID sets the "product_id" field if the given value is not nil.
func (vpdu *VendorProductDetailUpdate) SetNillableProductID(i *int) *VendorProductDetailUpdate {
	if i != nil {
		vpdu.SetProductID(*i)
	}
	return vpdu
}

// ClearProductID clears the value of the "product_id" field.
func (vpdu *VendorProductDetailUpdate) ClearProductID() *VendorProductDetailUpdate {
	vpdu.mutation.ClearProductID()
	return vpdu
}

// SetVendorParentID sets the "vendor_parent" edge to the Vendor entity by ID.
func (vpdu *VendorProductDetailUpdate) SetVendorParentID(id int) *VendorProductDetailUpdate {
	vpdu.mutation.SetVendorParentID(id)
	return vpdu
}

// SetNillableVendorParentID sets the "vendor_parent" edge to the Vendor entity by ID if the given value is not nil.
func (vpdu *VendorProductDetailUpdate) SetNillableVendorParentID(id *int) *VendorProductDetailUpdate {
	if id != nil {
		vpdu = vpdu.SetVendorParentID(*id)
	}
	return vpdu
}

// SetVendorParent sets the "vendor_parent" edge to the Vendor entity.
func (vpdu *VendorProductDetailUpdate) SetVendorParent(v *Vendor) *VendorProductDetailUpdate {
	return vpdu.SetVendorParentID(v.ID)
}

// SetProductParentID sets the "product_parent" edge to the Product entity by ID.
func (vpdu *VendorProductDetailUpdate) SetProductParentID(id int) *VendorProductDetailUpdate {
	vpdu.mutation.SetProductParentID(id)
	return vpdu
}

// SetNillableProductParentID sets the "product_parent" edge to the Product entity by ID if the given value is not nil.
func (vpdu *VendorProductDetailUpdate) SetNillableProductParentID(id *int) *VendorProductDetailUpdate {
	if id != nil {
		vpdu = vpdu.SetProductParentID(*id)
	}
	return vpdu
}

// SetProductParent sets the "product_parent" edge to the Product entity.
func (vpdu *VendorProductDetailUpdate) SetProductParent(p *Product) *VendorProductDetailUpdate {
	return vpdu.SetProductParentID(p.ID)
}

// Mutation returns the VendorProductDetailMutation object of the builder.
func (vpdu *VendorProductDetailUpdate) Mutation() *VendorProductDetailMutation {
	return vpdu.mutation
}

// ClearVendorParent clears the "vendor_parent" edge to the Vendor entity.
func (vpdu *VendorProductDetailUpdate) ClearVendorParent() *VendorProductDetailUpdate {
	vpdu.mutation.ClearVendorParent()
	return vpdu
}

// ClearProductParent clears the "product_parent" edge to the Product entity.
func (vpdu *VendorProductDetailUpdate) ClearProductParent() *VendorProductDetailUpdate {
	vpdu.mutation.ClearProductParent()
	return vpdu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (vpdu *VendorProductDetailUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, vpdu.sqlSave, vpdu.mutation, vpdu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (vpdu *VendorProductDetailUpdate) SaveX(ctx context.Context) int {
	affected, err := vpdu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (vpdu *VendorProductDetailUpdate) Exec(ctx context.Context) error {
	_, err := vpdu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vpdu *VendorProductDetailUpdate) ExecX(ctx context.Context) {
	if err := vpdu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (vpdu *VendorProductDetailUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(vendorproductdetail.Table, vendorproductdetail.Columns, sqlgraph.NewFieldSpec(vendorproductdetail.FieldID, field.TypeInt))
	if ps := vpdu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if vpdu.mutation.VendorParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   vendorproductdetail.VendorParentTable,
			Columns: []string{vendorproductdetail.VendorParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(vendor.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vpdu.mutation.VendorParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   vendorproductdetail.VendorParentTable,
			Columns: []string{vendorproductdetail.VendorParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(vendor.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if vpdu.mutation.ProductParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   vendorproductdetail.ProductParentTable,
			Columns: []string{vendorproductdetail.ProductParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(product.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vpdu.mutation.ProductParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   vendorproductdetail.ProductParentTable,
			Columns: []string{vendorproductdetail.ProductParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(product.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, vpdu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{vendorproductdetail.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	vpdu.mutation.done = true
	return n, nil
}

// VendorProductDetailUpdateOne is the builder for updating a single VendorProductDetail entity.
type VendorProductDetailUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *VendorProductDetailMutation
}

// SetVendorID sets the "vendor_id" field.
func (vpduo *VendorProductDetailUpdateOne) SetVendorID(i int) *VendorProductDetailUpdateOne {
	vpduo.mutation.SetVendorID(i)
	return vpduo
}

// SetNillableVendorID sets the "vendor_id" field if the given value is not nil.
func (vpduo *VendorProductDetailUpdateOne) SetNillableVendorID(i *int) *VendorProductDetailUpdateOne {
	if i != nil {
		vpduo.SetVendorID(*i)
	}
	return vpduo
}

// ClearVendorID clears the value of the "vendor_id" field.
func (vpduo *VendorProductDetailUpdateOne) ClearVendorID() *VendorProductDetailUpdateOne {
	vpduo.mutation.ClearVendorID()
	return vpduo
}

// SetProductID sets the "product_id" field.
func (vpduo *VendorProductDetailUpdateOne) SetProductID(i int) *VendorProductDetailUpdateOne {
	vpduo.mutation.SetProductID(i)
	return vpduo
}

// SetNillableProductID sets the "product_id" field if the given value is not nil.
func (vpduo *VendorProductDetailUpdateOne) SetNillableProductID(i *int) *VendorProductDetailUpdateOne {
	if i != nil {
		vpduo.SetProductID(*i)
	}
	return vpduo
}

// ClearProductID clears the value of the "product_id" field.
func (vpduo *VendorProductDetailUpdateOne) ClearProductID() *VendorProductDetailUpdateOne {
	vpduo.mutation.ClearProductID()
	return vpduo
}

// SetVendorParentID sets the "vendor_parent" edge to the Vendor entity by ID.
func (vpduo *VendorProductDetailUpdateOne) SetVendorParentID(id int) *VendorProductDetailUpdateOne {
	vpduo.mutation.SetVendorParentID(id)
	return vpduo
}

// SetNillableVendorParentID sets the "vendor_parent" edge to the Vendor entity by ID if the given value is not nil.
func (vpduo *VendorProductDetailUpdateOne) SetNillableVendorParentID(id *int) *VendorProductDetailUpdateOne {
	if id != nil {
		vpduo = vpduo.SetVendorParentID(*id)
	}
	return vpduo
}

// SetVendorParent sets the "vendor_parent" edge to the Vendor entity.
func (vpduo *VendorProductDetailUpdateOne) SetVendorParent(v *Vendor) *VendorProductDetailUpdateOne {
	return vpduo.SetVendorParentID(v.ID)
}

// SetProductParentID sets the "product_parent" edge to the Product entity by ID.
func (vpduo *VendorProductDetailUpdateOne) SetProductParentID(id int) *VendorProductDetailUpdateOne {
	vpduo.mutation.SetProductParentID(id)
	return vpduo
}

// SetNillableProductParentID sets the "product_parent" edge to the Product entity by ID if the given value is not nil.
func (vpduo *VendorProductDetailUpdateOne) SetNillableProductParentID(id *int) *VendorProductDetailUpdateOne {
	if id != nil {
		vpduo = vpduo.SetProductParentID(*id)
	}
	return vpduo
}

// SetProductParent sets the "product_parent" edge to the Product entity.
func (vpduo *VendorProductDetailUpdateOne) SetProductParent(p *Product) *VendorProductDetailUpdateOne {
	return vpduo.SetProductParentID(p.ID)
}

// Mutation returns the VendorProductDetailMutation object of the builder.
func (vpduo *VendorProductDetailUpdateOne) Mutation() *VendorProductDetailMutation {
	return vpduo.mutation
}

// ClearVendorParent clears the "vendor_parent" edge to the Vendor entity.
func (vpduo *VendorProductDetailUpdateOne) ClearVendorParent() *VendorProductDetailUpdateOne {
	vpduo.mutation.ClearVendorParent()
	return vpduo
}

// ClearProductParent clears the "product_parent" edge to the Product entity.
func (vpduo *VendorProductDetailUpdateOne) ClearProductParent() *VendorProductDetailUpdateOne {
	vpduo.mutation.ClearProductParent()
	return vpduo
}

// Where appends a list predicates to the VendorProductDetailUpdate builder.
func (vpduo *VendorProductDetailUpdateOne) Where(ps ...predicate.VendorProductDetail) *VendorProductDetailUpdateOne {
	vpduo.mutation.Where(ps...)
	return vpduo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (vpduo *VendorProductDetailUpdateOne) Select(field string, fields ...string) *VendorProductDetailUpdateOne {
	vpduo.fields = append([]string{field}, fields...)
	return vpduo
}

// Save executes the query and returns the updated VendorProductDetail entity.
func (vpduo *VendorProductDetailUpdateOne) Save(ctx context.Context) (*VendorProductDetail, error) {
	return withHooks(ctx, vpduo.sqlSave, vpduo.mutation, vpduo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (vpduo *VendorProductDetailUpdateOne) SaveX(ctx context.Context) *VendorProductDetail {
	node, err := vpduo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (vpduo *VendorProductDetailUpdateOne) Exec(ctx context.Context) error {
	_, err := vpduo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (vpduo *VendorProductDetailUpdateOne) ExecX(ctx context.Context) {
	if err := vpduo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (vpduo *VendorProductDetailUpdateOne) sqlSave(ctx context.Context) (_node *VendorProductDetail, err error) {
	_spec := sqlgraph.NewUpdateSpec(vendorproductdetail.Table, vendorproductdetail.Columns, sqlgraph.NewFieldSpec(vendorproductdetail.FieldID, field.TypeInt))
	id, ok := vpduo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "VendorProductDetail.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := vpduo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, vendorproductdetail.FieldID)
		for _, f := range fields {
			if !vendorproductdetail.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != vendorproductdetail.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := vpduo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if vpduo.mutation.VendorParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   vendorproductdetail.VendorParentTable,
			Columns: []string{vendorproductdetail.VendorParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(vendor.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vpduo.mutation.VendorParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   vendorproductdetail.VendorParentTable,
			Columns: []string{vendorproductdetail.VendorParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(vendor.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if vpduo.mutation.ProductParentCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   vendorproductdetail.ProductParentTable,
			Columns: []string{vendorproductdetail.ProductParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(product.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := vpduo.mutation.ProductParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   vendorproductdetail.ProductParentTable,
			Columns: []string{vendorproductdetail.ProductParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(product.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &VendorProductDetail{config: vpduo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, vpduo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{vendorproductdetail.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	vpduo.mutation.done = true
	return _node, nil
}
