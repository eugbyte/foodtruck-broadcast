// Code generated by ent, DO NOT EDIT.

package vendorproductdetail

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the vendorproductdetail type in the database.
	Label = "vendor_product_detail"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldVendorID holds the string denoting the vendor_id field in the database.
	FieldVendorID = "vendor_id"
	// FieldProductID holds the string denoting the product_id field in the database.
	FieldProductID = "product_id"
	// EdgeVendorParent holds the string denoting the vendor_parent edge name in mutations.
	EdgeVendorParent = "vendor_parent"
	// EdgeProductParent holds the string denoting the product_parent edge name in mutations.
	EdgeProductParent = "product_parent"
	// Table holds the table name of the vendorproductdetail in the database.
	Table = "vendor_product_details"
	// VendorParentTable is the table that holds the vendor_parent relation/edge.
	VendorParentTable = "vendor_product_details"
	// VendorParentInverseTable is the table name for the Vendor entity.
	// It exists in this package in order to avoid circular dependency with the "vendor" package.
	VendorParentInverseTable = "vendors"
	// VendorParentColumn is the table column denoting the vendor_parent relation/edge.
	VendorParentColumn = "vendor_id"
	// ProductParentTable is the table that holds the product_parent relation/edge.
	ProductParentTable = "vendor_product_details"
	// ProductParentInverseTable is the table name for the Product entity.
	// It exists in this package in order to avoid circular dependency with the "product" package.
	ProductParentInverseTable = "products"
	// ProductParentColumn is the table column denoting the product_parent relation/edge.
	ProductParentColumn = "product_id"
)

// Columns holds all SQL columns for vendorproductdetail fields.
var Columns = []string{
	FieldID,
	FieldVendorID,
	FieldProductID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the VendorProductDetail queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByVendorID orders the results by the vendor_id field.
func ByVendorID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVendorID, opts...).ToFunc()
}

// ByProductID orders the results by the product_id field.
func ByProductID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProductID, opts...).ToFunc()
}

// ByVendorParentField orders the results by vendor_parent field.
func ByVendorParentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newVendorParentStep(), sql.OrderByField(field, opts...))
	}
}

// ByProductParentField orders the results by product_parent field.
func ByProductParentField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProductParentStep(), sql.OrderByField(field, opts...))
	}
}
func newVendorParentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(VendorParentInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, VendorParentTable, VendorParentColumn),
	)
}
func newProductParentStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProductParentInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ProductParentTable, ProductParentColumn),
	)
}
