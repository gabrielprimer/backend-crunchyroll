package supabase

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nedpals/supabase-go"
)

var (
	instance *supabase.Client
	once     sync.Once
)

// ClientConfig configurações avançadas
type ClientConfig struct {
	Timeout        time.Duration
	MaxRetries     int
	EnableMetrics  bool
}

// DefaultConfig configuração padrão otimizada
var DefaultConfig = ClientConfig{
	Timeout:        5 * time.Second,
	MaxRetries:     3,
	EnableMetrics:  true,
}

// GetClient retorna uma instância singleton thread-safe do cliente Supabase
func GetClient() *supabase.Client {
	once.Do(func() {
		cfg := DefaultConfig

		url := mustGetEnv("SUPABASE_URL")
		key := mustGetEnv("SUPABASE_KEY")

		// Criação do cliente com timeout
		client := supabase.CreateClient(url, key)

		// Teste de conexão robusto
		if err := testConnection(client, cfg); err != nil {
			log.Fatalf("Supabase connection failed: %v", err)
		}

		instance = client
	})

	return instance
}

// testConnection verifica a conexão com tratamento de retry
func testConnection(client *supabase.Client, cfg ClientConfig) error {
	var lastErr error

	for i := 0; i < cfg.MaxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()

		var result []map[string]interface{}
		err := client.DB.From("animes").Select("*").Limit(1).ExecuteWithContext(ctx, &result)

		if err == nil {
			return nil
		}

		lastErr = err
		time.Sleep(time.Duration(i+1) * 500 * time.Millisecond) // Backoff exponencial
	}

	return lastErr
}

// mustGetEnv obtém variável de ambiente ou falha
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Variável de ambiente %s não definida", key)
	}
	return value
}

// Métodos adicionais para performance:

// WithContext adiciona timeout padrão às operações
func WithContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, DefaultConfig.Timeout)
}