package config

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Port        int
	DatabaseURL string
	JWTSecret   string
	AdminLogin  string
	AdminPass   string
}

// Load reads default.env into environment.
func Load() *Config {
	loadDotEnv()

	port, err := strconv.Atoi(getEnv("PORT", "0"))
	if err != nil || port == 0 {
		port = 8080
	}

	return &Config{
		Port:        port,
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/appdb?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecretkey"),
		AdminLogin:  getEnv("ADMIN_LOGIN", "admin"),
		AdminPass:   getEnv("ADMIN_PASSWORD", "securepassword123"),
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// loadDotEnv searches for default.env and sets its entries into environment.
func loadDotEnv() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("config: unable to get working directory: %v", err)
		return
	}

	var envPath string
	dir := cwd
	for {
		attempt := filepath.Join(dir, "default.env")
		if _, err := os.Stat(attempt); err == nil {
			envPath = attempt
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	if envPath == "" {
		return
	}

	file, err := os.Open(envPath)
	if err != nil {
		log.Printf("config: failed to open %s: %v", envPath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) >= 2 {
			value = value[1 : len(value)-1]
		}
		os.Setenv(key, value)
	}
	if err := scanner.Err(); err != nil {
		log.Printf("config: error reading %s: %v", envPath, err)
	}
}
