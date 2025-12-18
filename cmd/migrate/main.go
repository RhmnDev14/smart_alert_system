package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Get database URL
	dbURL := getDatabaseURL()
	if dbURL == "" {
		log.Fatal("DATABASE_URL or DB_* variables not set")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✓ Connected to database")
	fmt.Printf("Database URL: %s\n\n", maskPassword(dbURL))

	// Get migrations directory
	migrationsDir := "migrations"
	if len(os.Args) > 1 {
		migrationsDir = os.Args[1]
	}

	// Read and sort migration files
	files, err := getMigrationFiles(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migration files: %v", err)
	}

	if len(files) == 0 {
		log.Fatal("No migration files found")
	}

	fmt.Printf("Found %d migration files\n\n", len(files))

	// Run migrations
	for _, file := range files {
		fmt.Printf("Running: %s\n", filepath.Base(file))
		
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read file %s: %v", file, err)
		}

		// Execute migration
		if _, err := db.Exec(string(content)); err != nil {
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}

		fmt.Printf("✓ Success: %s\n\n", filepath.Base(file))
	}

	fmt.Println("✓ All migrations completed successfully!")
}

func getDatabaseURL() string {
	// Check for DATABASE_URL first
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	// Build from individual variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || user == "" || password == "" || dbname == "" {
		return ""
	}

	if port == "" {
		port = "5432"
	}
	if sslmode == "" {
		sslmode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
}

func getMigrationFiles(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		
		name := file.Name()
		// Only include .sql files, exclude README, drop_all_tables, and other non-migration files
		if strings.HasSuffix(name, ".sql") && 
		   name != "README.md" && 
		   name != "drop_all_tables.sql" {
			migrations = append(migrations, filepath.Join(dir, name))
		}
	}

	// Sort by filename
	sort.Strings(migrations)

	return migrations, nil
}

func maskPassword(url string) string {
	// Mask password in database URL for display
	parts := strings.Split(url, "@")
	if len(parts) != 2 {
		return url
	}
	
	auth := parts[0]
	if strings.Contains(auth, ":") {
		authParts := strings.Split(auth, ":")
		if len(authParts) >= 2 {
			authParts[1] = "***"
			auth = strings.Join(authParts, ":")
		}
	}
	
	return auth + "@" + parts[1]
}

