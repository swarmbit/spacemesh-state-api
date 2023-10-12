package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/price"
	"github.com/swarmbit/spacemesh-state-api/route"
	"github.com/swarmbit/spacemesh-state-api/sink"
)

func StartServer() {

	connection := "mongodb://spacemesh:<>@spacemesh-mongodb-svc.spacemesh.svc.cluster.local:27017/admin?replicaSet=spacemesh-mongodb&authMechanism=SCRAM-SHA-256"
	writeDB, err := database.NewWriteDB(connection)
	if err != nil {
		panic("Failed to open document write db")
	}
	readDB, err := database.NewReadDB(connection)
	if err != nil {
		panic("Failed to open document read db")
	}

	priceResolver := price.NewPriceResolver()

	sink := sink.NewSink(writeDB)
	sink.StartRewardsSink()
	sink.StartLayersSink()
	sink.StartAtxSink()
	sink.StartTransactionCreatedSink()
	sink.StartTransactionResultSink()
	sink.StartMalfeasanceSink()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}

	router.Use(cors.New(config))
	route.AddRoutes(readDB, router, priceResolver)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		writeDB.CloseWrite()
		readDB.CloseRead()
		log.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			log.Fatal("Server Close:", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}

	log.Println("Server exiting")
}
