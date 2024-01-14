package main

import (
	"fmt"
	postProto "github.com/tom-blog-app/blog-proto/post"
	"github.com/tom-blog-app/post-service/pkg/service"
	"google.golang.org/grpc/reflection"
	"os"
	"sync"

	//"os"

	//"github.com/tom-blog-app/post-service/pkg/models"
	db "github.com/tom-blog-app/blog-utils/database"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	mongoURL = os.Getenv("POST_MONGO_URL")
	gRpcPort = os.Getenv("POST_GRPC_PORT")
)

type PostServiceApp struct {
	server *grpc.Server
	client *mongo.Client
}

func NewPostServiceApp() *PostServiceApp {

	client, err := db.ConnectToMongo(mongoURL)

	if err != nil {
		log.Panic(err)
	}

	return &PostServiceApp{
		server: grpc.NewServer(),
		client: client,
	}
}

func main() {
	var wg sync.WaitGroup
	postApp := NewPostServiceApp()
	reflection.Register(postApp.server)

	wg.Add(1)
	go func() {
		defer wg.Done()
		postApp.register()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		postApp.checkHealth()
	}()

	wg.Wait()
}

func (app *PostServiceApp) register() {
	log.Println("Registering gRPC server..." + gRpcPort)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	postProto.RegisterPostServiceServer(app.server, &service.PostServer{
		Client: app.client,
	})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := app.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (app *PostServiceApp) checkHealth() {
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(app.server, healthServer)
}
