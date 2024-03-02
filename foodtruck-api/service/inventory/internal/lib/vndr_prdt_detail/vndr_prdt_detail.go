package vndrprdtdetail

import (
	"context"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent"
	vpd "foodtruck/service/inventory/internal/lib/db/code_gen/ent/vendorproductdetail"
)

type VendorProductDetailService struct {
	client *ent.Client
	ctx    context.Context
}

func New(client *ent.Client) *VendorProductDetailService {
	return &VendorProductDetailService{
		client: client,
		ctx:    context.Background(),
	}
}

func (v *VendorProductDetailService) Create(vendorID int, productID int) (*ent.VendorProductDetail, error) {
	if _, err := v.client.Vendor.Get(v.ctx, vendorID); err != nil {
		return nil, err
	}

	if _, err := v.client.Product.Get(v.ctx, productID); err != nil {
		return nil, err
	}

	// query by FK
	// https://entgo.io/docs/schema-edges/#o2o-bidirectional
	vpdExists, err := v.client.VendorProductDetail.
		Query().
		Where(
			vpd.And(
				vpd.VendorID(vendorID),
				vpd.ProductID(vendorID),
			),
		).
		Exist(v.ctx)

	if err != nil {
		return nil, err
	}
	if vpdExists {
		return nil, nil
	}

	return v.client.VendorProductDetail.
		Create().
		SetProductID(productID).
		SetVendorID(vendorID).
		Save(v.ctx)

}
