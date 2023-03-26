package main

// 使用 Hystrix 包装的反向代理
import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/longjoy/micro-go-book/ch6-discovery/discover"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"github.com/longjoy/micro-go-book/common/loadbalance"
)

type HystrixHandler struct {

	// 记录 Hystrix 是否已配置
	hystrixs     map[string]bool
	HystrixMutex *sync.Mutex

	discoverClient discover.DiscoveryClient
	loadbalance    loadbalance.LoadBalance
	logger         *log.Logger
}

func NewHystrixHandler(discoverClient discover.DiscoveryClient, loadbalance loadbalance.LoadBalance, logger *log.Logger) *HystrixHandler {
	return &HystrixHandler{
		discoverClient: discoverClient,
		logger:         logger,
		hystrixs:       make(map[string]bool),
		loadbalance:    loadbalance,
		HystrixMutex:   &sync.Mutex{},
	}
}

func (hystrixHandler *HystrixHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	if reqPath == "" {
		return
	}
	// 按照分隔符'/'对路经进行分割,获取服务名称 serviceName
	pathArray := strings.Split(reqPath, "/")
	serviceName := pathArray[1]

	if serviceName == "" {
		//路径不存在
		rw.WriteHeader(http.StatusNotFound) //404
		return
	}

	if _, ok := hystrixHandler.hystrixs[serviceName]; !ok {
		hystrixHandler.HystrixMutex.Lock()
		if _, ok := hystrixHandler.hystrixs[serviceName]; !ok {
			//把 serviceName 作为 hystrix 命令命名
			hystrix.ConfigureCommand(serviceName, hystrix.CommandConfig{
				//进行 hystrix 命令自定义
				RequestVolumeThreshold: 5,
			})
			hystrixHandler.hystrixs[serviceName] = true
		}
		hystrixHandler.HystrixMutex.Unlock()
	}

	err := hystrix.Do(serviceName, func() error {

		//调用 DiscoverClient 查询 serviceName 的服务实例列表
		instances := hystrixHandler.discoverClient.DiscoverServices(serviceName, hystrixHandler.logger)
		instanceList := make([]*api.AgentService, len(instances))
		for i := 0; i < len(instances); i++ {
			instanceList[i] = instances[i].(*api.AgentService)
		}
		// 使用负载均衡算法选取实例
		selectInstance, err := hystrixHandler.loadbalance.SelectService(instanceList)
		if err != nil {
			return errors.New("service instances are not existed.")
		}

		// 创建Director
		director := func(req *http.Request) {

			//重新组织请求路径, 去掉服务名称部分
			destPath := strings.Join(pathArray[2:], "/")

			hystrixHandler.logger.Println("service id ", selectInstance.ID)

			// 设置代理服务地址信息
			req.URL.Scheme = "http"
			req.URL.Host = fmt.Sprintf("%s:%d", selectInstance.Address, selectInstance.Port)
			req.URL.Path = "/" + destPath
		}

		var proxyError error

		// 返回代理异常, 用于记录 hystrix.Do 执行失败
		errorHandler := func(ew http.ResponseWriter, er *http.Request, err error) {
			proxyError = err
		}

		proxy := &httputil.ReverseProxy{
			Director:     director,
			ErrorHandler: errorHandler,
		}

		proxy.ServeHTTP(rw, req)

		// 将执行异常反馈 Hystrix
		return proxyError
	}, func(e error) error {
		hystrixHandler.logger.Println("proxy error ", e)
		return errors.New("fallback excute.")
	})

	// hystrix.Do 返回执行异常
	if err != nil {
		rw.WriteHeader(500)
		rw.Write([]byte(err.Error()))
	}
}
