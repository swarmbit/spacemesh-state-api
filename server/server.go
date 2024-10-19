package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/config"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/price"
	"github.com/swarmbit/spacemesh-state-api/route"
	"github.com/swarmbit/spacemesh-state-api/sync/processor"
)

func StartServer(configValues *config.Config) {

	connection := configValues.DB.Uri
	writeDB, err := database.NewWriteDB(connection)
	if err != nil {
		panic("Failed to open document write db")
	}
	readDB, err := database.NewReadDB(connection)
	if err != nil {
		panic("Failed to open document read db")
	}
	log.Println("Created dbs")

	priceResolver := price.NewPriceResolver(configValues)
	log.Println("Created price resolver")

	if configValues.Sync.Enabled {
		p := processor.NewSyncProcessor(configValues, writeDB)
		p.StartProcessingSync()
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	route.AddRoutes(readDB, router, priceResolver, configValues)

	server := &http.Server{
		Addr:    configValues.Server.Port,
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

	log.Println("Listen and serve")
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server closed under request")
		} else {
			log.Fatal("Server closed unexpect")
		}
	}

	log.Println("Server exiting")
}
