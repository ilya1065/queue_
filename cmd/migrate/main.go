package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

func main() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		// используй то же, что у тебя в config.yml по умолчанию
		dsn = "file:./data/app.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=foreign_keys(ON)"
	}

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./db/migrations"
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// таблица учёта миграций
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY);`); err != nil {
		log.Fatalf("create schema_migrations: %v", err)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("read migrations dir: %v", err)
	}

	var ups []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".up.sql") {
			ups = append(ups, filepath.Join(migrationsDir, name))
		}
	}
	sort.Strings(ups)

	for _, path := range ups {
		base := filepath.Base(path)
		version := strings.TrimSuffix(base, ".up.sql")

		var exists int
		err := db.QueryRow(`SELECT COUNT(1) FROM schema_migrations WHERE version = ?`, version).Scan(&exists)
		if err != nil {
			log.Fatalf("check migration %s: %v", version, err)
		}
		if exists > 0 {
			continue
		}

		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("read %s: %v", path, err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("begin tx: %v", err)
		}

		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			log.Fatalf("apply %s: %v", base, err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations(version) VALUES(?)`, version); err != nil {
			_ = tx.Rollback()
			log.Fatalf("record %s: %v", base, err)
		}
		if err := tx.Commit(); err != nil {
			log.Fatalf("commit %s: %v", base, err)
		}

		fmt.Printf("applied %s\n", base)
	}

	fmt.Println("migrations done")
}
