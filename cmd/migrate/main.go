package main

import (
	"log"
	"os"
	"slices"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/shivamkrch/olx-clone-api/internal/config"
)

func main() {
	commands := []string{"up", "down"}
	if len(os.Args) != 2 || !slices.Contains(commands, os.Args[1]) {
		log.Fatalf("usage: migrate <%s>", strings.Join(commands, " | "))
	}

	cfg := config.MustLoad()

	m, err := migrate.New("file://migrations", cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}

	cmd := os.Args[1]
	switch cmd {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}

	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("unknown command %s", cmd)
	}

	log.Println("running migration")
}
