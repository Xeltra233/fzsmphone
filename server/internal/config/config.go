package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	JWTSecret              string
	DiscordClientID        string
	DiscordSecret          string
	DiscordRedirect        string
	DiscordGuildID         string
	DiscordBotToken        string
	DiscordRequiredRoleIDs []string
	DiscordRoleMatch       string // "all" or "any"
	DiscordAuthFailMessage string
	FrontendURL            string
	Environment            string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:                   getEnv("PORT", "8080"),
		DatabaseURL:            getEnv("DATABASE_URL", ""),
		JWTSecret:              getEnv("JWT_SECRET", ""),
		DiscordClientID:        getEnv("DISCORD_CLIENT_ID", ""),
		DiscordSecret:          getEnv("DISCORD_CLIENT_SECRET", ""),
		DiscordRedirect:        getEnv("DISCORD_REDIRECT_URI", "http://localhost:3000/auth/callback"),
		DiscordGuildID:         getEnv("DISCORD_GUILD_ID", ""),
		DiscordBotToken:        getEnv("DISCORD_BOT_TOKEN", ""),
		DiscordRequiredRoleIDs: parseCSV(getEnv("DISCORD_REQUIRED_ROLE_IDS", "")),
		DiscordRoleMatch:       getEnv("DISCORD_ROLE_MATCH", "all"),
		DiscordAuthFailMessage: getEnv("DISCORD_AUTH_FAIL_MESSAGE", "请先加入指定服务器并完成身份组验证"),
		FrontendURL:            getEnv("FRONTEND_URL", "http://localhost:3000"),
		Environment:            getEnv("ENVIRONMENT", "development"),
	}

	if cfg.DatabaseURL == "" {
		// Build from individual params
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		pass := getEnv("DB_PASSWORD", "postgres")
		name := getEnv("DB_NAME", "fzsmphone")
		sslmode := getEnv("DB_SSLMODE", "disable")
		cfg.DatabaseURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, pass, host, port, name, sslmode,
		)
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}

func (c *Config) IsDev() bool {
	return c.Environment == "development"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
