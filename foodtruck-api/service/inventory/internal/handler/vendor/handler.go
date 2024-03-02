package vendor

import (
	"fmt"
	debug "foodtruck/pkg/logger"
	"foodtruck/service/inventory/internal/lib/db/code_gen/ent"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var logger = debug.Logger

type VendorServicer interface {
	Get(vendorID int) (*ent.Vendor, error)
}

type Handler struct {
	vendorService VendorServicer
}

func New(vendorService VendorServicer) *Handler {
	return &Handler{
		vendorService: vendorService,
	}
}

func (h *Handler) GetVendor(c *gin.Context) {
	id := c.Param("id")
	vendorID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"msg": "invalid vendorID",
		})
		return
	}

	vendor, err := h.vendorService.Get(vendorID)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, map[string]string{
			"msg": fmt.Sprintf("vendorID %d not found", vendorID),
		})
		return
	}

	c.JSON(http.StatusOK, vendor)
}
