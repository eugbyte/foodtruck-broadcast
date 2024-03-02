package vendor

import (
	"context"
	"fmt"
	dbconfig "foodtruck/pkg/db/db_config"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent/enttest"
	"testing"

	assrt "github.com/blend/go-sdk/assert"
	_ "github.com/lib/pq"
)

func TestVendorService(t *testing.T) {
	assert := assrt.New(t)

	config := dbconfig.Config{
		Host:     "localhost", // for linux: "172.17.0.1"
		Port:     "5200",
		User:     "postgres",
		Password: "postgres",
		DbName:   "inventory_staging",
		SSLMode:  "disable",
	}

	connstr := dbconfig.ConnString(dbconfig.Postgres, config)
	client := enttest.Open(t, "postgres", connstr)
	defer client.Close()

	ctx := context.Background()

	err := client.Schema.Create(ctx)
	assert.Nil(err)
	fmt.Println("successfully opened and seeded")

	vendorService := New(client)

	v, err := vendorService.Create("vendor_1", "one")
	assert.Nil(err)

	_, err = vendorService.Get(v.ID)
	assert.Nil(err)
}
