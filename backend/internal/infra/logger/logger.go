// Package logger builds the structured logger every layer takes as
// *slog.Logger (NFR-OBS-01): one JSON or text handler whose level comes
// from configuration, so a call site never parses its own log level.
//
// Logger is a TYPE ALIAS of slog.Logger on purpose: existing constructors
// read `logger *slog.Logger` (use cases, handlers, the recovery middleware),
// while AGENTS.md names the port `*logger.Logger` — the alias makes both
// syntactic forms the same runtime type, and no logging call site changes.
package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Logger is the shared logger type (see the package comment).
type Logger = slog.Logger

// levelNames maps the configured level names; lowercase is accepted.
var levelNames = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

// ParseLevel maps a configuration level name to a slog.Level. It accepts
// the four canonical names case-insensitively.
func ParseLevel(level string) (slog.Level, error) {
	if parsed, ok := levelNames[strings.ToLower(level)]; ok {
		return parsed, nil
	}
	return 0, fmt.Errorf("logger: unknown level %q (want debug, info, warn or error)", level)
}

// (default, machine-readable for the structured pipeline) or text
// (developer consoles); the returned value is safe for concurrent use and
// ready to share application-wide through .With/.WithGroup.
//
// Correlation stays where the transport put it (the request-id middleware
// writes it into the request context and the call sites add the
// "request_id" attribute — see internal/http/middlewares): this package
// builds the handler, it does not centralize another context key.
func New(format, level string) (*Logger, error) {
	if format == "" {
		format = "json"
	}
	if level == "" {
		level = "info"
	}
	lvl, err := ParseLevel(level)
	if err != nil {
		return nil, err
	}
	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: lvl}
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stderr, opts)
	default:
		return nil, fmt.Errorf("logger: unknown format %q (want json or text)", format)
	}
	return slog.New(handler), nil
}
