package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/AdityaP1502/Instant-Messaging/api/api/routes"
	"github.com/AdityaP1502/Instant-Messaging/api/api/util"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	path := "./app.config.json"
	config, err := util.ReadJSONConfiguration(path)

	if err != nil {
		log.Fatal(err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.Username, config.Database.Password, config.Database.Database)

	db, err := sqlx.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	routes.SetAccountRoute(r.PathPrefix("/v1").Subrouter(), db.DB, config)
	// // r.Handle("/", r)

	r.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("WTF is happening")
		w.Write([]byte("Hello, world!"))
		w.WriteHeader(200)
	}).Methods("GET")

	// a := s.PathPrefix("/account").Subrouter()

	// a.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello, world!"))
	// }).Methods("GET")

	http.Handle("/", r)

	// wait until the server has ended
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if strings.ToLower(config.Server.Secure) == "true" {
			http.ListenAndServeTLS(
				fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
				config.Certificate.CertFile,
				config.Certificate.KeyFile,
				r,
			)
		} else {
			err := http.ListenAndServe(
				fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
				r,
			)

			if err != nil {
				panic(err)
			}
		}
	}()

	fmt.Printf("Server is running on %s:%d\n", config.Server.Host, config.Server.Port)
	wg.Wait()

}
