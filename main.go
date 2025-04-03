package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/jefersonprimer/backend-crunchyroll/config"
	"github.com/jefersonprimer/backend-crunchyroll/graphql"
	"github.com/jefersonprimer/backend-crunchyroll/supabase"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: Arquivo .env não encontrado. Usando variáveis de ambiente do sistema.")
	}

	// Inicializar Supabase
	supabaseClient := supabase.NewClient()

	// Configurar GraphQL
	schema, err := graphql.NewSchema(supabaseClient)
	if err != nil {
		log.Fatalf("Erro ao criar schema GraphQL: %v", err)
	}

	// Configurar CORS
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			h.ServeHTTP(w, r)
		})
	}

	// Configurar servidor HTTP
	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   true,
	})

	// Obter porta do config
	port := config.GetPort()

	// Rota principal GraphQL
	http.Handle("/graphql", corsHandler(h))
	
	// Rota de verificação de saúde
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API Funcionando!"))
	})

	log.Printf("Servidor rodando em http://localhost:%s/graphql", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}