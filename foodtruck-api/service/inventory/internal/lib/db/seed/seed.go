package seed

import (
	"context"
	"fmt"
	dbconfig "foodtruck/pkg/db/db_config"
	debug "foodtruck/pkg/logger"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/product"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/vendor"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var logger = debug.Logger

// Opens a DB conn, creates the database, and then immediately closes the connections
func CreateDB(config dbconfig.Config) error {
	dbName := config.DbName
	config.DbName = "postgres"
	connstr := dbconfig.ConnString(dbconfig.Postgres, config)
	logger.Info(connstr)

	rawDbClient, err := sqlx.Connect("postgres", connstr)
	if err != nil {
		return errors.Wrap(err, "cannot create sqlx client")
	}

	if _, err := rawDbClient.Exec(fmt.Sprintf(`CREATE DATABASE %s`, dbName)); err != nil {
		logger.Info("Error creating database, database may already exist. ", err)
		return nil
	}

	// rawDbClient.MustExec(fmt.Sprintf(`CREATE DATABASE %s`, dbName))
	return rawDbClient.Close()
}

func InsertRows(config dbconfig.Config) error {
	connstr := dbconfig.ConnString(dbconfig.Postgres, config)
	client, err := ent.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
		return err
	}
	defer client.Close()

	logger.Info("conn successfully established")

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		return err
	}

	exist, err := client.Product.Query().Where(
		product.And(
			product.Name("product1"),
			product.Description("one"),
		)).Exist(ctx)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	var prdt *ent.Product
	prdt, err = client.Product.Create().SetName("product_1").SetDescription("one").Save(ctx)
	if err != nil {
		return err
	}
	if prdt == nil {
		return errors.New("product is nil")
	}

	if err := seedVendor(client, "vendor1", "one", prdt.ID); err != nil {
		return err
	}
	if err := seedVendor(client, "vendor2", "two", prdt.ID); err != nil {
		return err
	}
	return err
}

func seedVendor(client *ent.Client, name string, description string, productID int) (err error) {
	ctx := context.Background()

	exist, err := client.Vendor.Query().Where(
		vendor.And(
			vendor.Name(name),
			vendor.Description(description),
		)).Exist(ctx)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	var vendr *ent.Vendor
	vendr, err = client.Vendor.Create().
		SetName(name).
		SetDescription(description).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.VendorProductDetail.
		Create().
		SetProductID(productID).
		SetVendorID(vendr.ID).
		Save(ctx)
	return err
}
