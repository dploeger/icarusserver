package main

import (
	"github.com/gin-gonic/gin"
	"icarusserver/internal/endpoints"
	"log"
	"net/http"
	"os"
)

var endpointRegistry = []endpoints.Endpoint{
	endpoints.ProcessorsEndpoint{},
}

func main() {
	bindAddress := "127.0.0.1:8080"
	if a, ok := os.LookupEnv("BIND_ADDRESS"); ok {
		bindAddress = a
	}
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusFound, "/client")
	})

	r.Static("/client", "./client")

	api := r.Group("/api")

	for _, endpoint := range endpointRegistry {
		if err := endpoint.Register(api); err != nil {
			log.Fatalf("Error registering endpoint: %s", err)
		}
	}

	if err := r.Run(bindAddress); err != nil {
		log.Fatalf("Error running server: %s", err)
	}
}
