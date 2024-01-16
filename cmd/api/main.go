package main

import (
	postProto "github.com/tom-blog-app/blog-proto/post"
	"github.com/tom-blog-app/post-service/pkg/service"
	"os"
	"sync"

	//"os"

	//"github.com/tom-blog-app/post-service/pkg/models"
	db "github.com/tom-blog-app/blog-utils/database"
	microservice "github.com/tom-blog-app/blog-utils/microservice"
	"google.golang.org/grpc"
	"log"
)

var (
	mongoURL = os.Getenv("POST_MONGO_URL")
	gRpcPort = os.Getenv("POST_GRPC_PORT")
)

func NewPostMicroApp() *microservice.MicroApp {

	client, err := db.ConnectToMongo(mongoURL)

	if err != nil {
		log.Panic(err)
	}

	return &microservice.MicroApp{
		GrpcServer: grpc.NewServer(),
		GrpcPort:   gRpcPort,
		RegisterService: func(server *grpc.Server) {
			postProto.RegisterPostServiceServer(server, &service.PostServer{
				Client: client,
			})
		},
	}
}

func main() {
	var wg sync.WaitGroup
	postApp := NewPostMicroApp()

	wg.Add(1)
	go func() {
		defer wg.Done()
		postApp.Register()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		postApp.CheckHealth()
	}()

	wg.Wait()
}
