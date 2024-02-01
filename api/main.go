package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
	"github.com/gorilla/mux"
)

func main() {
	path := "./app.config.json"
	config, err := util.ReadJSONConfiguration(path)

	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	// wait until the server has ended
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if strings.ToLower(config.Server.Secure) == "true" {
			http.ListenAndServeTLS(
				fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
				config.Certificate.CertFile,
				config.Certificate.KeyFile,
				r,
			)
		} else {
			http.ListenAndServe(
				fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
				r,
			)
		}
	}()

	fmt.Printf("Server is running on %s:%s\n", config.Server.Host, config.Server.Port)
	wg.Wait()
}
