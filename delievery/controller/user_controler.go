package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/albar2305/payment-app/config"
	"github.com/albar2305/payment-app/delievery/middleware"
	"github.com/albar2305/payment-app/model"
	"github.com/albar2305/payment-app/usecase"
	"github.com/albar2305/payment-app/utils/common"
	"github.com/albar2305/payment-app/utils/token"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	router *gin.Engine
	userUC usecase.UserUseCase
	maker  token.Maker
	cfg    *config.Config
}

func newUserResponse(userRequest model.User) model.UserResponse {
	return model.UserResponse{
		ID:        userRequest.ID,
		Username:  userRequest.Username,
		Email:     userRequest.Email,
		CreatedAt: userRequest.CreatedAt,
	}
}

func (u *UserController) createUserHandler(c *gin.Context) {
	var userRequest model.CreateUserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	user, err := u.userUC.RegisterNewUser(userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, newUserResponse(user))
}

type getUserRequest struct {
	ID string `uri:"id" binding:"required"`
}

func (u *UserController) getUserHandler(c *gin.Context) {
	var req getUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}
	userRequest, err := u.userUC.GetUserById(req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, newUserResponse(userRequest))
}

func (u *UserController) updateUserHandler(c *gin.Context) {
	var userRequest model.User
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	user, err := u.userUC.UpdateUser(userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, newUserResponse(user))

}

func (u *UserController) listUserHandler(c *gin.Context) {
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

	users, err := u.userUC.ListUser(arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, users)
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken           string             `json:"access_token"`
	AccessTokenExpiresAt  time.Time          `json:"access_token_expires_at"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time          `json:"refresh_token_expires_at"`
	User                  model.UserResponse `json:"user"`
}

func (u *UserController) loginHandler(c *gin.Context) {
	var req loginUserRequest
	fmt.Println(req.Username)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse(err))
		return
	}

	user, err := u.userUC.GetUser(req.Username)
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.ErrorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	err = common.CheckPassword(req.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.ErrorResponse(err))
		return
	}

	accessToken, accessPayload, err := u.maker.CreateToken(
		user.ID,
		user.Username,
		user.Role,
		u.cfg.AccessTokenDuration,
	)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := u.maker.CreateToken(
		user.ID,
		user.Username,
		user.Role,
		u.cfg.RefreshTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	c.JSON(http.StatusOK, rsp)
}

func NewUserController(r *gin.Engine, usecase usecase.UserUseCase, cfg *config.Config) *UserController {
	tokenMaker, _ := token.NewJWTMaker(cfg.TokenSymetricKey)
	controller := UserController{
		router: r,
		userUC: usecase,
		maker:  tokenMaker,
		cfg:    cfg,
	}

	rg := r.Group("/api/v1")
	rg.POST("/users", controller.createUserHandler)
	rg.GET("/users", middleware.AuthMiddleware(tokenMaker, "admin"), controller.listUserHandler)
	rg.GET("/users"+"/:id", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.getUserHandler)
	rg.PUT("/users", middleware.AuthMiddleware(tokenMaker, "admin", "user"), controller.updateUserHandler)
	rg.POST("/users/login", controller.loginHandler)

	return &controller
}
