package database

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var L *zap.Logger

func InitLogger() {
	var err error
	env := os.Getenv("GO_ENV")

	if env == "production" {
		L, err = zap.NewProduction()
	} else {
		L, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatalf("Failed to initialize Zap logger: %v", err)
	}
	zap.ReplaceGlobals(L)
}

func SyncLogger() {
	if L != nil {
		L.Sync()
	}
}
