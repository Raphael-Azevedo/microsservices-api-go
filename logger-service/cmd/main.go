package main

import (
	"biz-hub-logger-service/configs"
	"biz-hub-logger-service/data"
	"biz-hub-logger-service/internal/web"
	"biz-hub-logger-service/internal/web/webserver"
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	mongoClient, err := connectToMong(configs)
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	data := data.New(client)

	err = rpc.Register(new(RPCServer))
	go rpcListen(configs)

	go gRPCListen(configs, data)

	var wg sync.WaitGroup
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webUserHandler := web.NewWebLoggerHandler(data)
	webserver.AddHandler("/log", webUserHandler.WriteLog)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	webserver.Start(&wg)

	wg.Wait()
}

func connectToMong(configs *configs.Conf) (*mongo.Client, error) {

	clientOptions := options.Client().ApplyURI(configs.MongUrl)
	clientOptions.SetAuth(options.Credential{
		Username: configs.DBName,
		Password: configs.DBPassword,
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Erro connecting:", err)
		return nil, err
	}
	log.Println("Conected to mongo!")
	return c, nil
}

func rpcListen(configs *configs.Conf) error {
	log.Println("Starting RPC server on port ", configs.RcpServerPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", configs.RcpServerPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
