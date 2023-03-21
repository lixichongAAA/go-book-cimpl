package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	endpts "github.com/longjoy/micro-go-book/ch6-discovery/endpoint"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpts.DiscoveryEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	// 定义处理器
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeErr),
	}

	// say-hello 接口 create say-hello handler
	r.Methods("GET").Path("/say-hello").Handler(kithttp.NewServer(
		endpoints.SayHelloEndpoint,
		decodeSayHelloRequest,
		encodeJsonResponse,
		options...,
	))

	// discovery 接口 create discovery handler
	r.Methods("GET").Path("/discovery").Handler(kithttp.NewServer(
		endpoints.DiscoveryEndpoint,
		decodeDiscoveryRequest,
		encodeJsonResponse,
		options...,
	))

	// create HealthCheck handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthRequest,
		encodeJsonResponse,
		options...,
	))

	return r
}

// 解码http请求为endpts的SayHelloRequest结构体
func decodeSayHelloRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpts.SayHelloRequest{}, nil
}

// 解码http请求为 endpts的DiscoveryRequest结构体
func decodeDiscoveryRequest(_ context.Context, r *http.Request) (interface{}, error) {
	serviceName := r.URL.Query().Get("servicename")
	if serviceName == "" {
		return nil, ErrorBadRequest
	}
	return endpts.DiscoveryRequest{
		ServiceName: serviceName,
	}, nil
}

// 解码http请求为endpts的HealthRequest结构体
func decodeHealthRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}

// 编码response结构体为http JSON 响应
func encodeJsonResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// 解码业务逻辑中出现的err到http响应
func encodeErr(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
