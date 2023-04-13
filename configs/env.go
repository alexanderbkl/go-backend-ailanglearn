package configs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGO_DB_URI")
}