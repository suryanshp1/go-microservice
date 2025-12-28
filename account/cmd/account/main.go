package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/suryanshp1/go-microservice/account"
	"github.com/tinrab/retry"
)

type config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Printf("failed to connect to database: %v", err)
		}
		return
	})
	defer r.Close()
	log.Println("Listening on Port 8080")

	s := account.NewService(r)
	log.Fatal(account.ListenGRPC(s, 8080))
}
