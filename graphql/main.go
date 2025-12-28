package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountURL string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogURL string `envconfig:"CATALOG_SERVICE_URL"`
	OrderURL   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)

	if err != nil {
		log.Fatalf("Failed to process envconfig: %v", err)
	}

	// Validate required environment variables
	if cfg.AccountURL == "" {
		log.Fatal("ACCOUNT_SERVICE_URL is required but not set")
	}
	if cfg.CatalogURL == "" {
		log.Fatal("CATALOG_SERVICE_URL is required but not set")
	}
	if cfg.OrderURL == "" {
		log.Fatal("ORDER_SERVICE_URL is required but not set")
	}

	log.Printf("Connecting to services:")
	log.Printf("  Account: %s", cfg.AccountURL)
	log.Printf("  Catalog: %s", cfg.CatalogURL)
	log.Printf("  Order:   %s", cfg.OrderURL)

	s, err := NewGraphQLServer(cfg.AccountURL, cfg.CatalogURL, cfg.OrderURL)

	if err != nil {
		log.Fatalf("Failed to create GraphQL server: %v", err)
	}

	http.Handle("/graphql", handler.NewDefaultServer(s.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
