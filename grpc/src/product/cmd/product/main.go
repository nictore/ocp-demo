package main

import (
	"sync"

	"github.com/iamdpastore/ocp-demo/grpc/src/product/internal/server/grpc"
	"github.com/iamdpastore/ocp-demo/grpc/src/product/internal/server/http"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go grpc.Serve(&wg, "5000")

	wg.Add(1)
	go http.Serve(&wg, "5000", "8080")

	wg.Wait()
}
