package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-crunchyroll/config"
	appgraphql "backend-crunchyroll/graphql"
	"backend-crunchyroll/supabase"

	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/handler"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Inicializar logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic("Erro ao iniciar logger: " + err.Error())
	}
	defer logger.Sync()

	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		logger.Warn(".env não encontrado, usando variáveis do sistema")
	}

	// Inicializar Supabase
	supabaseClient := supabase.GetClient()

	// Criar Schema (resolver é instanciado dentro de NewSchema)
	schema, err := appgraphql.NewSchema(supabaseClient)
	if err != nil {
		logger.Fatal("Erro ao criar schema GraphQL", zap.Error(err))
	}

	// Configurar servidor HTTP
	server := &http.Server{
		Addr:         ":" + config.GetPort(),
		Handler:      setupRouter(schema, logger),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Canal para shutdown gracioso
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Rodar servidor em goroutine
	go func() {
		logger.Info("Servidor rodando", zap.String("url", "http://localhost:"+config.GetPort()+"/graphql"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Erro ao iniciar servidor", zap.Error(err))
		}
	}()

	// Aguardar sinal para shutdown
	<-shutdownChan
	logger.Info("Desligando servidor...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Erro no desligamento", zap.Error(err))
	} else {
		logger.Info("Servidor desligado com sucesso")
	}
}

// setupRouter configura rotas e middlewares
func setupRouter(schema graphql.Schema, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	// Handler GraphQL
	handler := gqlhandler.New(&gqlhandler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	mux.Handle("/graphql", corsMiddleware(loggingMiddleware(handler, logger)))

	// Healthcheck
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
	})

	return mux
}

// corsMiddleware habilita CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware faz log de cada requisição HTTP
func loggingMiddleware(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		defer func() {
			logger.Info("Requisição",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", ww.status),
				zap.Duration("duration", time.Since(start)),
				zap.String("ip", r.RemoteAddr),
			)
		}()

		next.ServeHTTP(ww, r)
	})
}

// responseWriter captura o status de resposta
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
