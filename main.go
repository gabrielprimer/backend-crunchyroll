package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/seu-usuario/anime-api/graphql"
)

func main() {
	// Inicializar schema GraphQL
	schema, err := graphql.NewSchema()
	if err != nil {
		log.Fatalf("Erro ao criar schema: %v", err)
	}

	// Configurar handler GraphQL
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Configurar rota
	http.Handle("/graphql", h)

	// Iniciar servidor
	log.Println("Servidor GraphQL rodando em http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}