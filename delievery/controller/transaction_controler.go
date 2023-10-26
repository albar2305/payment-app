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

type TransactionController struct {
	router        *gin.Engine
	transactionUC usecase.TransactionUseCase
	maker         token.Maker
	cfg           *config.Config
}

func (t *TransactionController) createTransactionHandler(c *gin.Context) {
	var req model.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	authPayload := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)
	transactionRequest := model.CreateTransactionRequest{
		UserId:             authPayload.ID,
		ReceiverMerchantId: req.ReceiverMerchantId,
		Amount:             req.Amount,
	}

	user, err := t.transactionUC.RegisterNewTransaction(transactionRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (t *TransactionController) listTransactionHandler(c *gin.Context) {
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

	transactions, err := t.transactionUC.ListTransaction(arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (t *TransactionController) getTransactionHandlerByCustomerID(c *gin.Context) {
	id := c.Params.ByName("id")
	fmt.Println(id)

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

	transactions, err := t.transactionUC.GetTransactionByCustomerId(id, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func NewTransactionController(r *gin.Engine, usecase usecase.TransactionUseCase, cfg *config.Config) *TransactionController {
	tokenMaker, _ := token.NewJWTMaker(cfg.TokenSymetricKey)
	controller := TransactionController{
		router:        r,
		transactionUC: usecase,
		maker:         tokenMaker,
		cfg:           cfg,
	}

	rg := r.Group("/api/v1")
	rg.POST("/transactions", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.createTransactionHandler)
	rg.GET("/transactions/:id", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.getTransactionHandlerByCustomerID)
	rg.GET("/transactions", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.listTransactionHandler)
	return &controller
}
