package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const DefaultFile = "config.yaml"

type (
	Config struct {
		Env       string
		HTTP      HTTP
		DB        DB
		Seeder    Seeder
		Log       Log
		RateLimit RateLimit
	}

	HTTP struct {
		Host           string
		Port           int
		AllowedOrigins []string
		ReadTimeout    time.Duration
		WriteTimeout   time.Duration
	}

	DB struct {
		Host             string
		Port             int
		Name             string
		User             string
		Password         string
		SSLMode          string
		MaxConns         int
		MinConns         int
		LockTimeout      time.Duration
		StatementTimeout time.Duration
	}

	Seeder struct {
		MigrationsDir string
		Placeholders  []string
	}

	Log struct {
		Format string
		Level  string
	}

	RateLimit struct {
		MaxPerMinute int
	}
)

func (c *Config) Validate() error {
	switch c.Env {
	case "development", "staging", "production":
	case "":
		return errors.New("env is empty")
	default:
		return fmt.Errorf("unknown environment %q", c.Env)
	}

	if c.HTTP.Port <= 0 || c.HTTP.Port > 65535 {
		return fmt.Errorf("http.port must be 1..65535, got %d", c.HTTP.Port)
	}
	if c.Env == "production" && len(c.HTTP.AllowedOrigins) == 0 {
		return errors.New("http.allowed_origins is required in production")
	}

	if c.DB.Host == "" || c.DB.Name == "" || c.DB.User == "" || c.DB.Password == "" {
		return errors.New("db host, name, user and password are required")
	}
	if c.DB.Port <= 0 || c.DB.Port > 65535 {
		return fmt.Errorf("db.port must be 1..65535, got %d", c.DB.Port)
	}

	if c.Log.Level == "" {
		return errors.New("log.level is required")
	}
	switch c.Log.Format {
	case "json", "text":
	default:
		return fmt.Errorf("log.format must be json or text, got %q", c.Log.Format)
	}
	return nil
}

func (c *Config) Secret(name string) string {
	return os.Getenv(name)
}

func Load(path string) (*Config, error) {
	_ = godotenv.Load()

	v := viper.New()
	v.SetEnvPrefix("minning")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	if path == "" {
		path = DefaultFile
	}
	if _, statErr := os.Stat(path); statErr == nil {
		v.SetConfigFile(path)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("config: read %s: %w", path, err)
		}
	}

	v.SetDefault("env", "development")
	v.SetDefault("http.port", 8080)
	v.SetDefault("http.read_timeout", "30s")
	v.SetDefault("http.write_timeout", "30s")
	v.SetDefault("db.port", 5432)
	v.SetDefault("db.sslmode", "require")
	v.SetDefault("db.max_conns", 5)
	v.SetDefault("log.format", "text")
	v.SetDefault("log.level", "debug")
	v.SetDefault("ratelimit.max_per_minute", 20)

	c := &Config{
		Env: v.GetString("env"),
		HTTP: HTTP{
			Host:           v.GetString("http.host"),
			Port:           v.GetInt("http.port"),
			AllowedOrigins: parseAllowedOrigins(v.GetString("http.allowed_origins")),
			ReadTimeout:    v.GetDuration("http.read_timeout"),
			WriteTimeout:   v.GetDuration("http.write_timeout"),
		},
		DB: DB{
			Host:             v.GetString("db.host"),
			Port:             v.GetInt("db.port"),
			Name:             v.GetString("db.name"),
			User:             v.GetString("db.user"),
			Password:         v.GetString("db.password"),
			SSLMode:          v.GetString("db.sslmode"),
			MaxConns:         v.GetInt("db.max_conns"),
			MinConns:         v.GetInt("db.min_conns"),
			LockTimeout:      v.GetDuration("db.lock_timeout"),
			StatementTimeout: v.GetDuration("db.statement_timeout"),
		},
		Seeder: Seeder{
			MigrationsDir: v.GetString("migration.dir_path"),
			Placeholders:  v.GetStringSlice("migration.placeholders"),
		},
		Log: Log{
			Format: v.GetString("log.format"),
			Level:  v.GetString("log.level"),
		},
		RateLimit: RateLimit{
			MaxPerMinute: v.GetInt("ratelimit.max_per_minute"),
		},
	}

	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return c, nil
}

func parseAllowedOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}
	parts := strings.Split(origins, ",")
	var result []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (c *DB) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode,
	)
}

func (c *HTTP) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
