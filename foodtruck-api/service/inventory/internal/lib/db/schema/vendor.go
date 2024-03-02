package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Vendor holds the schema definition for the Vendor entity.
type Vendor struct {
	ent.Schema
}

// Fields of the Vendor.
func (Vendor) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.String("description"),
	}
}

// Edges of the Vendor.
func (Vendor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vendor_product_details", VendorProductDetail.Type),
	}
}
