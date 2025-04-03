package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/jefersonprimer/backend-crunchyroll/graphql"
	"github.com/jefersonprimer/backend-crunchyroll/supabase"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inicializar Supabase
	supabaseClient := supabase.NewClient()

	// Configurar GraphQL
	schema, err := graphql.NewSchema(supabaseClient)
	if err != nil {
		log.Fatalf("Error creating schema: %v", err)
	}

	// Configurar servidor HTTP
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	log.Println("Server running at http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}