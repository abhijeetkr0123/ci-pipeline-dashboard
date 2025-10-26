package config

import (
	"log"
	"os"
)

type Config struct {
	SupabaseUrl string
	SupabaseKey string
	Port        string
}

func LoadConfig() *Config {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if url == "" || key == "" {
		log.Fatal("SUPABASE_URL or SUPABASE_KEY is missing in .env")
	}

	return &Config{
		SupabaseUrl: url,
		SupabaseKey: key,
		Port:        port,
	}
}
