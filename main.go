package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Prameesh-P/user-handler/grpcs"
	"github.com/Prameesh-P/user-handler/http"
	"github.com/Prameesh-P/user-handler/internal/database"
	 "github.com/Prameesh-P/user-handler/internal/initializers"
	"github.com/Prameesh-P/user-handler/internal/user"
	"github.com/Prameesh-P/user-handler/pkg/pb"
)

func init(){

	initializers.LoadEnvVariable()
	

}

func main() {
	pgDSN := os.Getenv("CONNSTR")
    pgDB, err := database.NewPostgresDB(pgDSN)
    if err != nil {
        log.Fatal("Error connecting to PostgreSQL:", err)
        return
    }
	database.SycnDB()

    // Initialize Redis connection
    redisAddr := "localhost:6379"
	redisPassword:= ""
    redisClient, err := database.NewRedisClient(redisAddr, redisPassword)
    if err != nil {
        fmt.Println("Error connecting to Redis:", err)
        return
    }
	userRepository := user.NewUserRepository(pgDB,redisClient)
	userUsecase := user.NewUserUsecase(userRepository)
	// Start the gRPC server in a separate goroutine
	go grpcs.StartGRPCServer(userUsecase)

	// Create a gRPC connection for the HTTP server
	grpcConn := grpcs.CreateGRPCConnection()
	defer grpcConn.Close()

	// Create the gRPC client
	userClient := pb.NewUserServiceClient(grpcConn)

	// Start the HTTP server
	httpServer := http.NewHTTPServer(userClient)
	go httpServer.StartHTTPServer()

	// Wait for a termination signal to gracefully shut down the servers
	waitForTerminationSignal()
}

func waitForTerminationSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Received termination signal. Shutting down gracefully...")
}