package vendor

import (
	"context"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent"
	vpd "foodtruck/service/inventory/internal/lib/db/code_gen/ent/vendorproductdetail"
)

type VendorService struct {
	client *ent.Client
	ctx    context.Context
}

func New(client *ent.Client) *VendorService {
	return &VendorService{
		client: client,
		ctx:    context.Background(),
	}
}

func (v *VendorService) Get(vendorID int) (*ent.Vendor, error) {
	return v.client.Vendor.Get(v.ctx, vendorID)
}

func (v *VendorService) Create(name string, description string) (*ent.Vendor, error) {
	return v.client.Vendor.Create().
		SetName(name).
		SetDescription(description).
		Save(v.ctx)
}

func (v *VendorService) AddProduct(vendorID int, productID int) (*ent.VendorProductDetail, error) {
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
