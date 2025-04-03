package supabase

import (
	"log"
	"os"

	"github.com/nedpals/supabase-go"
)

func NewClient() *supabase.Client {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	if url == "" || key == "" {
		log.Fatal("Supabase credentials not found in environment variables")
	}

	client := supabase.CreateClient(url, key)
	
	// Teste a conexão
	var result []map[string]interface{}
	err := client.DB.From("animes").Select("*").Limit(1).Execute(&result)
	if err != nil {
		log.Fatalf("Failed to connect to Supabase: %v", err)
	}
	
	return client
}