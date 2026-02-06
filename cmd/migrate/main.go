package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"assetManager/internal/config"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	migrationsDir := flag.String("migrations", "migrations", "Path to migrations directory")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := sqlx.Connect("mysql", cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create migrations table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get applied migrations
	var applied []string
	err = db.Select(&applied, "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		log.Fatalf("Failed to get applied migrations: %v", err)
	}
	appliedMap := make(map[string]bool)
	for _, v := range applied {
		appliedMap[v] = true
	}

	// Read migration files
	files, err := os.ReadDir(*migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	var migrations []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".sql") {
			migrations = append(migrations, f.Name())
		}
	}
	sort.Strings(migrations)

	// Apply pending migrations
	for _, migration := range migrations {
		version := strings.TrimSuffix(migration, ".sql")
		if appliedMap[version] {
			log.Printf("Skipping %s (already applied)", migration)
			continue
		}

		log.Printf("Applying %s...", migration)

		content, err := os.ReadFile(filepath.Join(*migrationsDir, migration))
		if err != nil {
			log.Fatalf("Failed to read migration file: %v", err)
		}

		// Remove comments and split by semicolon
		cleanContent := removeComments(string(content))
		statements := strings.Split(cleanContent, ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			_, err = db.Exec(stmt)
			if err != nil {
				log.Fatalf("Failed to execute migration %s: %v\nStatement: %s", migration, err, stmt)
			}
		}

		// Record migration
		_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
		if err != nil {
			log.Fatalf("Failed to record migration: %v", err)
		}

		log.Printf("Applied %s", migration)
	}

	fmt.Println("Migrations complete!")
}

// removeComments strips SQL comments from content
func removeComments(content string) string {
	var result strings.Builder
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		// Remove inline comments (-- comment)
		if idx := strings.Index(line, "--"); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line != "" {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	return result.String()
}
