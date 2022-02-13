package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sinojin/wallet/internal/transaction/adapters"
	"github.com/sinojin/wallet/internal/transaction/ports"
	"github.com/sinojin/wallet/internal/transaction/services"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func HttpServer(addr string) {
	// Creates a router without any middleware by default
	r := gin.New()
	// r.Use(gin.Logger())
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		log.WithFields(log.Fields{
			"ClientIP":   param.ClientIP,
			"Time":       param.TimeStamp.Format(time.RFC1123),
			"Method":     param.Method,
			"Path":       param.Path,
			"UserAgent":  param.Request.UserAgent(),
			"Latency":    param.Latency,
			"StatusCode": param.StatusCode,
		}).Info("")
		// your custom format
		return ""
	}))

	dsn := "docker:docker@tcp(127.0.0.1:3307)/db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	mysqlRepository := adapters.NewMysqlWalletRepository(db)
	redisRepository := adapters.NewRedisRepository(mysqlRepository, rdb)

	walletService := services.NewWalletService(redisRepository)

	transactionHandler := ports.NewTransactionHandler(walletService)

	r.Use(gin.Recovery())
	api := r.Group("/api")
	v1 := api.Group("/v1")
	transactionHandler.Handle(v1)
	r.Run(fmt.Sprintf(":%v", addr))

}
