package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swarmbit/spacemesh-state-api/database"
	"github.com/swarmbit/spacemesh-state-api/node"
	"github.com/swarmbit/spacemesh-state-api/route"
	"github.com/swarmbit/spacemesh-state-api/sink"
	"github.com/swarmbit/spacemesh-state-api/state"
)

func StartServer() {

	nodeDB, err := node.NewNodeDB("/Users/brunovale/spacemesh-db/node-data/state.sql", 1)
	if err != nil {
		panic("Failed to open node db")
	}

	state := state.NewState(nodeDB.DB, nodeDB.DB)

	docDB, err := database.NewDocDB("mongodb://localhost:27017")
	if err != nil {
		panic("Failed to open document db")
	}

	sink := sink.NewSink(docDB, nodeDB)

	go func() {
		sink.StartRewardsSink()
		time.Sleep(3 * time.Second)
		sink.StartLayersSink()
		time.Sleep(3 * time.Second)
		sink.StartAccountsSink()
	}()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	route.AddRoutes(state, router)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		nodeDB.Close()
		docDB.Close()
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
