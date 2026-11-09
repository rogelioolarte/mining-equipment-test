package postgres

import (
	"context"
	"errors"
	"fmt"
	"io/fs"

	"minning/internal/domain/ports"
	"minning/internal/infra/config"
	"minning/internal/infra/logger"

	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// matches every gen_random_uuid() call in the SQL text
var genRandomUUIDPattern = regexp.MustCompile(`gen_random_uuid\(\)`)

type SeederService struct {
	m   ports.TransactionManager
	cfg *config.Seeder
	log *logger.Logger
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func NewSeederService(logger *logger.Logger, m ports.TransactionManager, cfg *config.Seeder) (*SeederService, error) {
	if logger == nil {
		return nil, fmt.Errorf("no logger provided")
	}
	if cfg == nil {
		return nil, fmt.Errorf("no config for seeder provided")
	}
	logger.Info(fmt.Sprintf("seeder service initialized: %v", cfg.MigrationsDir))
	return &SeederService{m: m, cfg: cfg, log: logger}, nil
}

func (s *SeederService) ApplyAll(ctx context.Context) error {
	s.log.Info("starting to apply all migrations")

	err := filepath.WalkDir(s.cfg.MigrationsDir, func(path string, info fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			s.log.Error("error walking migrations dir", "path", path, "error", walkErr)
			return walkErr
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".sql") {
			return nil
		}
		s.log.Info("processing migration file", "file", path)
		return s.applyFile(ctx, path)
	})
	if err != nil {
		s.log.Error("error applying migrations", "error", err)
		return err
	}

	s.log.Info("all migrations applied successfully")
	return nil
}

func (s *SeederService) applyFile(ctx context.Context, filePath string) error {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		s.log.Error("failed to read file", "file", filePath, "error", err)
		return err
	}

	sqlText := s.replacePlaceholders(string(contents))

	txCtx, err := s.m.Begin(ctx)
	if err != nil {
		s.log.Error("failed to begin transaction", "file", filePath, "error", err)
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer s.m.Rollback(txCtx)

	if _, err := s.m.DB(txCtx).Exec(txCtx, "SAVEPOINT seeder_sp"); err != nil {
		return fmt.Errorf("failed to create savepoint: %w", err)
	}

	if _, execErr := s.m.DB(txCtx).Exec(txCtx, sqlText); execErr != nil {
		if isUniqueViolation(execErr) {
			if _, rbErr := s.m.DB(txCtx).Exec(txCtx, "ROLLBACK TO SAVEPOINT seeder_sp"); rbErr != nil {
				return fmt.Errorf("failed to rollback to savepoint: %w", rbErr)
			}
			s.log.Warn("record already exists, skipping", "file", filePath)
			// Commit limpio — la tx sigue válida tras el rollback al savepoint
			return s.m.Commit(txCtx)
		}
		s.log.Error("error executing SQL file", "file", filePath, "error", execErr)
		return execErr // defer ejecuta Rollback
	}

	if _, err := s.m.DB(txCtx).Exec(txCtx, "RELEASE SAVEPOINT seeder_sp"); err != nil {
		return fmt.Errorf("failed to release savepoint: %w", err)
	}

	if err := s.m.Commit(txCtx); err != nil {
		s.log.Error("failed to commit transaction", "file", filePath, "error", err)
		return fmt.Errorf("commit transaction: %w", err)
	}

	s.log.Info("migration applied", "file", filePath)
	return nil
}

// replacePlaceholders handles two kinds of substitution:
//  1. Config-defined named placeholders  → replaced with a fresh UUIDv7 each.
//  2. gen_random_uuid() outside DO $$ blocks → replaced by UUIDv7
func (s *SeederService) replacePlaceholders(sqlText string) string {
	// 1. Named placeholders from config
	for _, ph := range s.cfg.Placeholders {
		if strings.Contains(sqlText, ph) {
			newID := uuid.Must(uuid.NewV7()).String()
			sqlText = strings.ReplaceAll(sqlText, ph, fmt.Sprintf("'%s'", newID))
			s.log.Debug("replaced placeholder", "placeholder", ph, "uuidv7", newID)
		}
	}

	// 2. Reemplazo gen_random_uuid() **solo fuera de bloques DO $$**
	parts := regexp.MustCompile(`(?s)(DO \$\$.*?\$\$)`).Split(sqlText, -1)
	matches := regexp.MustCompile(`(?s)(DO \$\$.*?\$\$)`).FindAllString(sqlText, -1)

	for i, part := range parts {
		parts[i] = genRandomUUIDPattern.ReplaceAllStringFunc(part, func(_ string) string {
			newID := uuid.Must(uuid.NewV7()).String()
			s.log.Debug("replaced gen_random_uuid() outside DO $$", "uuidv7", newID)
			return fmt.Sprintf("'%s'", newID)
		})
	}

	// recombinar partes: intercalar bloques DO $$ originales sin tocarlos
	result := ""
	for i := 0; i < len(parts); i++ {
		result += parts[i]
		if i < len(matches) {
			result += matches[i] // bloques DO $$ intactos
		}
	}

	return result
}
