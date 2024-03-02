package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// VendorProductDetail holds the schema definition for the VendorProductDetail entity.
type VendorProductDetail struct {
	ent.Schema
}

// Fields of the VendorProductDetail.
func (VendorProductDetail) Fields() []ent.Field {
	return []ent.Field{
		field.Int("vendor_id").
			Optional(),
		field.Int("product_id").
			Optional(),
	}
}

// Edges of the VendorProductDetail.
func (VendorProductDetail) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("vendor_parent", Vendor.Type).
			Ref("vendor_product_details").
			Unique().
			Field("vendor_id"),

		edge.From("product_parent", Product.Type).
			Ref("vendor_product_details").
			Unique().
			Field("product_id"),
	}
}
