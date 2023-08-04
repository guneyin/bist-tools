package app

import (
	"github.com/guneyin/bist-tools/pkg/config"
	"github.com/guneyin/bist-tools/pkg/database"
	"github.com/guneyin/bist-tools/pkg/httpserver"
	"github.com/guneyin/bist-tools/pkg/logger"
	"github.com/guneyin/bist-tools/pkg/migration"
	"github.com/guneyin/bist-tools/pkg/redis"
	"log"
)

func Run() {
	err := config.Init()
	if err != nil {
		log.Fatalf("configuration error: %s", err)
	}

	logger.Init()

	err = redis.Init()
	if err != nil {
		log.Fatalf("redis error: %s", err)
	}

	err = database.Init()
	if err != nil {
		log.Fatalf("database error: %s", err)
	}

	err = migration.Init()
	if err != nil {
		log.Fatalf("migration error: %s", err)
	}

	err = httpserver.Init()
	if err != nil {
		log.Fatalf("http server error: %s", err)
	}
}
