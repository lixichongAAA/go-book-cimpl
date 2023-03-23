package main

import (
	"fmt"
	"log"
	"net/rpc"

	service "github.com/longjoy/micro-go-book/ch7-rpc/basic/string-service"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	stringReq := &service.StringRequest{A: "A", B: "B"}
	// Synchronous call
	var reply string
	err = client.Call("StringService.Concat", stringReq, &reply)
	if err != nil {
		log.Fatal("StringService error:", err)
	}
	fmt.Printf("StringService Concat : %s concat %s = %s\n", stringReq.A, stringReq.B, reply)
	// Asynchronous call
	stringReq = &service.StringRequest{A: "ACD", B: "BDF"}
	call := client.Go("StringService.Diff", stringReq, &reply, nil)
	_ = <-call.Done
	fmt.Printf("StringService Diff : %s diff %s = %s\n", stringReq.A, stringReq.B, reply)

}
