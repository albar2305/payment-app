package controller

import (
	"fmt"
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

type CustomerController struct {
	router     *gin.Engine
	customerUC usecase.CustomerUseCase
	maker      token.Maker
	cfg        *config.Config
}

func (u *CustomerController) createCustomerHandler(c *gin.Context) {
	var req model.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	authPayload := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	arg := model.CreateCustomerRequest{
		UserID:  authPayload.ID,
		Name:    req.Name,
		Balance: 0,
	}

	fmt.Println(authPayload.ID)

	customer, err := u.customerUC.RegisterNewCustomer(arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, customer)
}

type getCustomerRequest struct {
	ID string `uri:"id" binding:"required"`
}

type deleteCustomerRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (u *CustomerController) getCustomerHandler(c *gin.Context) {
	var req getCustomerRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}
	customer, err := u.customerUC.GetCustomerById(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (u *CustomerController) listCustomerHandler(c *gin.Context) {
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

	customers, err := u.customerUC.ListCustomer(arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (u *CustomerController) deleteCustomerHandler(c *gin.Context) {
	var req deleteCustomerRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	err := u.customerUC.DeleteCustomer(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusNoContent, "")
}

func (u *CustomerController) addCustomerBalanceHandler(c *gin.Context) {
	var req model.AddCustomerBalanceParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	authPayload := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	customerResponse, err := u.customerUC.AddCustomerBalance(authPayload.ID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, customerResponse)
}

func NewCustomerController(r *gin.Engine, usecase usecase.CustomerUseCase, cfg *config.Config) *CustomerController {
	tokenMaker, _ := token.NewJWTMaker(cfg.TokenSymetricKey)
	controller := CustomerController{
		router:     r,
		customerUC: usecase,
		maker:      tokenMaker,
		cfg:        cfg,
	}

	rg := r.Group("/api/v1")
	rg.POST("/customers", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.createCustomerHandler)
	rg.GET("/customers/:name", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.getCustomerHandler)
	rg.DELETE("/customers/:id", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.deleteCustomerHandler)
	rg.GET("/customers", middleware.AuthMiddleware(tokenMaker, "admin"), controller.listCustomerHandler)
	rg.POST("/customers/top-up", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.addCustomerBalanceHandler)

	return &controller
}
