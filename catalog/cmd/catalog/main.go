package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/suryanshp1/go-microservice/catalog"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatal(err)
	}

	var r catalog.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = catalog.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println("Error connecting to elastic database:", err)
		}
		return
	})

	defer r.Close()

	log.Println("Connected to ElasticSearch, Listening on port 8080")

	s := catalog.NewService(r)

	log.Fatal(catalog.ListenGRPC(s, 8080))

}
