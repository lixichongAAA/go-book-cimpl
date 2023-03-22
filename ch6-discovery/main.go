package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"ch6-discovery/config"
	"ch6-discovery/discover"
	"ch6-discovery/service"
	"ch6-discovery/transport"

	"github.com/longjoy/micro-go-book/ch6-discovery/endpoint"

	uuid "github.com/satori/go.uuid"
)

func main() {

	//从命令行读取参数,没有则使用默认值
	var (
		// 服务地址和服务名
		servicePort = flag.Int("service.port", 10086, "service port")
		serviceHost = flag.String("service.host", "127.0.0.1", "service host")
		serviceName = flag.String("service.name", "SayHello", "service name")
		// consul 地址
		consulPort = flag.Int("consul.port", 8500, "consul port")
		consulHost = flag.String("consul.host", "127.0.0.1", "service host")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	//声明服务发现客户端
	var discoveryClient discover.DiscoveryClient
	//直接使用HTTP与consul交互
	//discoveryClient, err := discover.NewMyDiscoverClient(*consulHost, *consulPort)
	//借助Go-kit服务注册与发现包和Consul交互
	discoveryClient, err := discover.NewKitDiscoverClient(*consulHost, *consulPort)
	//获取服务发现客户端失败,直接关闭服务
	if err != nil {
		config.Logger.Println("Get Consul Client failed.")
		os.Exit(-1)
	}

	//声明并初始化Service
	var svc = service.NewDiscoveryServiceImpl(discoveryClient)
	//创建打招呼的Endpoint
	sayHelloEndpoint := endpoint.MakeSayHelloEndpoint(svc)
	//创建服务发现的Endpoint
	discoveryEndpoint := endpoint.MakeDiscoveryEndpoint(svc)
	//创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	endpts := endpoint.DiscoveryEndpoints{
		SayHelloEndpoint:    sayHelloEndpoint,
		DiscoveryEndpoint:   discoveryEndpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.KitLogger)
	//定义服务实例ID
	instanceID := *serviceName + "-" + uuid.NewV4().String()
	//启动http Server
	go func() {
		config.Logger.Println("Http server start at port: ", strconv.Itoa(*servicePort))
		//启动前注册
		if !discoveryClient.Register(*serviceName, instanceID, "/health", *serviceHost, *servicePort, nil, config.Logger) {
			config.Logger.Printf("string-service for service %s failed.", *serviceName)
			os.Exit(-1)
		}
		handler := r
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), handler)
	}()

	go func() {
		//监控系统信号,等待Ctrl+C 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	discoveryClient.DeRegister(instanceID, config.Logger)
	config.Logger.Println(error)
}
