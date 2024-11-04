package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"mossT8.github.com/device-backend/internal/application/types"
	"mossT8.github.com/device-backend/internal/domain/customer"
	"mossT8.github.com/device-backend/internal/domain/device"
	"mossT8.github.com/device-backend/internal/infrastructure/config/aws"
	"mossT8.github.com/device-backend/internal/infrastructure/config/local"
	"mossT8.github.com/device-backend/internal/infrastructure/env"
	envConstants "mossT8.github.com/device-backend/internal/infrastructure/env/constants"
	"mossT8.github.com/device-backend/internal/infrastructure/logger"
	"mossT8.github.com/device-backend/internal/infrastructure/persistence/datastore"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http"
	httpConstants "mossT8.github.com/device-backend/internal/infrastructure/transport/http/constants"
	"mossT8.github.com/device-backend/internal/infrastructure/transport/http/middleware"
)

var sqlStoreConn *datastore.MySqlDataStore

var config *types.ConfigModel

var axxessLogs *accesslog.AccessLog

var customerDomain customer.CustomerDomain

var deviceDomain device.DeviceDomain

var irisServer *iris.Application

var port string

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Infof(httpConstants.DefaultRequestId, "Starting Application Version %s", envConstants.DefaultVersion)

	err := setup()
	if err != nil {
		logger.Errorf(httpConstants.DefaultRequestId, "unable to Setup Application: %s", err.Error())
		return
	}

	go start()

	<-ctx.Done()
	logger.Info(httpConstants.DefaultRequestId, "shutdown signalled...")
	shutdown()
	logger.Info(httpConstants.DefaultRequestId, "shutdown complete->")
}

func shutdown() {
	//Close Sql Conn to DB
	writerCloseError, readerCloseError := sqlStoreConn.Close()
	if writerCloseError != nil {
		logger.Errorf(httpConstants.CTXRequestIdKey, "Unable to close DB writer: %s", writerCloseError.Error())
	}
	if readerCloseError != nil {
		logger.Errorf(httpConstants.CTXRequestIdKey, "Unable to close DB reader: %s", writerCloseError.Error())
	}

	//Close Accesslogs to Server
	err := axxessLogs.Close()

	if err != nil {
		logger.Errorf(httpConstants.CTXRequestIdKey, "Unable to stop HTTP Server gracefully: %s", err.Error())
	}
}
func setup() error {
	var err error
	configManager := local.NewLocalConfigManager()
	logger.Infof(httpConstants.DefaultRequestId, "Env: %s", env.Getenv(os.Getenv(envConstants.Env), "[not specified]"))
	if os.Getenv(envConstants.Env) != "" {
		configManager = aws.NewSecretConfigManager()
	}

	logger.Infof(httpConstants.DefaultRequestId, "Secret Username: %s", env.Getenv(os.Getenv(envConstants.SecretName), "[not specified]"))
	config, err = configManager.GetConfig(os.Getenv("SECRET_NAME"))
	if err != nil {
		return fmt.Errorf("unable to load secret config: %s, exiting", err.Error())
	}

	sqlStoreConn = datastore.NewMysqlDataStore(config.Database)
	cErr := sqlStoreConn.Connect()
	if cErr != nil {
		return fmt.Errorf("unable to connect to db: %s, exiting", cErr.Error())
	}

	customerDomain = customer.NewCustomerDomain(sqlStoreConn)
	deviceDomain = device.NewDeviceDomain(sqlStoreConn)

	irisServer = iris.New()
	axxessLogs = middleware.MakeAccessLog()

	middlewareFunction := middleware.NewJWTMiddleware(nil)

	irisServer.Use(
		axxessLogs.Handler,
		middleware.CaselessMatcherMiddleware,
		middleware.RequestIDMiddleware,
		middlewareFunction([]string{"/login", "/logout", "/refresh", "/health"}),
	)

	http.NewAuthController(irisServer, customerDomain)
	http.NewCustomerController(sqlStoreConn, irisServer, customerDomain)
	http.NewDeviceController(sqlStoreConn, irisServer, deviceDomain, customerDomain)

	port = env.Getenv(envConstants.Port, envConstants.DefaultPort)

	return nil
}

func start() {
	if err := irisServer.Listen(fmt.Sprintf(":%s", port)); err != nil {
		logger.Errorf("failed to start server reason: %s", err.Error())
	}
}
