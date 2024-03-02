package main

import (
	"flag"
	dbconfig "foodtruck/pkg/db/db_config"
	vendorhandler "foodtruck/service/inventory/internal/handler/vendor"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent"
	"foodtruck/service/inventory/internal/lib/db/seed"
	"foodtruck/service/inventory/internal/lib/vendor"
	"log"

	debug "foodtruck/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var logger = debug.Logger

var addr = flag.String("addr", ":7000", "http service address")

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	config := dbconfig.Config{
		Host:     "host.docker.internal",
		Port:     "5200", // 5432
		User:     "postgres",
		Password: "postgres",
		DbName:   "inventory_staging",
		SSLMode:  "disable",
	}

	if err := seed.CreateDB(config); err != nil {
		log.Fatal(err)
	}
	logger.Info("successfully created DB")
	if err := seed.InsertRows(config); err != nil {
		log.Fatal(err)
	}
	logger.Info("successfully seeded")
	connstr := dbconfig.ConnString(dbconfig.Postgres, config)

	dbClient, err := ent.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer dbClient.Close()

	vendorService := vendor.New(dbClient)
	handler := vendorhandler.New(vendorService)

	logger.Info("listening to ", *addr)
	router.GET("/vendor/:id", func(c *gin.Context) {
		handler.GetVendor(c)
	})
	router.Run(*addr)
}
