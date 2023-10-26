package controller

import (
	"net/http"
	"strconv"

	"github.com/albar2305/payment-app/config"
	"github.com/albar2305/payment-app/delievery/middleware"
	"github.com/albar2305/payment-app/model"
	"github.com/albar2305/payment-app/usecase"
	"github.com/albar2305/payment-app/utils/common"
	"github.com/albar2305/payment-app/utils/token"
	"github.com/gin-gonic/gin"
)

type MerchantController struct {
	router     *gin.Engine
	merchantUC usecase.MerchantUseCase
	maker      token.Maker
	cfg        *config.Config
}

type getMerchantRequest struct {
	ID string `uri:"id" binding:"required"`
}

type deleteMerchantRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (u *MerchantController) createMerchantHandler(c *gin.Context) {
	var req model.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	merchant, err := u.merchantUC.RegisterNewMerchant(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, merchant)
}

func (u *MerchantController) deleteMerchantHandler(c *gin.Context) {
	var req deleteMerchantRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	err := u.merchantUC.DeleteMerchant(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func (u *MerchantController) listMerchantHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if page == 0 || limit == 0 {
		page = 1
		limit = 5
	}

	arg := model.PaginationParams{
		Limit:  int32(limit),
		Offset: int32((page - 1) * limit),
	}

	customers, err := u.merchantUC.ListMerchant(arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (u *MerchantController) getMerchantHandler(c *gin.Context) {
	var req getMerchantRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}
	merchant, err := u.merchantUC.GetMerchant(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, merchant)
}

func NewMerchantController(r *gin.Engine, usecase usecase.MerchantUseCase, cfg *config.Config) *MerchantController {
	tokenMaker, _ := token.NewJWTMaker(cfg.TokenSymetricKey)
	controller := MerchantController{
		router:     r,
		merchantUC: usecase,
		maker:      tokenMaker,
		cfg:        cfg,
	}

	rg := r.Group("/api/v1")
	rg.POST("/merchants", middleware.AuthMiddleware(tokenMaker, "admin"), controller.createMerchantHandler)
	rg.GET("/merchants", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.listMerchantHandler)
	rg.DELETE("/merchants/:id", middleware.AuthMiddleware(tokenMaker, "admin"), controller.deleteMerchantHandler)
	rg.GET("/merchants/:id", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.getMerchantHandler)
	return &controller
}
