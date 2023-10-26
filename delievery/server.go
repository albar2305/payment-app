package delievery

import (
	"fmt"

	"github.com/albar2305/payment-app/config"
	"github.com/albar2305/payment-app/delievery/controller"
	"github.com/albar2305/payment-app/manager"
	"github.com/albar2305/payment-app/utils/exception"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	useCaseManager manager.UseCaseManager
	engine         *gin.Engine
	host           string
	log            *logrus.Logger
}

func (s *Server) Run() {
	s.setupControllers()
	err := s.engine.Run(s.host)
	if err != nil {
		panic(err)
	}
}

func (s *Server) setupControllers() {
	cfg, _ := config.NewConfig()
	controller.NewUserController(s.engine, s.useCaseManager.UserUseCase(), cfg)
	controller.NewCustomerController(s.engine, s.useCaseManager.CustomerUseCase(), cfg)
	controller.NewTransactionController(s.engine, s.useCaseManager.TransactionUseCase(), cfg)
	controller.NewMerchantController(s.engine, s.useCaseManager.MerchantUseCase(), cfg)
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	exception.CheckErr(err)
	infraManager, _ := manager.NewInfraManager(cfg)
	repoManager := manager.NewRepoManager(infraManager)
	useCaseManager := manager.NewUseCaseManager(repoManager)
	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)
	return &Server{
		useCaseManager: useCaseManager,
		engine:         engine,
		host:           host,
		log:            logrus.New(),
	}
}
