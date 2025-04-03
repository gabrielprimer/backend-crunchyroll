package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/jefersonprimer/backend-crunchyroll/graphql"
)

func main() {
	schema, err := graphql.NewSchema()
	if err != nil {
		log.Fatalf("Erro ao criar schema: %v", err)
	}

	// Configurar o handler GraphQL
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Configurar o servidor HTTP
	http.Handle("/graphql", h)
	
	// Servir arquivos estáticos (opcional, para o frontend)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Servidor GraphQL rodando em http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))
}